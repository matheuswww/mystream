package user_service

import (
	user_response "github.com/matheuswww/mystream/src/controller/model/user/response"
	user_service_util "github.com/matheuswww/mystream/src/model/user/service/util"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (us *userService) Signin(email, password string) (*user_response.Token, *rest_err.RestErr) {
	id,restErr := us.user_repository.Signin(email, password)
	if restErr != nil {
		return nil, restErr
	}
	userToken, restErr := user_service_util.GetUserToken(id, email)
	if restErr != nil {
		return nil,restErr
	}
	return userToken, nil
}