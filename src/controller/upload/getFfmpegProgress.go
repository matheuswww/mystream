package upload_controller

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	upload_request "github.com/matheuswww/mystream/src/controller/model/upload/request"
	upload_controller_util "github.com/matheuswww/mystream/src/controller/upload/util"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
	"github.com/matheuswww/mystream/src/router"
)

func (uc *uploadController) GetFfmpegProgress(w http.ResponseWriter, r *http.Request) {
	logger.Log("Init GetFfmpegProgress")
	defer r.Body.Close()  
	var upgrader = websocket.Upgrader {
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
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
		if err := router.BindJson(bytes.NewReader(msg), &fileHash); err != nil {
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