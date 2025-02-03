package admin_repository

import (
	"database/sql"

	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func NewAdminRepository(sql *sql.DB) AdminRepository {
	return &adminRepository {
		sql,
	}
}

type AdminRepository interface {
	Signin(email, password string) (string, *rest_err.RestErr)
}

type adminRepository struct {
	sql *sql.DB
}