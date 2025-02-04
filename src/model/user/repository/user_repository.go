package user_repository

import (
	"database/sql"

	user_response "github.com/matheuswww/mystream/src/controller/model/user/response"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func NewUserRepository(sql *sql.DB) UserRepository {
	return &userRepository { sql }
}

type UserRepository interface {
	Signup(id, email, name, password string) *rest_err.RestErr
	Signin(email, password string) (string, *rest_err.RestErr)
	GetEmailById(id string) (string, *rest_err.RestErr)
	GetVideo(cursor string) ([]user_response.GetVideo, *rest_err.RestErr)
}

type userRepository struct {
	sql *sql.DB
}
