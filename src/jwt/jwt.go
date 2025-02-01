package jwt_service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var ExpToken = time.Minute*60
var ExpRefreshToken = time.Hour*24*2

type UserClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func NewAccessToken(claims UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return accessToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func NewRefreshToken(claims jwt.StandardClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return refreshToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func ParseAccessToken(accessToken string) *UserClaims {
	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil || !parsedAccessToken.Valid {
		return nil
	}
	claims, ok := parsedAccessToken.Claims.(*UserClaims)
	if !ok {
		return nil
	}
	return claims
}

func ParseRefreshToken(refreshToken string) *jwt.StandardClaims {
	parsedRefreshToken, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil || !parsedRefreshToken.Valid {
		return nil
	}
	claims, ok := parsedRefreshToken.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil
	}
	return claims
}
