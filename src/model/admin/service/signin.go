package admin_service

import (
	admin_response "github.com/matheuswww/mystream/src/controller/model/admin/response"
	admin_service_util "github.com/matheuswww/mystream/src/model/admin/service/util"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (as *adminService) Signin(email, password string) (*admin_response.Token, *rest_err.RestErr) {
	id,restErr := as.admin_repository.Signin(email, password)
	if restErr != nil {
		return nil, restErr
	}
	adminToken, restErr := admin_service_util.GetAdminToken(id, email)
	if restErr != nil {
		return nil,restErr
	}
	return adminToken, nil
}