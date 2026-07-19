package config

import (
	"errors"
	"fmt"
	mailaddress "net/mail"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Config struct {
	Environment            string
	Port                   string
	AppURL                 string
	DatabaseURL            string
	AppSecret              string
	TrustProxy             bool
	CookieSecure           bool
	SessionTTL             time.Duration
	SessionIdleTTL         time.Duration
	SupabaseURL            string
	SupabaseServiceKey     string
	SupabaseStorageBucket  string
	ResendAPIKey           string
	ResendAPIURL           string
	EmailFrom              string
	PhoneWebhookURL        string
	PhoneWebhookToken      string
	BootstrapAdminUsername string
	BootstrapAdminPassword string
	BootstrapAdminName     string
	MaxExcelBytes          int64
	MaxDocumentBytes       int64
	WorkerInterval         time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		Environment:            env("APP_ENV", "production"),
		Port:                   env("PORT", "10000"),
		AppURL:                 strings.TrimRight(env("APP_URL", "http://localhost:10000"), "/"),
		DatabaseURL:            strings.TrimSpace(os.Getenv("DATABASE_URL")),
		AppSecret:              strings.TrimSpace(os.Getenv("APP_SECRET")),
		TrustProxy:             envBool("TRUST_PROXY", true),
		CookieSecure:           envBool("COOKIE_SECURE", true),
		SessionTTL:             envDuration("SESSION_TTL", 12*time.Hour),
		SessionIdleTTL:         envDuration("SESSION_IDLE_TTL", 2*time.Hour),
		SupabaseURL:            strings.TrimRight(strings.TrimSpace(os.Getenv("SUPABASE_URL")), "/"),
		SupabaseServiceKey:     strings.TrimSpace(os.Getenv("SUPABASE_SERVICE_ROLE_KEY")),
		SupabaseStorageBucket:  env("SUPABASE_STORAGE_BUCKET", "pantas-appeals"),
		ResendAPIKey:           strings.TrimSpace(os.Getenv("RESEND_API_KEY")),
		ResendAPIURL:           strings.TrimRight(env("RESEND_API_URL", "https://api.resend.com/emails"), "/"),
		EmailFrom:              strings.TrimSpace(os.Getenv("EMAIL_FROM")),
		PhoneWebhookURL:        strings.TrimSpace(os.Getenv("PHONE_OTP_WEBHOOK_URL")),
		PhoneWebhookToken:      strings.TrimSpace(os.Getenv("PHONE_OTP_WEBHOOK_TOKEN")),
		BootstrapAdminUsername: strings.ToLower(strings.TrimSpace(os.Getenv("BOOTSTRAP_ADMIN_USERNAME"))),
		BootstrapAdminPassword: os.Getenv("BOOTSTRAP_ADMIN_PASSWORD"),
		BootstrapAdminName:     strings.TrimSpace(os.Getenv("BOOTSTRAP_ADMIN_NAME")),
		MaxExcelBytes:          envInt64("MAX_EXCEL_BYTES", 20<<20),
		MaxDocumentBytes:       envInt64("MAX_DOCUMENT_BYTES", 5<<20),
		WorkerInterval:         envDuration("WORKER_INTERVAL", 5*time.Second),
	}

	var problems []string
	if cfg.DatabaseURL == "" {
		problems = append(problems, "DATABASE_URL wajib diisi")
	}
	if len(cfg.AppSecret) < 32 {
		problems = append(problems, "APP_SECRET wajib berisi minimal 32 karakter acak")
	}
	if parsed, err := url.Parse(cfg.AppURL); err != nil || parsed.Host == "" {
		problems = append(problems, "APP_URL tidak valid")
	} else if cfg.Environment == "production" && parsed.Scheme != "https" {
		problems = append(problems, "APP_URL production wajib menggunakan https")
	}
	adminConfigured := cfg.BootstrapAdminUsername != "" || cfg.BootstrapAdminPassword != "" || cfg.BootstrapAdminName != ""
	if cfg.Environment == "production" && !adminConfigured {
		problems = append(problems, "BOOTSTRAP_ADMIN_USERNAME, BOOTSTRAP_ADMIN_PASSWORD, dan BOOTSTRAP_ADMIN_NAME wajib diisi")
	}
	if adminConfigured && (cfg.BootstrapAdminUsername == "" || cfg.BootstrapAdminPassword == "" || cfg.BootstrapAdminName == "") {
		problems = append(problems, "konfigurasi bootstrap admin harus diisi lengkap")
	}
	if cfg.BootstrapAdminUsername != "" && !validAdminUsername(cfg.BootstrapAdminUsername) {
		problems = append(problems, "BOOTSTRAP_ADMIN_USERNAME harus 3-64 karakter, diawali huruf, dan hanya memakai huruf kecil, angka, titik, garis bawah, atau tanda hubung")
	}
	if cfg.BootstrapAdminPassword != "" && !validAdminPassword(cfg.BootstrapAdminPassword, cfg.BootstrapAdminUsername) {
		problems = append(problems, "BOOTSTRAP_ADMIN_PASSWORD harus 12-128 karakter, tidak memuat username, dan memakai sedikitnya tiga jenis karakter")
	}
	if cfg.SupabaseURL == "" || cfg.SupabaseServiceKey == "" {
		problems = append(problems, "SUPABASE_URL dan SUPABASE_SERVICE_ROLE_KEY wajib diisi untuk dokumen banding")
	}
	emailPartiallyConfigured := cfg.ResendAPIKey != "" || cfg.EmailFrom != ""
	if emailPartiallyConfigured && (cfg.ResendAPIKey == "" || cfg.EmailFrom == "") {
		problems = append(problems, "RESEND_API_KEY dan EMAIL_FROM harus diisi bersamaan agar OTP email dapat dikirim")
	}
	if cfg.EmailFrom != "" {
		if _, err := mailaddress.ParseAddress(cfg.EmailFrom); err != nil {
			problems = append(problems, "EMAIL_FROM tidak valid")
		}
	}
	if cfg.ResendAPIKey != "" && !validHTTPURL(cfg.ResendAPIURL) {
		problems = append(problems, "RESEND_API_URL tidak valid")
	} else if cfg.Environment == "production" && cfg.ResendAPIKey != "" && !isHTTPSURL(cfg.ResendAPIURL) {
		problems = append(problems, "RESEND_API_URL production wajib menggunakan https")
	}
	if cfg.PhoneWebhookURL != "" && !validHTTPURL(cfg.PhoneWebhookURL) {
		problems = append(problems, "PHONE_OTP_WEBHOOK_URL tidak valid")
	} else if cfg.Environment == "production" && cfg.PhoneWebhookURL != "" && !isHTTPSURL(cfg.PhoneWebhookURL) {
		problems = append(problems, "PHONE_OTP_WEBHOOK_URL production wajib menggunakan https")
	}
	if cfg.PhoneWebhookToken != "" && cfg.PhoneWebhookURL == "" {
		problems = append(problems, "PHONE_OTP_WEBHOOK_URL wajib diisi jika PHONE_OTP_WEBHOOK_TOKEN digunakan")
	}
	if len(problems) > 0 {
		return Config{}, errors.New(strings.Join(problems, "; "))
	}
	return cfg, nil
}

