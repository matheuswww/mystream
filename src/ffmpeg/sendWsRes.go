package ffmpeg

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/matheuswww/mystream/src/logger"
)

func sendWsRes(msg any, fileHash string) {
	conn := conns[fileHash]
	if conn == nil {
		return
	}
	err := conn.WriteJSON(msg)
	if err != nil {
		if websocket.IsCloseError(err) {
			conns[fileHash] = nil
			return
		}
		logger.Error(fmt.Sprintf("Error trying SendWsRes: %v", err))
	}
}