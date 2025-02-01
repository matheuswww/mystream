package user_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
	"golang.org/x/crypto/bcrypt"
)

func (ur *userRepository) Signup(id, email, name, password string) *rest_err.RestErr {
	ctx,cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var count int
	query := "SELECT COUNT(*) FROM users WHERE email = $1"
	err := ur.sql.QueryRowContext(ctx, query, email).Scan(&count)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying QueryRowContext: %v", err))
		return rest_err.NewInternalServerError("server error")
	}
	if count != 0 {
		return rest_err.NewBadRequestError("email already exists")
	}

	encryptedPassword,err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying GeneratePassword: %v", err))
		return rest_err.NewInternalServerError("server error")
	}
	query = "INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)"
	_,err = ur.sql.ExecContext(ctx, query, id, name, email, encryptedPassword)
	if err != nil {
		logger.Error(fmt.Sprintf("Error tyring ExecContext: %v", err))
		return rest_err.NewInternalServerError("server error")
	}
	return nil
}