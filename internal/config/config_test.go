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
