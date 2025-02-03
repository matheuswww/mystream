package admin_service

import (
	"fmt"
	"os"

	admin_response "github.com/matheuswww/mystream/src/controller/model/admin/response"
	jwt_service "github.com/matheuswww/mystream/src/jwt"
	"github.com/matheuswww/mystream/src/logger"
	admin_service_util "github.com/matheuswww/mystream/src/model/admin/service/util"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (as *adminService) RefreshToken(refreshToken string) (*admin_response.Token, *rest_err.RestErr) {
	refreshClaims := jwt_service.ParseAdminRefreshToken(refreshToken)
	if refreshClaims == nil || refreshClaims.Valid() != nil {
		logger.Error(fmt.Sprintf("Error trying ParseAdminRefreshToken: invalid claims"))
		return nil, rest_err.NewBadRequestError("invalid refresh token")
	}
	id := refreshClaims.Subject
	email := os.Getenv("ADMIN_EMAIL")
	if email == "" {
		logger.Error("Error trying get env")
		return nil, rest_err.NewInternalServerError("server error")
	}
	adminToken, restErr := admin_service_util.GetAdminToken(id, email)
	if restErr != nil {
		return nil,restErr
	}
	return adminToken, nil
}