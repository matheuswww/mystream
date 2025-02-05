package upload_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	upload_request "github.com/matheuswww/mystream/src/controller/model/upload/request"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (uc *uploadController) UpdateVideo(c *gin.Context) {
	logger.Log("Init UpdateVideo")
	var updateVideoRequest upload_request.UpdateVideo
	if err := c.ShouldBindJSON(&updateVideoRequest); err != nil {
		restErr := rest_err.NewBadRequestError("invalid fields")
		c.JSON(restErr.Code, restErr)
		return
	}
	if updateVideoRequest.Description == "" && updateVideoRequest.Title == "" && updateVideoRequest.Uploaded == nil {
		logger.Error("Error trying UpdateVideo: no parameters sent")
		restErr := rest_err.NewBadRequestError("no parameters sent")
		c.JSON(restErr.Code, restErr)
		return
	}
	restErr := uc.uploadService.UpdateVideo(updateVideoRequest.FileHash, updateVideoRequest.Title, updateVideoRequest.Description, updateVideoRequest.Uploaded, BeingProcessed)
	if restErr != nil {
		c.JSON(restErr.Code, restErr)
		return
	}
	c.Status(http.StatusOK)
}