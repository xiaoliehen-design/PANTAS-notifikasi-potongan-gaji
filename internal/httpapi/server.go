package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"time"

	"github.com/bcpriok/pantas/internal/auth"
	"github.com/bcpriok/pantas/internal/config"
	"github.com/bcpriok/pantas/internal/importer"
	"github.com/bcpriok/pantas/internal/storage"
	webassets "github.com/bcpriok/pantas/web"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	pool    *pgxpool.Pool
	cfg     config.Config
	auth    *auth.Service
	imports *importer.Service
	storage *storage.Client
	log     *slog.Logger
	static  http.Handler
	index   []byte
}

type handlerFunc func(http.ResponseWriter, *http.Request, auth.Principal)

func New(pool *pgxpool.Pool, cfg config.Config, authService *auth.Service, importService *importer.Service, storageClient *storage.Client, logger *slog.Logger) (*App, error) {
	staticFS, err := fs.Sub(webassets.Static, "static")
	if err != nil {
		return nil, err
	}
	index, err := fs.ReadFile(staticFS, "index.html")
	if err != nil {
		return nil, err
	}
	app := &App{
		pool: pool, cfg: cfg, auth: authService, imports: importService,
		storage: storageClient, log: logger, static: http.FileServer(http.FS(staticFS)), index: index,
	}
	return app, nil
}

func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", a.health)
	mux.HandleFunc("GET /api/auth/captcha", a.captcha)
	mux.HandleFunc("POST /api/auth/login", a.login)
	mux.HandleFunc("POST /api/auth/forgot-password", a.forgotPassword)
	mux.HandleFunc("POST /api/auth/reset-password", a.resetPassword)
	mux.Handle("GET /api/auth/me", a.withAuth(a.me))
	mux.Handle("POST /api/auth/logout", a.withAuth(a.logout))
	mux.Handle("POST /api/auth/change-password", a.withAuth(a.changePassword))
	mux.Handle("POST /api/profile/contact/start", a.withUser(a.startContactChange))
	mux.Handle("POST /api/profile/contact/verify", a.withUser(a.verifyContactChange))
	mux.Handle("GET /api/dashboard", a.withUser(a.dashboard))
	mux.Handle("GET /api/history", a.withUser(a.history))
	mux.Handle("GET /api/deductions", a.withUser(a.deductions))
	mux.Handle("GET /api/notifications", a.withAuth(a.notifications))
	mux.Handle("POST /api/notifications/read-all", a.withAuth(a.markAllNotificationsRead))
	mux.Handle("POST /api/notifications/{id}/read", a.withAuth(a.markNotificationRead))
	mux.Handle("GET /api/monitoring", a.withSupervisor(a.monitoring))
	mux.Handle("GET /api/warnings", a.withSupervisor(a.warnings))
	mux.Handle("GET /api/appeals", a.withUser(a.appeals))
	mux.Handle("GET /api/appeals/options", a.withUser(a.appealOptions))
	mux.Handle("POST /api/appeals", a.withUser(a.createAppeal))
	mux.Handle("POST /api/appeals/items/{id}/document", a.withUser(a.uploadAppealDocument))
	mux.Handle("GET /api/appeals/items/{id}/documents", a.withAuth(a.appealDocuments))
	mux.Handle("GET /api/reviews/supervisor", a.withSupervisor(a.supervisorQueue))
	mux.Handle("POST /api/reviews/supervisor/{id}", a.withSupervisor(a.supervisorDecision))
	mux.Handle("GET /api/reviews/admin", a.withAdmin(a.adminReviewQueue))
	mux.Handle("POST /api/reviews/admin/{id}", a.withAdmin(a.adminDecision))
	mux.Handle("GET /api/documents/{id}", a.withAuth(a.downloadDocument))
	mux.Handle("GET /api/admin/users", a.withAdmin(a.adminUsers))
	mux.Handle("POST /api/admin/users", a.withAdmin(a.adminCreateUser))
	mux.Handle("PATCH /api/admin/users/{id}", a.withAdmin(a.adminUpdateUser))
	mux.Handle("DELETE /api/admin/users/{id}", a.withAdmin(a.adminDeleteUser))
	mux.Handle("POST /api/admin/users/{id}/reset-password", a.withAdmin(a.adminResetUserPassword))
	mux.Handle("GET /api/admin/units", a.withAdmin(a.adminUnits))
	mux.Handle("GET /api/admin/parameters", a.withAdmin(a.adminParameters))
	mux.Handle("PATCH /api/admin/parameters/{key}", a.withAdmin(a.adminUpdateParameter))
	mux.Handle("GET /api/admin/rules", a.withAdmin(a.adminRules))
	mux.Handle("POST /api/admin/rules", a.withAdmin(a.adminCreateRule))
	mux.Handle("PATCH /api/admin/rules/{id}", a.withAdmin(a.adminUpdateRule))
	mux.Handle("GET /api/admin/reasons", a.withAdmin(a.adminReasons))
	mux.Handle("POST /api/admin/reasons", a.withAdmin(a.adminCreateReason))
	mux.Handle("PATCH /api/admin/reasons/{id}", a.withAdmin(a.adminUpdateReason))
	mux.Handle("GET /api/admin/imports", a.withAdmin(a.adminImports))
	mux.Handle("POST /api/admin/imports/preview", a.withAdmin(a.adminImportPreview))
	mux.Handle("POST /api/admin/imports/{id}/publish", a.withAdmin(a.adminImportPublish))
	mux.Handle("DELETE /api/admin/imports/{id}", a.withAdmin(a.adminImportReject))
	mux.HandleFunc("/", a.serveSPA)
	return a.securityHeaders(a.requestLog(a.recoverer(a.sameOrigin(mux))))
}

