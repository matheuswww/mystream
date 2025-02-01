package user_service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	jwt_service "github.com/matheuswww/mystream/src/jwt"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (us *userService) Signin(email, password string) (string, *rest_err.RestErr) {
	id,restErr := us.user_repository.Signin(email, password)
	if restErr != nil {
		return "", restErr
	}
	token, err := jwt_service.NewAccessToken(jwt_service.UserClaims{
		Id: id,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
			ExpiresAt: time.Now().Add(expToken).Unix(),
		},
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying NewAccessToken: %v", err))
		return "", rest_err.NewInternalServerError("server error")
	}
	return token,nil
}