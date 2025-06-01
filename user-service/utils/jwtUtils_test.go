package utils

import (
	"productivity-planner/user-service/models"
	"testing"
	"time"

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
