package admin_service

import (
	admin_response "github.com/matheuswww/mystream/src/controller/model/admin/response"
	admin_repository "github.com/matheuswww/mystream/src/model/admin/repository"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func NewAdminService(repository admin_repository.AdminRepository) AdminService {
	return &adminService {
		repository,
	}
}

type AdminService interface {
	Signin(email, password string) (*admin_response.Token, *rest_err.RestErr)
	RefreshToken(refreshToken string) (*admin_response.Token, *rest_err.RestErr)
}

type adminService struct {
	admin_repository admin_repository.AdminRepository
}