func validHTTPURL(value string) bool {
	parsed, err := url.Parse(value)
	return err == nil && parsed.Host != "" && (parsed.Scheme == "https" || parsed.Scheme == "http")
}

func isHTTPSURL(value string) bool {
	parsed, err := url.Parse(value)
	return err == nil && parsed.Host != "" && parsed.Scheme == "https"
}

func (c Config) Address() string { return ":" + c.Port }

func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func envDuration(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func envInt64(key string, fallback int64) int64 {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func validAdminUsername(value string) bool {
	if len(value) < 3 || len(value) > 64 || value[0] < 'a' || value[0] > 'z' {
		return false
	}
	for _, char := range value {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '.' || char == '_' || char == '-') {
			return false
		}
	}
	return true
}

func validAdminPassword(value, username string) bool {
	if len(value) < 12 || len(value) > 128 || (username != "" && strings.Contains(strings.ToLower(value), username)) {
		return false
	}
	var upper, lower, digit, symbol bool
	for _, char := range value {
		switch {
		case unicode.IsUpper(char):
			upper = true
		case unicode.IsLower(char):
			lower = true
		case unicode.IsDigit(char):
			digit = true
		default:
			symbol = true
		}
	}
	classes := 0
	for _, present := range []bool{upper, lower, digit, symbol} {
		if present {
			classes++
		}
	}
	return classes >= 3
}

func (c Config) String() string {
	return fmt.Sprintf("env=%s addr=%s app_url=%s", c.Environment, c.Address(), c.AppURL)
}
