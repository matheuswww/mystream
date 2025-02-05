package upload_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matheuswww/mystream/src/logger"
)

func (uc *uploadController) GetVideo(c *gin.Context) {
	logger.Log("Init GetVideo")
	cursor := c.Query("cursor")
	response, restErr := uc.uploadService.GetVideo(cursor)
	if restErr != nil {
		c.JSON(restErr.Code, restErr)
		return
	}
	c.JSON(http.StatusOK, response)
}