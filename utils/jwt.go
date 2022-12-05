package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var UserRolesEnums = map[string]string{
	"user": "user",
	"admin": "admin",
}

var JwtSigningKey string = "12345678"

type JwtClaims struct {
	ID uint `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJwtToken (ID uint, role string) *jwt.Token {
	claims := JwtClaims{
		ID,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "rt-server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token
}

func ParseJwtToken (tokenString  string) (JwtClaims, error) {
	claims := JwtClaims{}
	_, error := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSigningKey), nil
	})

	return claims, error
}