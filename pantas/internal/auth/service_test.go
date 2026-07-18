package auth

import "testing"

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
