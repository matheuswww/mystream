package file_routes

import (
	"github.com/gin-gonic/gin"
	file_controller "github.com/matheuswww/mystream/src/controller/video"
)

func InitFileRoutes(r *gin.Engine) {
	controller := getFileController()
	file := r.Group("/file")
	file.Use(controller.ServeFile)
	r.Static("/file", "../upload")
}

func getFileController() file_controller.FileController {
	return file_controller.NewFileoController()
}