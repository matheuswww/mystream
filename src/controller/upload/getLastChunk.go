package upload_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	upload_request "github.com/matheuswww/mystream/src/controller/model/upload/request"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (uc *uploadController) GetLastChunk(c *gin.Context) {
	logger.Log("Init GetLastChunk")
	var getLastChunkRequest upload_request.FileHash
	if err := c.ShouldBindJSON(&getLastChunkRequest); err != nil {
		restErr := rest_err.NewBadRequestError("invalid fields")
		c.JSON(restErr.Code, restErr)
		return
	}
	fileName, restErr := uc.uploadService.GetLastChunk(getLastChunkRequest)
	if restErr != nil {
		c.JSON(restErr.Code, restErr)
		return
	}
	c.JSON(http.StatusOK, struct{ Chunk string }{ Chunk: fileName })
}

