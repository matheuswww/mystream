package upload_controller

import (
	"github.com/gin-gonic/gin"
	upload_service "github.com/matheuswww/mystream/src/model/upload/service"
)

func NewUploadController(uploadService upload_service.UploadService ) UploadController {
	return &uploadController {
		uploadService,
	}
}

type uploadController struct {
	uploadService upload_service.UploadService
}

type UploadController interface {
	UploadFile(c *gin.Context)
	GetLastChunk(c *gin.Context)
	GetFfmpegProgress(c *gin.Context)
	RetryFfmpeg(c *gin.Context)
	GetStatus(c *gin.Context)
}