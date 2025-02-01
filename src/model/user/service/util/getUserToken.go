package user_service_util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	user_response "github.com/matheuswww/mystream/src/controller/model/user/response"
	jwt_service "github.com/matheuswww/mystream/src/jwt"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func GetUserToken(id, email string) (*user_response.Token, *rest_err.RestErr) {
	token, err := jwt_service.NewAccessToken(jwt_service.UserClaims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			Subject: id,
			IssuedAt: time.Now().Unix(),
			ExpiresAt: time.Now().Add(jwt_service.ExpToken).Unix(),
		},
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying NewAccessToken: %v", err))
		return nil, rest_err.NewInternalServerError("server error")
	}
	refreshToken, err := jwt_service.NewRefreshToken(jwt.StandardClaims{
		Subject: id,
		IssuedAt: time.Now().Unix(),
		ExpiresAt: time.Now().Add(jwt_service.ExpRefreshToken).Unix(),
	},)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying NewRefreshToken: %v", err))
		return nil, rest_err.NewInternalServerError("server error")
	}
	return &user_response.Token{
		Token: token,
		RefreshToken: refreshToken,
	}, nil
}