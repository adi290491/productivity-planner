package util

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "testpassword123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false, // bcrypt can hash empty strings
		},
		{
			name:     "password at 72 byte limit",
			password: strings.Repeat("a", 72),
			wantErr:  false,
		},
		{
			name:     "password exceeding 72 bytes",
			password: strings.Repeat("a", 100),
			wantErr:  true, // bcrypt will fail with passwords > 72 bytes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if hash == "" {
					t.Error("HashPassword() returned empty hash")
				}
				if hash == tt.password {
					t.Error("HashPassword() returned unhashed password")
				}
				if !strings.HasPrefix(hash, "$2a$") {
					t.Error("HashPassword() did not return bcrypt hash")
				}
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "correct password",
			password: password,
			hash:     hash,
			wantErr:  false,
		},
		{
			name:     "incorrect password",
			password: "wrongpassword",
			hash:     hash,
			wantErr:  true,
		},
		{
			name:     "empty password",
			password: "",
			hash:     hash,
			wantErr:  true,
		},
		{
			name:     "invalid hash",
			password: password,
			hash:     "invalid_hash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyPassword(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPasswordHashingRoundTrip(t *testing.T) {
	passwords := []string{
		"simple",
		"complex!@#$%^&*()",
		"with spaces",
		"CaseSensitive123",
		"unicode🔐password",
	}

	for _, password := range passwords {
		t.Run(password, func(t *testing.T) {
			// Hash
			hash, err := HashPassword(password)
			if err != nil {
				t.Fatalf("HashPassword() failed: %v", err)
			}

			// Verify correct password
			if err := VerifyPassword(password, hash); err != nil {
				t.Errorf("VerifyPassword() failed for correct password: %v", err)
			}

			// Verify wrong password fails
			if err := VerifyPassword(password+"wrong", hash); err == nil {
				t.Error("VerifyPassword() should fail for wrong password")
			}
		})
	}
}
