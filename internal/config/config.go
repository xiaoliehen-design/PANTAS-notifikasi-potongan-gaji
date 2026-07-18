package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Environment           string
	Port                  string
	AppURL                string
	DatabaseURL           string
	AppSecret             string
	TrustProxy            bool
	CookieSecure          bool
	SessionTTL            time.Duration
	SessionIdleTTL        time.Duration
	SupabaseURL           string
	SupabaseServiceKey    string
	SupabaseStorageBucket string
	ResendAPIKey          string
	EmailFrom             string
	PhoneWebhookURL       string
	PhoneWebhookToken     string
	BootstrapAdminNIP     string
	BootstrapAdminName    string
	MaxExcelBytes         int64
	MaxDocumentBytes      int64
	WorkerInterval        time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		Environment:           env("APP_ENV", "production"),
		Port:                  env("PORT", "10000"),
		AppURL:                strings.TrimRight(env("APP_URL", "http://localhost:10000"), "/"),
		DatabaseURL:           strings.TrimSpace(os.Getenv("DATABASE_URL")),
		AppSecret:             strings.TrimSpace(os.Getenv("APP_SECRET")),
		TrustProxy:            envBool("TRUST_PROXY", true),
		CookieSecure:          envBool("COOKIE_SECURE", true),
		SessionTTL:            envDuration("SESSION_TTL", 12*time.Hour),
		SessionIdleTTL:        envDuration("SESSION_IDLE_TTL", 2*time.Hour),
		SupabaseURL:           strings.TrimRight(strings.TrimSpace(os.Getenv("SUPABASE_URL")), "/"),
		SupabaseServiceKey:    strings.TrimSpace(os.Getenv("SUPABASE_SERVICE_ROLE_KEY")),
		SupabaseStorageBucket: env("SUPABASE_STORAGE_BUCKET", "pantas-appeals"),
		ResendAPIKey:          strings.TrimSpace(os.Getenv("RESEND_API_KEY")),
		EmailFrom:             strings.TrimSpace(os.Getenv("EMAIL_FROM")),
		PhoneWebhookURL:       strings.TrimSpace(os.Getenv("PHONE_OTP_WEBHOOK_URL")),
		PhoneWebhookToken:     strings.TrimSpace(os.Getenv("PHONE_OTP_WEBHOOK_TOKEN")),
		BootstrapAdminNIP:     strings.TrimSpace(os.Getenv("BOOTSTRAP_ADMIN_NIP")),
		BootstrapAdminName:    strings.TrimSpace(os.Getenv("BOOTSTRAP_ADMIN_NAME")),
		MaxExcelBytes:         envInt64("MAX_EXCEL_BYTES", 20<<20),
		MaxDocumentBytes:      envInt64("MAX_DOCUMENT_BYTES", 5<<20),
		WorkerInterval:        envDuration("WORKER_INTERVAL", 5*time.Second),
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
	if cfg.BootstrapAdminNIP != "" && (len(cfg.BootstrapAdminNIP) != 18 || !digitsOnly(cfg.BootstrapAdminNIP)) {
		problems = append(problems, "BOOTSTRAP_ADMIN_NIP wajib tepat 18 digit")
	}
	if (cfg.BootstrapAdminNIP == "") != (cfg.BootstrapAdminName == "") {
		problems = append(problems, "BOOTSTRAP_ADMIN_NIP dan BOOTSTRAP_ADMIN_NAME harus diisi bersamaan")
	}
	if cfg.SupabaseURL == "" || cfg.SupabaseServiceKey == "" {
		problems = append(problems, "SUPABASE_URL dan SUPABASE_SERVICE_ROLE_KEY wajib diisi untuk dokumen banding")
	}
	if len(problems) > 0 {
		return Config{}, errors.New(strings.Join(problems, "; "))
	}
	return cfg, nil
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

func digitsOnly(value string) bool {
	for _, char := range value {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func (c Config) String() string {
	return fmt.Sprintf("env=%s addr=%s app_url=%s", c.Environment, c.Address(), c.AppURL)
}
