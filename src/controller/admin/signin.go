package admin_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	admin_request "github.com/matheuswww/mystream/src/controller/model/admin/request"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (ac *adminController) Signin(c *gin.Context) {
	logger.Log("Init Admin Signin")
	var adminSigninRequest admin_request.Signin
	if err := c.ShouldBindJSON(&adminSigninRequest); err != nil {
		restErr := rest_err.NewBadRequestError("invalid fields")
		c.JSON(restErr.Code, restErr)
		return
	}
	token,restErr := ac.admin_service.Signin(adminSigninRequest.Email, adminSigninRequest.Password)
	if restErr != nil {
		c.JSON(restErr.Code, restErr)
		return
	}
	c.JSON(http.StatusCreated, token)
}
