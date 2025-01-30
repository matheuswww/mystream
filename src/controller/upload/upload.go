package upload_controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	upload_request "github.com/matheuswww/mystream/src/controller/model/upload/request"
	upload_controller_util "github.com/matheuswww/mystream/src/controller/upload/util"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (uc *uploadController) UploadFile(w http.ResponseWriter, r *http.Request) {
	logger.Log("Init UploadFile Controller")
	var upgrader = websocket.Upgrader {
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	var conn *websocket.Conn
	var err error
	var wg sync.WaitGroup
	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying Upgrade: %v", err))
		restErr := rest_err.NewInternalServerError("server error")
		upload_controller_util.SendWsRes(restErr, conn)
		conn.Close()
		return
	}
	defer conn.Close()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logger.Error(fmt.Sprintf("Error trying ReadMessage: %v", err))
			restErr := rest_err.NewInternalServerError("server error")
			upload_controller_util.SendWsRes(restErr, conn)
			conn.Close()
			break 
		}
		var uploadRequest upload_request.UploadFile
		if err := json.Unmarshal(msg, &uploadRequest); err != nil {
			logger.Error(fmt.Sprintf("Error trying Unmarshal: %v", err))
			restErr := rest_err.NewBadRequestError("campos inválidos")
			upload_controller_util.SendWsRes(restErr, conn)
			conn.Close()
			break
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			uc.uploadService.UploadFile(conn, uploadRequest)
		}()
	}
	wg.Wait()
}