func (a *App) withAuth(next handlerFunc) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		principal, session, err := a.auth.AuthenticateRequest(request.Context(), request)
		if err != nil {
			writeError(response, http.StatusUnauthorized, "Sesi berakhir. Silakan login kembali.", "unauthenticated")
			return
		}
		if isUnsafe(request.Method) {
			if err := a.auth.VerifyCSRF(request, session); err != nil {
				writeError(response, http.StatusForbidden, "Token keamanan tidak valid. Muat ulang halaman.", "invalid_csrf")
				return
			}
		}
		if principal.MustChangePassword && request.URL.Path != "/api/auth/change-password" && request.URL.Path != "/api/auth/logout" && request.URL.Path != "/api/auth/me" {
			writeError(response, http.StatusPreconditionRequired, "Ganti password awal sebelum menggunakan PANTAS.", "password_change_required")
			return
		}
		request = request.WithContext(auth.WithPrincipal(request.Context(), principal))
		next(response, request, principal)
	})
}

func (a *App) withAdmin(next handlerFunc) http.Handler {
	return a.withAuth(func(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
		if !principal.IsAdmin {
			writeError(response, http.StatusForbidden, "Akses administrator diperlukan.", "forbidden")
			return
		}
		next(response, request, principal)
	})
}

func (a *App) withUser(next handlerFunc) http.Handler {
	return a.withAuth(func(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
		if principal.AccountType != "user" {
			writeError(response, http.StatusForbidden, "Menu pribadi hanya tersedia bagi akun pegawai.", "forbidden")
			return
		}
		next(response, request, principal)
	})
}

func (a *App) withSupervisor(next handlerFunc) http.Handler {
	return a.withAuth(func(response http.ResponseWriter, request *http.Request, principal auth.Principal) {
		if !principal.IsSupervisor() {
			writeError(response, http.StatusForbidden, "Menu ini hanya tersedia bagi atasan atau administrator.", "forbidden")
			return
		}
		next(response, request, principal)
	})
}

