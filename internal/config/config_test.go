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
