package upload_routes

import (
	"database/sql"

	upload_controller "github.com/matheuswww/mystream/src/controller/upload"
	upload_repository "github.com/matheuswww/mystream/src/model/upload/repository"
	upload_service "github.com/matheuswww/mystream/src/model/upload/service"
	"github.com/matheuswww/mystream/src/router"
)

func InitUploadRoutes(r *router.Router, db *sql.DB) {
	controller := getUploadController(db)
	r.Route("GET", "/upload/uploadFile", controller.UploadFile)
	r.Route("GET", "/upload/getLastChunk", controller.GetLastChunk)
	r.Route("GET", "/upload/getFfmpegProgress", controller.GetFfmpegProgress)
	r.Route("GET", "/upload/getStatus", controller.GetStatus)
	r.Route("PATCH", "/upload/retryFfmpeg", controller.RetryFfmpeg)
}

func getUploadController(db *sql.DB) upload_controller.UploadController {
	repository := upload_repository.NewUploadRepository(db)
	service := upload_service.NewUploadService(repository)
	return upload_controller.NewUploadController(service)
}