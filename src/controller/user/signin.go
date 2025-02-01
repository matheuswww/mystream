package user_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	user_request "github.com/matheuswww/mystream/src/controller/model/user/request"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (uc *userController) Signin(c *gin.Context) {
	logger.Log("Init Signin")
	var signinRequest user_request.Signin
	if err := c.ShouldBindJSON(&signinRequest); err != nil {
		restErr := rest_err.NewBadRequestError("invalid fields")
		c.JSON(restErr.Code, restErr)
		return
	}
	token,restErr := uc.user_service.Signin(signinRequest.Email, signinRequest.Password)
	if restErr != nil {
		c.JSON(restErr.Code, restErr)
		return
	}
	c.JSON(http.StatusCreated, token)
}