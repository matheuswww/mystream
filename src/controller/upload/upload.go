package upload_controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"github.com/gorilla/websocket"
	upload_request "github.com/matheuswww/mystream/src/controller/model/upload/request"
	"github.com/matheuswww/mystream/src/logger"
)

var upgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (uc *uploadController) UploadFile(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err)
	}
	defer conn.Close()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logger.Error(err)
			return 
		}
		var uploadRequest upload_request.UploadFile
		if err := json.Unmarshal(msg, &uploadRequest); err != nil {
			logger.Error(err)
			return
		}
		go saveChunk(uploadRequest)
	}
}

func saveChunk(uploadFile upload_request.UploadFile) {
	dir := fmt.Sprintf("/home/virus/mystream/upload/%s/temp", uploadFile.FileHash)
	if _,err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			logger.Error(err)
			return
		}
	} else if err != nil {
		logger.Error(err)
		return
	}
	for _,chunk := range uploadFile.Chunks {
		filePath := fmt.Sprintf("%s/chunk%d", dir, chunk.Chunk)
		file,err := os.Create(filePath)
		if err != nil {
			logger.Error(err)
			return
		}
		defer file.Close()
		_,err = file.Write(chunk.Data)
		if err != nil {
			logger.Error(err)
			return
		}
		if chunk.Chunk == uploadFile.TotalChunk - 1 {
			combineChunk(uploadFile.TotalChunk, uploadFile.Filename, uploadFile.FileHash)
		}
	}
}

func combineChunk(totalChunk int, fileName, fileHash string) {
	filePath := fmt.Sprintf("/home/virus/mystream/upload/%s/%s", fileHash, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		logger.Error(err)
		return
	}
	defer file.Close()

	dir := fmt.Sprintf("/home/virus/mystream/upload/%s/temp", fileHash)
	for i := 0; i < totalChunk; i++ {
		chunkFileName := fmt.Sprintf("/%s/chunk%d", dir, i)
		chukData, err := os.ReadFile(chunkFileName)
		if err != nil {
			logger.Error(err)
			return
		}
		_, err = file.Write(chukData)
		if err != nil {
			logger.Error(err)
			return
		}
	}
	return
}
