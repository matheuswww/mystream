package admin_controller

import (
	"github.com/gin-gonic/gin"
	admin_service "github.com/matheuswww/mystream/src/model/admin/service"
)

func NewAdminController(service admin_service.AdminService) AdminController {
	return &adminController{
		service,
	} 
}

type AdminController interface {
	Signin(c *gin.Context)
}

type adminController struct {
	admin_service admin_service.AdminService
}