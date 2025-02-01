package user_service

import (
	"fmt"
	user_response "github.com/matheuswww/mystream/src/controller/model/user/response"
	jwt_service "github.com/matheuswww/mystream/src/jwt"
	"github.com/matheuswww/mystream/src/logger"
	user_service_util "github.com/matheuswww/mystream/src/model/user/service/util"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (us *userService) RefreshToken(refreshToken string) (*user_response.Token, *rest_err.RestErr) {
	refreshClaims := jwt_service.ParseRefreshToken(refreshToken)
	if refreshClaims == nil || refreshClaims.Valid() != nil {
		logger.Error(fmt.Sprintf("Error trying ParseRefreshToken: invalid claims"))
		return nil, rest_err.NewBadRequestError("invalid refresh token")
	}
	id := refreshClaims.Subject
	email, restErr := us.user_repository.GetEmailById(id)
	if restErr != nil {
		return nil,restErr
	}
	userToken, restErr := user_service_util.GetUserToken(id, email)
	if restErr != nil {
		return nil,restErr
	}
	return userToken, nil
}