package user_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (ur *userRepository) GetEmailById(id string) (string, *rest_err.RestErr) {
	ctx,cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	query := "SELECT email FROM users WHERE id = $1"
	var email string
	err := ur.sql.QueryRowContext(ctx, query, id).Scan(&email)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying QueryRowContext: %v", err))
		return "", rest_err.NewInternalServerError("server error")
	}
	return email,nil
}