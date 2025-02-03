package jwt_service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var AdminExpToken = time.Minute * 60
var AdminExpRefreshToken = time.Hour * 24 * 2

type AdminClaims struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func NewAdminAccessToken(claims AdminClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return accessToken.SignedString([]byte(os.Getenv("ADMIN_TOKEN_SECRET")))
}

func NewAdminRefreshToken(claims jwt.StandardClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return refreshToken.SignedString([]byte(os.Getenv("ADMIN_TOKEN_SECRET")))
}

func ParseAdminAccessToken(accessToken string) *AdminClaims {
	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ADMIN_TOKEN_SECRET")), nil
	})
	if err != nil || !parsedAccessToken.Valid {
		return nil
	}
	claims, ok := parsedAccessToken.Claims.(*AdminClaims)
	if !ok {
		return nil
	}
	return claims
}

func ParseAdminRefreshToken(refreshToken string) *jwt.StandardClaims {
	parsedRefreshToken, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ADMIN_TOKEN_SECRET")), nil
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
