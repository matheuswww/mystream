package upload_controller

import (
	"net/http"

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
	UploadFile(w http.ResponseWriter, r *http.Request)
}