func (a *App) health(response http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), 2*time.Second)
	defer cancel()
	if err := a.pool.Ping(ctx); err != nil {
		writeJSON(response, http.StatusServiceUnavailable, map[string]any{"status": "degraded", "database": "unavailable"})
		return
	}
	writeJSON(response, http.StatusOK, map[string]any{"status": "ok", "database": "ok"})
}

func (a *App) serveSPA(response http.ResponseWriter, request *http.Request) {
	if strings.HasPrefix(request.URL.Path, "/api/") {
		writeError(response, http.StatusNotFound, "Endpoint tidak ditemukan.", "not_found")
		return
	}
	cleanPath := strings.TrimPrefix(request.URL.Path, "/")
	if cleanPath != "" && strings.Contains(cleanPath, ".") {
		response.Header().Set("Cache-Control", "no-cache")
		a.static.ServeHTTP(response, request)
		return
	}
	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	response.Header().Set("Cache-Control", "no-cache")
	response.WriteHeader(http.StatusOK)
	if request.Method != http.MethodHead {
		_, _ = response.Write(a.index)
	}
}

func (a *App) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("X-Content-Type-Options", "nosniff")
		response.Header().Set("X-Frame-Options", "DENY")
		response.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		response.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		response.Header().Set("Content-Security-Policy", "default-src 'self'; base-uri 'self'; form-action 'self'; frame-ancestors 'none'; img-src 'self' data:; style-src 'self'; script-src 'self'; connect-src 'self'")
		if a.cfg.Environment == "production" {
			response.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		next.ServeHTTP(response, request)
	})
}

func (a *App) sameOrigin(next http.Handler) http.Handler {
	appURL, _ := url.Parse(a.cfg.AppURL)
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if isUnsafe(request.Method) {
			origin := request.Header.Get("Origin")
			if origin != "" {
				parsed, err := url.Parse(origin)
				if err != nil || !strings.EqualFold(parsed.Host, appURL.Host) || parsed.Scheme != appURL.Scheme {
					writeError(response, http.StatusForbidden, "Origin permintaan tidak diizinkan.", "invalid_origin")
					return
				}
			}
		}
		next.ServeHTTP(response, request)
	})
}

func (a *App) recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				a.log.Error("panic", "error", recovered, "stack", string(debug.Stack()))
				writeError(response, http.StatusInternalServerError, "Terjadi kesalahan internal.", "internal_error")
			}
		}()
		next.ServeHTTP(response, request)
	})
}

func (a *App) requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		started := time.Now()
		requestID := randomRequestID()
		response.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(response, request)
		a.log.Info("request", "id", requestID, "method", request.Method, "path", request.URL.Path, "duration_ms", time.Since(started).Milliseconds())
	})
}

func decodeJSON(response http.ResponseWriter, request *http.Request, target any) bool {
	request.Body = http.MaxBytesReader(response, request.Body, 1<<20)
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeError(response, http.StatusBadRequest, "Isi permintaan tidak valid.", "invalid_request")
		return false
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeError(response, http.StatusBadRequest, "Isi permintaan hanya boleh memuat satu objek JSON.", "invalid_request")
		return false
	}
	return true
}

func writeJSON(response http.ResponseWriter, status int, value any) {
	response.Header().Set("Content-Type", "application/json; charset=utf-8")
	response.Header().Set("Cache-Control", "no-store")
	response.WriteHeader(status)
	_ = json.NewEncoder(response).Encode(value)
}

func writeError(response http.ResponseWriter, status int, message, code string) {
	writeJSON(response, status, map[string]any{"error": map[string]string{"message": message, "code": code}})
}

func isUnsafe(method string) bool {
	return method != http.MethodGet && method != http.MethodHead && method != http.MethodOptions
}

func randomRequestID() string {
	buffer := make([]byte, 8)
	if _, err := rand.Read(buffer); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buffer)
}
