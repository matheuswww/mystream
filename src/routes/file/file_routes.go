package file_routes

import (
	file_controller "github.com/matheuswww/mystream/src/controller/video"
	"github.com/matheuswww/mystream/src/router"
)

func InitFileRoutes(r *router.Router) {
	controller := getFileController()
	r.Route("GET", "/file", controller.ServeFile)
}

func getFileController() file_controller.FileController {
	return file_controller.NewFileoController()
}