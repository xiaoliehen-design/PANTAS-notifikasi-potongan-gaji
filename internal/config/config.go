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
	EmailProvider          string
	ResendAPIKey           string
	ResendAPIURL           string
	EmailFrom              string
	SMTPHost               string
	SMTPPort               string
	SMTPUsername           string
	SMTPPassword           string
	SMTPTLSMode            string
	PhoneProvider          string
	PhoneWebhookURL        string
	PhoneWebhookToken      string
	TwilioAccountSID       string
	TwilioAuthToken        string
	TwilioAPIKey           string
	TwilioAPISecret        string
	TwilioMessagingSID     string
	TwilioFrom             string
	TwilioAPIBaseURL       string
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
		SessionIdleTTL:         envDuration("SESSION_IDLE_TTL", 30*time.Minute),
		SupabaseURL:            strings.TrimRight(strings.TrimSpace(os.Getenv("SUPABASE_URL")), "/"),
		SupabaseServiceKey:     strings.TrimSpace(os.Getenv("SUPABASE_SERVICE_ROLE_KEY")),
		SupabaseStorageBucket:  env("SUPABASE_STORAGE_BUCKET", "pantas-appeals"),
		EmailProvider:          strings.ToLower(env("EMAIL_PROVIDER", "auto")),
		ResendAPIKey:           strings.TrimSpace(os.Getenv("RESEND_API_KEY")),
		ResendAPIURL:           strings.TrimRight(env("RESEND_API_URL", "https://api.resend.com/emails"), "/"),
		EmailFrom:              envMailAddress("EMAIL_FROM"),
		SMTPHost:               strings.TrimSpace(os.Getenv("SMTP_HOST")),
		SMTPPort:               env("SMTP_PORT", "587"),
		SMTPUsername:           strings.TrimSpace(os.Getenv("SMTP_USERNAME")),
		SMTPPassword:           strings.TrimSpace(os.Getenv("SMTP_PASSWORD")),
		SMTPTLSMode:            smtpTLSMode(),
		PhoneProvider:          strings.ToLower(env("PHONE_PROVIDER", "auto")),
		PhoneWebhookURL:        strings.TrimSpace(os.Getenv("PHONE_OTP_WEBHOOK_URL")),
		PhoneWebhookToken:      strings.TrimSpace(os.Getenv("PHONE_OTP_WEBHOOK_TOKEN")),
		TwilioAccountSID:       strings.TrimSpace(os.Getenv("TWILIO_ACCOUNT_SID")),
		TwilioAuthToken:        strings.TrimSpace(os.Getenv("TWILIO_AUTH_TOKEN")),
		TwilioAPIKey:           strings.TrimSpace(os.Getenv("TWILIO_API_KEY")),
		TwilioAPISecret:        strings.TrimSpace(os.Getenv("TWILIO_API_SECRET")),
		TwilioMessagingSID:     strings.TrimSpace(os.Getenv("TWILIO_MESSAGING_SERVICE_SID")),
		TwilioFrom:             strings.TrimSpace(os.Getenv("TWILIO_FROM_NUMBER")),
		TwilioAPIBaseURL:       strings.TrimRight(env("TWILIO_API_BASE_URL", "https://api.twilio.com/2010-04-01"), "/"),
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
	if cfg.EmailProvider != "auto" && cfg.EmailProvider != "resend" && cfg.EmailProvider != "smtp" {
		problems = append(problems, "EMAIL_PROVIDER harus berisi auto, resend, atau smtp")
	}
	resendConfigured := cfg.ResendAPIKey != ""
	smtpConfigured := cfg.SMTPHost != "" || cfg.SMTPUsername != "" || cfg.SMTPPassword != ""
	if cfg.EmailProvider == "resend" || (cfg.EmailProvider == "auto" && resendConfigured && !smtpConfigured) {
		if cfg.ResendAPIKey == "" || cfg.EmailFrom == "" {
			problems = append(problems, "RESEND_API_KEY dan EMAIL_FROM wajib diisi untuk EMAIL_PROVIDER=resend")
		}
	}
	if cfg.EmailProvider == "smtp" || (cfg.EmailProvider == "auto" && smtpConfigured) {
		if cfg.SMTPHost == "" || cfg.SMTPPort == "" || cfg.SMTPUsername == "" || cfg.SMTPPassword == "" || cfg.EmailFrom == "" {
			problems = append(problems, "SMTP_HOST, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD, dan EMAIL_FROM wajib diisi untuk EMAIL_PROVIDER=smtp")
		}
		if port, err := strconv.Atoi(cfg.SMTPPort); err != nil || port < 1 || port > 65535 {
			problems = append(problems, "SMTP_PORT tidak valid")
		}
		if cfg.SMTPTLSMode != "starttls" && cfg.SMTPTLSMode != "implicit" {
			problems = append(problems, "SMTP_TLS_MODE harus berisi starttls atau implicit")
		}
	}
	if cfg.EmailProvider == "auto" && resendConfigured && smtpConfigured {
		problems = append(problems, "EMAIL_PROVIDER wajib dipilih jika konfigurasi Resend dan SMTP sama-sama terisi")
	}
	if cfg.EmailFrom != "" {
		if _, err := mailaddress.ParseAddress(cfg.EmailFrom); err != nil {
			problems = append(problems, "EMAIL_FROM tidak valid")
		}
	}
	if (cfg.EmailProvider == "resend" || resendConfigured) && !validHTTPURL(cfg.ResendAPIURL) {
		problems = append(problems, "RESEND_API_URL tidak valid")
	} else if cfg.Environment == "production" && (cfg.EmailProvider == "resend" || resendConfigured) && !isHTTPSURL(cfg.ResendAPIURL) {
		problems = append(problems, "RESEND_API_URL production wajib menggunakan https")
	}
	if cfg.PhoneProvider != "auto" && cfg.PhoneProvider != "webhook" && cfg.PhoneProvider != "twilio" {
		problems = append(problems, "PHONE_PROVIDER harus berisi auto, webhook, atau twilio")
	}
	twilioConfigured := cfg.TwilioAccountSID != "" || cfg.TwilioAuthToken != "" || cfg.TwilioAPIKey != "" || cfg.TwilioAPISecret != "" || cfg.TwilioMessagingSID != "" || cfg.TwilioFrom != ""
	webhookConfigured := cfg.PhoneWebhookURL != "" || cfg.PhoneWebhookToken != ""
	if cfg.PhoneProvider == "webhook" || (cfg.PhoneProvider == "auto" && webhookConfigured && !twilioConfigured) {
		if cfg.PhoneWebhookURL == "" {
			problems = append(problems, "PHONE_OTP_WEBHOOK_URL wajib diisi untuk PHONE_PROVIDER=webhook")
		}
	}
	if cfg.PhoneWebhookURL != "" && !validHTTPURL(cfg.PhoneWebhookURL) {
		problems = append(problems, "PHONE_OTP_WEBHOOK_URL tidak valid")
	} else if cfg.Environment == "production" && cfg.PhoneWebhookURL != "" && !isHTTPSURL(cfg.PhoneWebhookURL) {
		problems = append(problems, "PHONE_OTP_WEBHOOK_URL production wajib menggunakan https")
	}
	if cfg.PhoneWebhookToken != "" && cfg.PhoneWebhookURL == "" {
		problems = append(problems, "PHONE_OTP_WEBHOOK_URL wajib diisi jika PHONE_OTP_WEBHOOK_TOKEN digunakan")
	}
	if cfg.PhoneProvider == "twilio" || (cfg.PhoneProvider == "auto" && twilioConfigured) {
		credentialOK := (cfg.TwilioAPIKey != "" && cfg.TwilioAPISecret != "") || cfg.TwilioAuthToken != ""
		if cfg.TwilioAccountSID == "" || !credentialOK || (cfg.TwilioMessagingSID == "" && cfg.TwilioFrom == "") {
			problems = append(problems, "Twilio memerlukan TWILIO_ACCOUNT_SID, kredensial API/Auth Token, serta TWILIO_MESSAGING_SERVICE_SID atau TWILIO_FROM_NUMBER")
		}
		if (cfg.TwilioAPIKey == "") != (cfg.TwilioAPISecret == "") {
			problems = append(problems, "TWILIO_API_KEY dan TWILIO_API_SECRET harus diisi bersamaan")
		}
		if !validHTTPURL(cfg.TwilioAPIBaseURL) {
			problems = append(problems, "TWILIO_API_BASE_URL tidak valid")
		} else if cfg.Environment == "production" && !isHTTPSURL(cfg.TwilioAPIBaseURL) {
			problems = append(problems, "TWILIO_API_BASE_URL production wajib menggunakan https")
		}
	}
	if cfg.PhoneProvider == "auto" && webhookConfigured && twilioConfigured {
		problems = append(problems, "PHONE_PROVIDER wajib dipilih jika konfigurasi webhook dan Twilio sama-sama terisi")
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

// envMailAddress accepts the value exactly as Render stores it, while also
// tolerating a common dashboard copy/paste mistake: quoting the whole value as
// "PANTAS <address@example.com>". The quote characters are not part of the
// RFC mailbox and would otherwise be sent to the SMTP server as the address.
func envMailAddress(key string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if len(value) >= 2 {
		first, last := value[0], value[len(value)-1]
		if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
			value = strings.TrimSpace(value[1 : len(value)-1])
		}
	}
	return value
}

// smtpTLSMode keeps compatibility with the legacy SMTP_TLS key found in an
// earlier Render setup. SMTP_TLS_MODE remains the canonical setting.
func smtpTLSMode() string {
	if value := strings.TrimSpace(os.Getenv("SMTP_TLS_MODE")); value != "" {
		return strings.ToLower(value)
	}
	return strings.ToLower(env("SMTP_TLS", "starttls"))
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
