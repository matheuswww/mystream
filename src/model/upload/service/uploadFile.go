package upload_service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
	upload_request "github.com/matheuswww/mystream/src/controller/model/upload/request"
	upload_controller_util "github.com/matheuswww/mystream/src/controller/upload/util"
	"github.com/matheuswww/mystream/src/ffmpeg"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (us *uploadService) UploadFile(conn *websocket.Conn, uploadFile upload_request.UploadFile) {
	path,err := filepath.Abs("upload")
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying get abs path for upload: %v", err))
		restErr := rest_err.NewInternalServerError("server error")
		upload_controller_util.SendWsRes(restErr, conn)
		conn.Close()
		return
	}
	dir := fmt.Sprintf("%s/%s/temp", path, uploadFile.FileHash)
	if _,err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			logger.Error(fmt.Sprintf("Error trying MkdirAll: %v", err))
			restErr := rest_err.NewInternalServerError("server error")
			upload_controller_util.SendWsRes(restErr, conn)
			conn.Close()
			return
		}
	} else if err != nil {
		logger.Error(fmt.Sprintf("Error trying get stat of file: %v", err))
		restErr := rest_err.NewInternalServerError("server error")
		upload_controller_util.SendWsRes(restErr, conn)
		conn.Close()
		return
	}
	for _,chunk := range uploadFile.Chunks {
		hash := sha256.Sum256(chunk.Data) 
		if hex.EncodeToString(hash[:]) != chunk.Hash {
			logger.Error("chunk hash is different")
			restErr := rest_err.NewBadRequestError("chunck hash is different")
			upload_controller_util.SendWsRes(restErr, conn)
			conn.Close()
			return
		}
		filePath := fmt.Sprintf("%s/chunk%d", dir, chunk.Chunk)
		file,err := os.Create(filePath)
		if err != nil {
			logger.Error(fmt.Sprintf("Error trying Create file: %v", err))
			restErr := rest_err.NewInternalServerError("server error")
			upload_controller_util.SendWsRes(restErr, conn)
			conn.Close()
			return
		}
		defer file.Close()
		_,err = file.Write(chunk.Data)
		if err != nil {
			logger.Error(fmt.Sprintf("Error trying Write file: %v", err))
			restErr := rest_err.NewInternalServerError("server error")
			upload_controller_util.SendWsRes(restErr, conn)
			conn.Close()
			return
		}
		if chunk.Chunk == uploadFile.TotalChunk - 1 {
			combineChunk(uploadFile.TotalChunk, uploadFile.Filename, uploadFile.FileHash, conn)
		}
	}	
}

func combineChunk(totalChunk int, fileName, fileHash string, conn *websocket.Conn) {
	logger.Log("Init combineChunk")
	path,err := filepath.Abs("upload")
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying get abs path for upload: %v", err))
		restErr := rest_err.NewInternalServerError("server error")
		upload_controller_util.SendWsRes(restErr, conn)
		conn.Close()
		return 
	}
	filePath := fmt.Sprintf("%s/%s/%s", path, fileHash, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying Create file: %v", err))
		restErr := rest_err.NewInternalServerError("server error")
		upload_controller_util.SendWsRes(restErr, conn)
		conn.Close()
		return
	}
	defer file.Close()

	dir := fmt.Sprintf("%s/%s/temp", path, fileHash)
	for i := 0; i < totalChunk; i++ {
		chunkFileName := fmt.Sprintf("/%s/chunk%d", dir, i)
		chukData, err := os.ReadFile(chunkFileName)
		if err != nil {
			logger.Error(fmt.Sprintf("Error trying ReadFile: %v", err))
			restErr := rest_err.NewInternalServerError("server error")
			upload_controller_util.SendWsRes(restErr, conn)
			conn.Close()
			return
		}
		_, err = file.Write(chukData)
		if err != nil {
			logger.Error(fmt.Sprintf("Error trying Write file: %v", err))
			restErr := rest_err.NewInternalServerError("server error")
			upload_controller_util.SendWsRes(restErr, conn)
			conn.Close()
			return
		}
	}
	res := struct{ Message string }{
		Message: "sucesso",
	}
	upload_controller_util.SendWsRes(res, conn)
	err = os.RemoveAll(dir)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying RemoveAll: %v", err))
	}
	err = ffmpeg.SaveVideo(path, filePath, fileHash, conn)
	if err != nil {
		restErr := rest_err.NewInternalServerError("server error")
		upload_controller_util.SendWsRes(restErr, conn)
		conn.Close()
		return
	}
}