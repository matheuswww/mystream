package user_service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	user_response "github.com/matheuswww/mystream/src/controller/model/user/response"
	jwt_service "github.com/matheuswww/mystream/src/jwt"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (us *userService) Signin(email, password string) (*user_response.Token, *rest_err.RestErr) {
	id,restErr := us.user_repository.Signin(email, password)
	if restErr != nil {
		return nil, restErr
	}
	standardClaims := jwt.StandardClaims{
		IssuedAt: time.Now().Unix(),
		ExpiresAt: time.Now().Add(expToken).Unix(),
	}
	token, err := jwt_service.NewAccessToken(jwt_service.UserClaims{
		Id: id,
		Email: email,
		StandardClaims: standardClaims,
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying NewAccessToken: %v", err))
		return nil, rest_err.NewInternalServerError("server error")
	}
	refreshToken, err := jwt_service.NewRefreshToken(standardClaims)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying NewRefreshToken: %v", err))
		return nil, rest_err.NewInternalServerError("server error")
	}
	return &user_response.Token{
		Token: token,
		RefreshToken: refreshToken,
	}, nil
}