package upload_service

import (
	"fmt"
	jwt_service "github.com/matheuswww/mystream/src/jwt"
	"github.com/matheuswww/mystream/src/logger"
)

func (as *uploadService) CheckToken(token string) bool {
	refreshClaims := jwt_service.ParseAdminAccessToken(token)
	if refreshClaims == nil || refreshClaims.Valid() != nil {
		logger.Error(fmt.Sprintf("Error trying ParseAdminAccessToken: invalid claims"))
		return false
	}
	return true
}
