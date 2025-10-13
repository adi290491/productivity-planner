package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/adi290491/productivity-planner/user-service/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestGenerateToken(t *testing.T) {
	j := JWTUtil{Secret: []byte("test-secret")}
	user := &models.User{
		ID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
	}

	token, err := j.GenerateToken(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected a non-empty token")
	}
}

func TestValidateToken_Valid(t *testing.T) {
	j := JWTUtil{Secret: []byte("test-secret")}
	user := &models.User{
		ID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
	}
	token, _ := j.GenerateToken(user)

	err := j.ValidateToken(token)
	if err != nil {
		t.Errorf("expected token to be valid, got error: %v", err)
	}
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	// Token created with different secret
	jGood := JWTUtil{Secret: []byte("correct-secret")}
	jBad := JWTUtil{Secret: []byte("wrong-secret")}

	user := &models.User{ID: uuid.New()}
	token, _ := jGood.GenerateToken(user)

	err := jBad.ValidateToken(token)
	if err == nil {
		t.Error("expected validation to fail due to invalid signature, got nil")
	}
}

func TestValidateToken_Expired(t *testing.T) {
	secret := []byte("test-secret")
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": "123",
		"exp":    time.Now().Add(-time.Hour).Unix(), // expired 1 hour ago
	})
	tokenStr, _ := expiredToken.SignedString(secret)

	j := JWTUtil{Secret: secret}
	err := j.ValidateToken(tokenStr)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestValidateToken_InvalidFormat(t *testing.T) {
	j := JWTUtil{Secret: []byte("test-secret")}

	tests := []struct {
		name  string
		token string
	}{
		{"empty token", ""},
		{"malformed token", "not.a.valid.jwt"},
		{"missing parts", "invalid"},
		{"too many parts", "too.many.parts.in.token"},
		{"invalid characters", "invalid$characters!"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := j.ValidateToken(test.token)
			if err == nil {
				t.Errorf("expected error for %s, got nil", test.name)
			}
		})
	}
}

func TestGenerateToken_DifferentUsers(t *testing.T) {
	j := JWTUtil{Secret: []byte("test-secret")}

	user1 := &models.User{ID: uuid.New()}
	user2 := &models.User{ID: uuid.New()}

	token1, err1 := j.GenerateToken(user1)
	token2, err2 := j.GenerateToken(user2)

	if err1 != nil || err2 != nil {
		t.Fatal("expected no errors generating tokens")
	}

	if token1 == token2 {
		t.Error("expected different tokens for different users")
	}

	// Both tokens should be valid
	if err := j.ValidateToken(token1); err != nil {
		t.Errorf("token1 should be valid: %v", err)
	}
	if err := j.ValidateToken(token2); err != nil {
		t.Errorf("token2 should be valid: %v", err)
	}
}

func TestJWTUtil_EmptySecret(t *testing.T) {
	j := JWTUtil{Secret: []byte("")}
	user := &models.User{ID: uuid.New()}

	token, err := j.GenerateToken(user)

	// Generation might succeed even with empty secret
	if err == nil && token != "" {
		// If generation succeeds, validation should also work
		err = j.ValidateToken(token)
		if err != nil {
			t.Errorf("if generation succeeds with empty secret, validation should too: %v", err)
		}
	}
}

func TestJWTUtil_RoundTrip(t *testing.T) {
	secrets := [][]byte{
		[]byte("simple-secret"),
		[]byte("complex!@#$%^&*()_+secret"),
		[]byte("very-long-secret-key-that-is-much-longer-than-typical"),
		[]byte("1"),
	}

	for i, secret := range secrets {
		t.Run(fmt.Sprintf("secret-%d", i), func(t *testing.T) {
			j := JWTUtil{Secret: secret}
			user := &models.User{
				ID:    uuid.New(),
				Email: "test@example.com",
				Name:  "Test User",
			}

			// Generate token
			token, err := j.GenerateToken(user)
			if err != nil {
				t.Errorf("Failed to generate token: %v", err)
				return
			}

			// Validate token
			err = j.ValidateToken(token)
			if err != nil {
				t.Errorf("Failed to validate token: %v", err)
			}
		})
	}
}
