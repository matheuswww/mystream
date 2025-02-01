package user_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	user_request "github.com/matheuswww/mystream/src/controller/model/user/request"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (uc *userController) Signup(c *gin.Context) {
	logger.Log("Init Signup")
	var signupRequest user_request.Signup
	if err := c.ShouldBindJSON(&signupRequest); err != nil {
		restErr := rest_err.NewBadRequestError("invalid fields")
		c.JSON(restErr.Code, restErr)
		return
	}
	token,restErr := uc.user_service.Signup(signupRequest.Email, signupRequest.Name, signupRequest.Password)
	if restErr != nil {
		c.JSON(restErr.Code, restErr)
		return
	}
	c.JSON(http.StatusCreated, token)
}