package utils

import "testing"

func Test_VerifyPassword(t *testing.T) {

	password := "1234"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		err      bool
	}{
		{
			name:     "matching hash and password",
			password: "1234",
			hash:     hash,
			err:      false,
		},
		{
			name:     "non-matching hash and password",
			password: "1233",
			hash:     hash,
			err:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := VerifyPassword(test.password, test.hash)
			if test.err && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !test.err && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
