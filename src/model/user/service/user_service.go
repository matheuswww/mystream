package user_service

import (
	user_repository "github.com/matheuswww/mystream/src/model/user/repository"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func NewUserService(userRepository user_repository.UserRepository) UserService {
	return &userService {
		userRepository,
	}
}

type UserService interface {
	Signup(email, name, password string) (string, *rest_err.RestErr)
}

type userService struct {
	user_repository user_repository.UserRepository
}