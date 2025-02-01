package user_repository

import (
	"database/sql"

	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func NewUserRepository(sql *sql.DB) UserRepository {
	return &userRepository { sql }
}

type UserRepository interface {
	Signup(id, email, name, password string) *rest_err.RestErr
	Signin(email, password string) (string, *rest_err.RestErr)
}

type userRepository struct {
	sql *sql.DB
}
