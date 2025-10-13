package utils

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	t.Run("generates hash for valid password", func(t *testing.T) {
		password := "testpassword123"
		hash, err := HashPassword(password)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if hash == "" {
			t.Error("Expected non-empty hash")
		}
		if hash == password {
			t.Error("Hash should not equal original password")
		}
		if !strings.HasPrefix(hash, "$2a$") {
			t.Error("Expected bcrypt hash format")
		}
	})

	t.Run("generates different hashes for same password", func(t *testing.T) {
		password := "samepassword"
		hash1, err1 := HashPassword(password)
		hash2, err2 := HashPassword(password)

		if err1 != nil || err2 != nil {
			t.Error("Expected no errors generating hashes")
		}
		if hash1 == hash2 {
			t.Error("Expected different hashes for same password (salt should be different)")
		}
	})

	t.Run("handles empty password", func(t *testing.T) {
		password := ""
		hash, err := HashPassword(password)

		if err != nil {
			t.Errorf("Expected no error for empty password, got %v", err)
		}
		if hash == "" {
			t.Error("Expected non-empty hash even for empty password")
		}
	})
}

func TestVerifyPassword(t *testing.T) {
	password := "1234"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	tests := []struct {
		name        string
		password    string
		hash        string
		expectError bool
	}{
		{
			name:        "matching hash and password",
			password:    "1234",
			hash:        hash,
			expectError: false,
		},
		{
			name:        "non-matching hash and password",
			password:    "1233",
			hash:        hash,
			expectError: true,
		},
		{
			name:        "empty password with valid hash",
			password:    "",
			hash:        hash,
			expectError: true,
		},
		{
			name:        "valid password with empty hash",
			password:    "1234",
			hash:        "",
			expectError: true,
		},
		{
			name:        "valid password with invalid hash format",
			password:    "1234",
			hash:        "invalidhash",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := VerifyPassword(test.password, test.hash)
			if test.expectError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !test.expectError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestPasswordHashingRoundTrip(t *testing.T) {
	passwords := []string{
		"simple",
		"complex!@#$%^&*()_+",
		"withNumbers123",
		"CaseSensitive",
		"   spaces   ",
		"unicode🔐password",
	}

	for _, password := range passwords {
		t.Run("password: "+password, func(t *testing.T) {
			// Hash the password
			hash, err := HashPassword(password)
			if err != nil {
				t.Errorf("Failed to hash password: %v", err)
			}

			// Verify the password matches the hash
			err = VerifyPassword(password, hash)
			if err != nil {
				t.Errorf("Password verification failed: %v", err)
			}

			// Verify wrong password doesn't match
			wrongPassword := password + "wrong"
			err = VerifyPassword(wrongPassword, hash)
			if err == nil {
				t.Error("Wrong password should not verify successfully")
			}
		})
	}
}
