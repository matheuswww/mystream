package upload_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	upload_request "github.com/matheuswww/mystream/src/controller/model/upload/request"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (uc *uploadController) RetryFfmpeg(c *gin.Context) {
	logger.Log("Init RetryFfmpeg")
	var retryFfmpeg upload_request.FileHash
	if err := c.ShouldBindJSON(&retryFfmpeg); err != nil {
		restErr := rest_err.NewBadRequestError("invalid fields")
		c.JSON(restErr.Code, restErr)
		return
	}
	restErr := uc.uploadService.RetryFfmpeg(retryFfmpeg.FileHash)
	if restErr != nil {
		c.JSON(restErr.Code, restErr)
		return
	}
	c.JSON(http.StatusOK, struct{ Message string }{ Message: "success" })
}