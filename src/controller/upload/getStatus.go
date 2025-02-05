package upload_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	upload_request "github.com/matheuswww/mystream/src/controller/model/upload/request"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (uc *uploadController) GetStatus(c *gin.Context) {
	logger.Log("Init GetStatus")
	var getStatus upload_request.FileHash
	if err := c.ShouldBindJSON(&getStatus); err != nil {
		restErr := rest_err.NewBadRequestError("invalid fields")
		c.JSON(restErr.Code, restErr)
		return
	}
	status,restErr := uc.uploadService.GetStatus(getStatus.FileHash, BeingProcessed)
	if restErr != nil {
		c.JSON(restErr.Code, restErr)
		return
	}
	c.JSON(http.StatusOK, struct{Status string}{ Status: status })
}