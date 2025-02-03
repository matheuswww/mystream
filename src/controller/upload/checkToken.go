package upload_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	admin_controller_util "github.com/matheuswww/mystream/src/controller/admin/util"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (uc *uploadController) CheckToken(c *gin.Context) {
	logger.Log("Init CheckToken")
	authHeader := c.GetHeader("Authorization")
	token, err := admin_controller_util.GetToken(authHeader)
	if err != nil {
		restErr := rest_err.NewBadRequestError(err.Error())
		c.AbortWithStatusJSON(restErr.Code, restErr)
		return
	}
	valid := uc.uploadService.CheckToken(token)
	if !valid {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Next()
}