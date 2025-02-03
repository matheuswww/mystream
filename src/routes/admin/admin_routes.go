package admin_routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	admin_controller "github.com/matheuswww/mystream/src/controller/admin"
	admin_repository "github.com/matheuswww/mystream/src/model/admin/repository"
	admin_service "github.com/matheuswww/mystream/src/model/admin/service"
)

func InitAdminRouter(r *gin.Engine, sql *sql.DB) {
	controller := getAdminController(sql)
	admin := r.Group("/admin")
  admin.POST("/signin", controller.Signin)
}

func getAdminController(sql *sql.DB) admin_controller.AdminController {
	repository := admin_repository.NewAdminRepository(sql)
	service := admin_service.NewAdminService(repository)
	controller := admin_controller.NewAdminController(service)
	return controller
}