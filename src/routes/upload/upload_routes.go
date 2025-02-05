package upload_routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	upload_controller "github.com/matheuswww/mystream/src/controller/upload"
	upload_repository "github.com/matheuswww/mystream/src/model/upload/repository"
	upload_service "github.com/matheuswww/mystream/src/model/upload/service"
)

func InitUploadRoutes(r *gin.Engine, db *sql.DB) {
	controller := getUploadController(db)
	upload := r.Group("/upload")
	upload.GET("/getLastChunk", controller.CheckToken, controller.GetLastChunk)
	upload.GET("/getStatus", controller.CheckToken, controller.GetStatus)
	upload.PATCH("/retryFfmpeg", controller.CheckToken, controller.RetryFfmpeg)
	upload.PATCH("/updateVideo", controller.CheckToken, controller.UpdateVideo)
	upload.GET("/getVideo", controller.CheckToken, controller.GetVideo)
	
	upload.GET("/uploadFile", controller.UploadFile)
	upload.GET("/getFfmpegProgress", controller.GetFfmpegProgress)
}

func getUploadController(db *sql.DB) upload_controller.UploadController {
	repository := upload_repository.NewUploadRepository(db)
	service := upload_service.NewUploadService(repository)
	return upload_controller.NewUploadController(service)
}