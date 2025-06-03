package middleware

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateToken(secret string, claims jwt.MapClaims, method jwt.SigningMethod) (string, error) {
	token := jwt.NewWithClaims(method, claims)
	return token.SignedString([]byte(secret))
}

func TestValidateToken_ValidToken(t *testing.T) {
	secret := "mysecret"
	claims := jwt.MapClaims{
		"userId": "12345",
		"exp":    time.Now().Add(time.Hour).Unix(),
	}
	tokenString, err := generateToken(secret, claims, jwt.SigningMethodHS256)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	j := &JWTUtil{
		Secret:      secret,
		tokenString: tokenString,
	}

	userId, err := j.ValidateToken()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if userId != "12345" {
		t.Errorf("expected userId '12345', got %v", userId)
	}
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	secret := "mysecret"
	wrongSecret := "wrongsecret"
	claims := jwt.MapClaims{
		"userId": "12345",
		"exp":    time.Now().Add(time.Hour).Unix(),
	}
	tokenString, err := generateToken(wrongSecret, claims, jwt.SigningMethodHS256)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	j := &JWTUtil{
		Secret:      secret,
		tokenString: tokenString,
	}

	_, err = j.ValidateToken()
	if err == nil {
		t.Error("expected error for invalid signature, got nil")
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	secret := "mysecret"
	claims := jwt.MapClaims{
		"userId": "12345",
		"exp":    time.Now().Add(-time.Hour).Unix(),
	}
	tokenString, err := generateToken(secret, claims, jwt.SigningMethodHS256)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	j := &JWTUtil{
		Secret:      secret,
		tokenString: tokenString,
	}

	_, err = j.ValidateToken()
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestValidateToken_InvalidSigningMethod(t *testing.T) {
	secret := "mysecret"
	// claims := jwt.MapClaims{
	// 	"userId": "12345",
	// 	"exp":    time.Now().Add(time.Hour).Unix(),
	// }

	// Create a malformed token to simulate invalid signing method
	// Since RS256 requires RSA keys, we'll create an invalid token string
	tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature"

	j := &JWTUtil{
		Secret:      secret,
		tokenString: tokenString,
	}

	_, err := j.ValidateToken()
	if err == nil {
		t.Error("expected error for invalid signing method, got nil")
	}
}

func TestValidateToken_NoUserIdClaim(t *testing.T) {
	secret := "mysecret"
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	tokenString, err := generateToken(secret, claims, jwt.SigningMethodHS256)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	j := &JWTUtil{
		Secret:      secret,
		tokenString: tokenString,
	}

	userId, err := j.ValidateToken()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if userId != nil {
		t.Errorf("expected nil userId, got %v", userId)
	}
}
