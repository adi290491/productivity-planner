package utils

import (
	"fmt"
	"productivity-planner/user-service/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTUtil struct {
	Secret []byte
}

func (j *JWTUtil) GenerateToken(user *models.User) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": user.ID,
			"exp":    time.Now().Add(24 * time.Hour).Unix(),
		})

	token, err := t.SignedString(j.Secret)
	if err != nil {
		return "", fmt.Errorf("error while generating token: %v", err)
	}

	return token, nil
}

func (j *JWTUtil) ValidateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.Secret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
