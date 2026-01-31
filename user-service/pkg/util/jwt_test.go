package util

import (
	"testing"
	"time"

	"github.com/adi290491/productivity-planner/user-service/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestNewJWTUtil(t *testing.T) {
	secret := []byte("test-secret")
	jwtUtil := NewJWTUtil(secret)

	if jwtUtil == nil {
		t.Error("NewJWTUtil() returned nil")
	}
	if string(jwtUtil.Secret) != string(secret) {
		t.Error("NewJWTUtil() did not set secret correctly")
	}
}

func TestGenerateToken(t *testing.T) {
	jwtUtil := NewJWTUtil([]byte("test-secret"))

	user := &model.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}

	token, err := jwtUtil.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	if token == "" {
		t.Error("GenerateToken() returned empty token")
	}

	// Verify token structure (3 parts separated by dots)
	parts := 0
	for _, ch := range token {
		if ch == '.' {
			parts++
		}
	}
	if parts != 2 {
		t.Errorf("GenerateToken() returned invalid token format, got %d parts, want 2", parts)
	}
}

func TestValidateToken(t *testing.T) {
	jwtUtil := NewJWTUtil([]byte("test-secret"))

	user := &model.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}

	token, err := jwtUtil.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken() failed: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   token,
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "invalid token",
			token:   "invalid.token.here",
			wantErr: true,
		},
		{
			name:    "malformed token",
			token:   "notavalidtoken",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := jwtUtil.ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && claims == nil {
				t.Error("ValidateToken() returned nil claims for valid token")
			}
		})
	}
}

func TestValidateToken_WithWrongSecret(t *testing.T) {
	jwtUtil1 := NewJWTUtil([]byte("secret1"))
	jwtUtil2 := NewJWTUtil([]byte("secret2"))

	user := &model.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}

	token, err := jwtUtil1.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken() failed: %v", err)
	}

	// Try to validate with different secret
	_, err = jwtUtil2.ValidateToken(token)
	if err == nil {
		t.Error("ValidateToken() should fail with wrong secret")
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	secret := []byte("test-secret")
	jwtUtil := NewJWTUtil(secret)

	// Create expired token manually
	claims := jwt.MapClaims{
		"userId": uuid.New().String(),
		"email":  "test@example.com",
		"exp":    time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
		"iat":    time.Now().Add(-2 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	_, err = jwtUtil.ValidateToken(tokenString)
	if err == nil {
		t.Error("ValidateToken() should fail for expired token")
	}
}

func TestJWTRoundTrip(t *testing.T) {
	jwtUtil := NewJWTUtil([]byte("test-secret"))

	users := []*model.User{
		{
			ID:    uuid.New(),
			Email: "user1@example.com",
			Name:  "User One",
		},
		{
			ID:    uuid.New(),
			Email: "user2@example.com",
			Name:  "User Two",
		},
	}

	for _, user := range users {
		t.Run(user.Email, func(t *testing.T) {
			// Generate token
			token, err := jwtUtil.GenerateToken(user)
			if err != nil {
				t.Fatalf("GenerateToken() failed: %v", err)
			}

			// Validate token
			claims, err := jwtUtil.ValidateToken(token)
			if err != nil {
				t.Fatalf("ValidateToken() failed: %v", err)
			}

			// Verify claims
			if (*claims)["userId"] != user.ID.String() {
				t.Errorf("userId mismatch: got %v, want %v", (*claims)["userId"], user.ID.String())
			}
			if (*claims)["email"] != user.Email {
				t.Errorf("email mismatch: got %v, want %v", (*claims)["email"], user.Email)
			}
		})
	}
}
