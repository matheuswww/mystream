package user_controller

import (
	"github.com/gin-gonic/gin"
	user_service "github.com/matheuswww/mystream/src/model/user/service"
)



func NewUserController(user_service user_service.UserService) UserController {
	return &userController{
		user_service,
	}
}

type UserController interface {
	Signup(c *gin.Context)
	Signin(c *gin.Context)
}

type userController struct {
	user_service user_service.UserService
}
