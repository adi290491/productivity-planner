package middleware

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWTUtil struct {
	Secret      string
	tokenString string
}

func (j *JWTUtil) ValidateToken() (interface{}, error) {
	token, err := jwt.Parse(j.tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	return claims["userId"], nil
}
