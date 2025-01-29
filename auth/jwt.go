package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY = []byte("secret_key")

func GenerateJsonWebToken(username string, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role": role,
		"expiration_date": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SECRET_KEY)
}

func ValidateJsonWebToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}
	return claims, nil
}