package upload_controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	upload_request "github.com/matheuswww/mystream/src/controller/model/upload/request"
	upload_controller_util "github.com/matheuswww/mystream/src/controller/upload/util"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (uc *uploadController) GetFfmpegProgress(c *gin.Context) {
	logger.Log("Init GetFfmpegProgress")
	token := c.DefaultQuery("token", "")
	if token == "" {
		c.Status(http.StatusForbidden)
		return
	}
	valid := uc.uploadService.CheckToken(token)
	if !valid {
		c.Status(http.StatusForbidden)
		return
	}
	var upgrader = websocket.Upgrader {
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying Upgrade: %v", err))
		restErr := rest_err.NewInternalServerError("server error")
		upload_controller_util.SendWsRes(restErr, conn)
		conn.Close()
		return
	}

	var wg sync.WaitGroup
	defer conn.Close()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			break 
		}
		var fileHash upload_request.FileHash
		if err := json.Unmarshal([]byte(msg), &fileHash); err != nil || fileHash.FileHash == "" {
			logger.Error(fmt.Sprintf("Error trying Unmarshal: %v", err))
			restErr := rest_err.NewBadRequestError("invalid fields")
			upload_controller_util.SendWsRes(restErr, conn)
			conn.Close()
			return
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			uc.uploadService.GetFfmpegProgress(fileHash.FileHash, conn)
		}()
	}
	wg.Wait()
}