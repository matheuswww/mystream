package upload_controller_util

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/matheuswww/mystream/src/logger"
)

func SendWsRes(msg any, conn *websocket.Conn) {
	b,err := json.Marshal(msg)
	if err != nil {
		logger.Error(err)
		return
	}
	err = conn.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}