package admin_controller_util

import (
	"errors"
	"strings"

	"github.com/matheuswww/mystream/src/logger"
)

func GetToken(token string) (string, error) {
	logger.Log("Init GetToken")
	if token == "" {
		logger.Error("Error trying get refreshToken,authorization is nil")
		err := errors.New("authorization is nil")
		return "", err
	}
	parts := strings.Split(token, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		logger.Error("Error trying get refreshToken,bad authorization header")
		err := errors.New("bad authorization header")
		return "", err
	}
	refreshToken := parts[1]
	return refreshToken, nil
}