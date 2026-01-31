package util

import (
	"fmt"
	"time"

	"github.com/adi290491/productivity-planner/user-service/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

// JWTUtil handles JWT token generation and validation
type JWTUtil struct {
	Secret []byte
}

// NewJWTUtil creates a new JWT utility instance
func NewJWTUtil(secret []byte) *JWTUtil {
	return &JWTUtil{Secret: secret}
}

// GenerateToken generates a JWT token for the given user
func (j *JWTUtil) GenerateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"userId": user.ID.String(),
		"email":  user.Email,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.Secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (j *JWTUtil) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.Secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return &claims, nil
}
