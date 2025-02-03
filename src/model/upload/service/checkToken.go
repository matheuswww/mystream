package upload_service

import (
	"fmt"

	jwt_service "github.com/matheuswww/mystream/src/jwt"
	"github.com/matheuswww/mystream/src/logger"
)

func (as *uploadService) CheckToken(refreshToken string) bool {
	refreshClaims := jwt_service.ParseAdminRefreshToken(refreshToken)
	if refreshClaims == nil || refreshClaims.Valid() != nil {
		logger.Error(fmt.Sprintf("Error trying ParseAdminRefreshToken: invalid claims"))
		return false
	}
	return true
}