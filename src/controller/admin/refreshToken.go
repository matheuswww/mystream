package admin_controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (ac *adminController) RefreshToken(c *gin.Context) {
	logger.Log("Init RefreshToken")
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		logger.Error("Error trying get refreshToken,authorization is nil")
		restErr := rest_err.NewBadRequestError("authorization is nil")
		c.JSON(restErr.Code, restErr)
		return
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		restErr := rest_err.NewBadRequestError("bad authorization header")
		c.JSON(restErr.Code, restErr)
	}
	refreshToken := parts[1]
	token,restErr := ac.admin_service.RefreshToken(refreshToken)
	if restErr != nil {
		c.JSON(restErr.Code, restErr)
		return
	}
	c.JSON(http.StatusOK, token)
}