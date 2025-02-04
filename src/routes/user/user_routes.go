package user_routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	user_controller "github.com/matheuswww/mystream/src/controller/user"
	user_repository "github.com/matheuswww/mystream/src/model/user/repository"
	user_service "github.com/matheuswww/mystream/src/model/user/service"
)

func InitUserRoutes(r *gin.Engine, sql *sql.DB) {
	controller := getUserController(sql)
	user := r.Group("/user")
  user.POST("/signup", controller.Signup)
	user.POST("/signin", controller.Signin)
	user.GET("/refreshToken", controller.RefreshToken)
	user.GET("/getVideo", controller.GetVideo)
}

func getUserController(sql *sql.DB) user_controller.UserController {
	repository := user_repository.NewUserRepository(sql)
	service := user_service.NewUserService(repository)
	controller := user_controller.NewUserController(service)
	return controller
}