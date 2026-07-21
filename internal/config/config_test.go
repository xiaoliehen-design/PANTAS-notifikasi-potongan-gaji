package config

import "testing"

func TestValidAdminUsername(t *testing.T) {
	tests := map[string]bool{
		"admin.pantas":  true,
		"admin_utama-2": true,
		"Admin.Pantas":  false,
		"2admin":        false,
		"admin pantas":  false,
		"ad":            false,
	}
	for value, expected := range tests {
		if actual := validAdminUsername(value); actual != expected {
			t.Errorf("validAdminUsername(%q) = %v, want %v", value, actual, expected)
		}
	}
}

func TestValidAdminPassword(t *testing.T) {
	if !validAdminPassword("Password-Awal-2026!", "admin.pantas") {
		t.Fatal("strong password rejected")
	}
	for _, value := range []string{"terlalupendek", "admin.pantas-Aman-2026!", "semuahurufkecilpanjang"} {
		if validAdminPassword(value, "admin.pantas") {
			t.Errorf("weak password %q accepted", value)
		}
	}
}

func TestValidHTTPURL(t *testing.T) {
	for _, value := range []string{"https://api.example.go.id/otp", "http://localhost:8080/test"} {
		if !validHTTPURL(value) {
			t.Errorf("validHTTPURL(%q) = false", value)
		}
	}
	for _, value := range []string{"", "api.example.go.id", "ftp://example.go.id/file"} {
		if validHTTPURL(value) {
			t.Errorf("validHTTPURL(%q) = true", value)
		}
	}
	if !isHTTPSURL("https://provider.example.go.id/otp") || isHTTPSURL("http://localhost:8080/otp") {
		t.Fatal("isHTTPSURL did not enforce https")
	}
}

func TestEnvMailAddressRemovesAccidentalOuterQuotes(t *testing.T) {
	t.Setenv("EMAIL_FROM", `"PANTAS <pantas@example.com>"`)
	if actual := envMailAddress("EMAIL_FROM"); actual != "PANTAS <pantas@example.com>" {
		t.Fatalf("envMailAddress() = %q", actual)
	}
}

func TestSMTPTLSModeAcceptsLegacyRenderKey(t *testing.T) {
	t.Setenv("SMTP_TLS_MODE", "")
	t.Setenv("SMTP_TLS", "STARTTLS")
	if actual := smtpTLSMode(); actual != "starttls" {
		t.Fatalf("smtpTLSMode() = %q, want starttls", actual)
	}
	t.Setenv("SMTP_TLS_MODE", "implicit")
	if actual := smtpTLSMode(); actual != "implicit" {
		t.Fatalf("smtpTLSMode() = %q, want implicit", actual)
	}
}

func TestLoadAcceptsBrevoHTTPSProvider(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("APP_URL", "http://localhost:10000")
	t.Setenv("APP_SECRET", "01234567890123456789012345678901")
	t.Setenv("DATABASE_URL", "postgresql://example.invalid/pantas")
	t.Setenv("SUPABASE_URL", "https://project.supabase.co")
	t.Setenv("SUPABASE_SERVICE_ROLE_KEY", "test-service-role-key")
	t.Setenv("BOOTSTRAP_ADMIN_USERNAME", "")
	t.Setenv("BOOTSTRAP_ADMIN_PASSWORD", "")
	t.Setenv("BOOTSTRAP_ADMIN_NAME", "")
	t.Setenv("BOOTSTRAP_TREASURY_USERNAME", "")
	t.Setenv("BOOTSTRAP_TREASURY_PASSWORD", "")
	t.Setenv("BOOTSTRAP_TREASURY_NAME", "")
	t.Setenv("EMAIL_PROVIDER", "brevo")
	t.Setenv("EMAIL_FROM", "PANTAS <noreply@example.go.id>")
	t.Setenv("BREVO_API_KEY", "brevo-test-key")
	t.Setenv("BREVO_API_URL", "https://api.brevo.com/v3/smtp/email")
	t.Setenv("RESEND_API_KEY", "")
	t.Setenv("SMTP_HOST", "")
	t.Setenv("SMTP_USERNAME", "")
	t.Setenv("SMTP_PASSWORD", "")
	t.Setenv("PHONE_PROVIDER", "auto")
	t.Setenv("PHONE_OTP_WEBHOOK_URL", "")
	t.Setenv("PHONE_OTP_WEBHOOK_TOKEN", "")
	t.Setenv("TWILIO_ACCOUNT_SID", "")
	t.Setenv("TWILIO_AUTH_TOKEN", "")
	t.Setenv("TWILIO_API_KEY", "")
	t.Setenv("TWILIO_API_SECRET", "")
	t.Setenv("TWILIO_MESSAGING_SERVICE_SID", "")
	t.Setenv("TWILIO_FROM_NUMBER", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.EmailProvider != "brevo" || cfg.BrevoAPIKey == "" {
		t.Fatalf("Brevo configuration was not loaded: %#v", cfg)
	}
}
