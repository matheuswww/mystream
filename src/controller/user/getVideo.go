package user_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matheuswww/mystream/src/logger"
)

func (uc *userController) GetVideo(c *gin.Context) {
	logger.Log("Init GetVideo")
	cursor := c.Query("cursor")
	response, restErr := uc.user_service.GetVideo(cursor)
	if restErr != nil {
		c.JSON(restErr.Code, restErr)
		return
	}
	c.JSON(http.StatusOK, response)
}