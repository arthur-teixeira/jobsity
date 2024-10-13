package service

import (
	"github.com/golang-jwt/jwt"
	"time"
)

var SECRET = []byte("VERY_SECRET_SECRET")

func CreateJWTToken(name, email string) (string, error) {
	claims := jwt.MapClaims{
		"name":  name,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(SECRET)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) *jwt.Token {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SECRET, nil
	})

	if err != nil {
		return nil
	}

	return token
}
