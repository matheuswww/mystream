package upload_service

import (
	"github.com/gorilla/websocket"
	upload_controller_util "github.com/matheuswww/mystream/src/controller/upload/util"
	"github.com/matheuswww/mystream/src/ffmpeg"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (us *uploadService) GetFfmpegProgress(fileHash string, conn *websocket.Conn) {
	logger.Log("Init NewFfmpegSession")
	err := ffmpeg.UpdateConn(fileHash, conn)
	if err != nil {
		restErr := rest_err.NewInternalServerError("arquivo não encontrado")
		upload_controller_util.SendWsRes(restErr, conn)
		conn.Close()
	}
}