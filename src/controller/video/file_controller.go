package file_controller

import (
	"github.com/gin-gonic/gin"
)

func NewFileoController() FileController {
	return &fileController{}
}

type FileController interface {
	ServeFile(c *gin.Context)
}

type fileController struct {}