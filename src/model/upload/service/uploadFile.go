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

func (us *uploadService) UploadFile(conn *websocket.Conn, uploadFile upload_request.UploadFile, id string) {
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
			err := combineChunk(uploadFile.TotalChunk, uploadFile.Filename, uploadFile.FileHash, conn)
			if err == nil {
				updated := true
				us.uploadRepository.UpdateVideo(uploadFile.FileHash, "", "", &updated)
			}
		}
	}	
}

func combineChunk(totalChunk int, fileName, fileHash string, conn *websocket.Conn) error {
	logger.Log("Init combineChunk")
	path,err := filepath.Abs("upload")
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying get abs path for upload: %v", err))
		restErr := rest_err.NewInternalServerError("server error")
		upload_controller_util.SendWsRes(restErr, conn)
		conn.Close()
		return err
	}
	filePath := fmt.Sprintf("%s/%s/%s", path, fileHash, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying Create file: %v", err))
		restErr := rest_err.NewInternalServerError("server error")
		upload_controller_util.SendWsRes(restErr, conn)
		conn.Close()
		return err
	}
	defer file.Close()

	var cerr error
	dir := fmt.Sprintf("%s/%s/temp", path, fileHash)
	for i := 0; i < totalChunk; i++ {
		chunkFileName := fmt.Sprintf("/%s/chunk%d", dir, i)
		var chunkData []byte
		chunkData, cerr = os.ReadFile(chunkFileName)
		if cerr != nil {
			logger.Error(fmt.Sprintf("Error trying ReadFile: %v", cerr))
			restErr := rest_err.NewInternalServerError("server error")
			upload_controller_util.SendWsRes(restErr, conn)
			conn.Close()
			break
		}
		_, cerr = file.Write(chunkData)
		if cerr != nil {
			logger.Error(fmt.Sprintf("Error trying Write file: %v", cerr))
			restErr := rest_err.NewInternalServerError("server error")
			upload_controller_util.SendWsRes(restErr, conn)
			conn.Close()
			break
		}
	}
	if cerr != nil {
		return cerr
	}
	res := struct{ Message string }{
		Message: "sucesso",
	}
	upload_controller_util.SendWsRes(res, conn)
	err = os.RemoveAll(dir)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying RemoveAll: %v", err))
		return err
	}
	err = ffmpeg.SaveVideo(path, filePath, fileHash, conn)
	if err != nil {
		restErr := rest_err.NewInternalServerError("server error")
		upload_controller_util.SendWsRes(restErr, conn)
		conn.Close()
		return err
	}
	return nil
}