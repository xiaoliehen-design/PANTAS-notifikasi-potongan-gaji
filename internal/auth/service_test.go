package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bcpriok/pantas/internal/config"
)

func TestValidatePassword(t *testing.T) {
	nip := "199001012010011001"
	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{name: "strong", password: "Pantas-Aman-2026", valid: true},
		{name: "too short", password: "Ab1!", valid: false},
		{name: "contains NIP", password: "Aman!" + nip, valid: false},
		{name: "too few classes", password: "semuahurufkecil", valid: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidatePassword(test.password, nip)
			if (err == nil) != test.valid {
				t.Fatalf("ValidatePassword() error = %v, valid want %v", err, test.valid)
			}
		})
	}
}

func TestNormalizePhoneIndonesia(t *testing.T) {
	for input, expected := range map[string]string{
		"0812-3456-7890": "+6281234567890",
		"81234567890":    "+6281234567890",
		"+6281234567890": "+6281234567890",
	} {
		actual, err := normalizePhone(input)
		if err != nil || actual != expected {
			t.Fatalf("normalizePhone(%q) = %q, %v; want %q", input, actual, err, expected)
		}
	}
}

func TestValidLoginIdentifier(t *testing.T) {
	tests := map[string]bool{
		"199001012010011001": true,
		"admin.pantas":       true,
		"admin_pantas-2":     true,
		"Admin.Pantas":       false,
		"12":                 false,
		"admin pantas":       false,
		"123456789012345678": true,
	}
	for input, expected := range tests {
		if actual := validLoginIdentifier(input); actual != expected {
			t.Errorf("validLoginIdentifier(%q) = %v, want %v", input, actual, expected)
		}
	}
}

func TestPrincipalLoginIdentifier(t *testing.T) {
	user := Principal{AccountType: "user", NIP: "199001012010011001"}
	admin := Principal{AccountType: "admin", Username: "admin.pantas"}
	if user.LoginIdentifier() != user.NIP {
		t.Fatalf("user identifier = %q", user.LoginIdentifier())
	}
	if admin.LoginIdentifier() != admin.Username {
		t.Fatalf("admin identifier = %q", admin.LoginIdentifier())
	}
}

func TestTabSessionProofRequiresMatchingHeaderAndCookie(t *testing.T) {
	service := &Service{cfg: config.Config{AppSecret: "01234567890123456789012345678901"}}
	token, _, err := randomToken(32)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest("GET", "/api/auth/me", nil)
	request.Header.Set(TabTokenHeader, token)
	request.AddCookie(&http.Cookie{Name: TabProofCookieName, Value: service.tabProof(token)})
	if !service.verifyTabSession(request) {
		t.Fatal("matching tab session was rejected")
	}

	request.Header.Set(TabTokenHeader, token+"tampered")
	if service.verifyTabSession(request) {
		t.Fatal("tampered tab token was accepted")
	}
}

func TestAuthenticationCookiesAreSessionScopedAndTabBound(t *testing.T) {
	service := &Service{cfg: config.Config{AppSecret: "01234567890123456789012345678901", CookieSecure: true}}
	sessionToken, _, _ := randomToken(32)
	csrfToken, _, _ := randomToken(32)
	tabToken, _, _ := randomToken(32)
	recorder := httptest.NewRecorder()
	service.SetCookies(recorder, LoginResult{SessionToken: sessionToken, CSRFToken: csrfToken, TabToken: tabToken})

	cookies := recorder.Result().Cookies()
	if len(cookies) != 3 {
		t.Fatalf("cookie count = %d, want 3", len(cookies))
	}
	values := make(map[string]*http.Cookie, len(cookies))
	for _, cookie := range cookies {
		values[cookie.Name] = cookie
		if cookie.MaxAge != 0 || !cookie.Secure {
			t.Fatalf("cookie %s is not a secure session cookie: %#v", cookie.Name, cookie)
		}
	}
	if values[TabProofCookieName] == nil || values[TabProofCookieName].Value != service.tabProof(tabToken) {
		t.Fatal("tab proof cookie is missing or invalid")
	}
	if !values[SessionCookieName].HttpOnly || values[CSRFCookieName].HttpOnly || !values[TabProofCookieName].HttpOnly {
		t.Fatal("cookie HttpOnly flags are incorrect")
	}
}

func TestProviderDeliveryErrorMatchesPublicSentinel(t *testing.T) {
	err := providerDeliveryError(errors.New("koneksi SMTP gagal: dial tcp: i/o timeout"))
	if !errors.Is(err, ErrDeliveryFailed) {
		t.Fatal("provider error does not match ErrDeliveryFailed")
	}
	if err.Error() == ErrDeliveryFailed.Error() {
		t.Fatal("provider error did not provide an actionable public message")
	}
}
