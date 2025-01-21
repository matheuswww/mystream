package file_controller

import "net/http"

func NewFileoController() FileController {
	return &fileController{}
}

type FileController interface {
	ServeFile(w http.ResponseWriter, r *http.Request)
}

type fileController struct {}