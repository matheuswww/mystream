package upload_controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	upload_request "github.com/matheuswww/mystream/src/controller/model/upload/request"
	upload_controller_util "github.com/matheuswww/mystream/src/controller/upload/util"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

var beingProcessed map[string]bool = make(map[string]bool)

func GetBeingProcessed(fileHash string) bool {
	if _,v := beingProcessed[fileHash]; v {
		return true
	}
	return false
}

func (uc *uploadController) UploadFile(c *gin.Context) {
	logger.Log("Init UploadFile Controller")
	var upgrader = websocket.Upgrader {
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	var conn *websocket.Conn
	var err error
	conn, err = upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying Upgrade: %v", err))
		restErr := rest_err.NewInternalServerError("server error")
		upload_controller_util.SendWsRes(restErr, conn)
		conn.Close()
		return
	}
	var fileHash string
	var wg sync.WaitGroup
	defer func() {
		delete(beingProcessed, fileHash)
		conn.Close()
	}()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			break 
		}
		var uploadRequest upload_request.UploadFile
		if err := json.Unmarshal([]byte(msg), &fileHash); err != nil || uploadRequest.Chunks == nil || uploadRequest.FileHash == "" || uploadRequest.Filename == "" {
			logger.Error(fmt.Sprintf("Error trying Unmarshal: %v", err))
			restErr := rest_err.NewBadRequestError("invalid fields")
			upload_controller_util.SendWsRes(restErr, conn)
			conn.Close()
			break
		}
		folder, tempFolder, file, err := checkFile(uploadRequest.FileHash, conn)
		if err != nil {
			break
		}
		if (folder && !tempFolder && !file) || (folder && file) {
			restErr := rest_err.NewBadRequestError("file was already send")
			upload_controller_util.SendWsRes(restErr, conn)
			conn.Close()
			break
		}
		if fileHash == "" {
			beingProcessed[uploadRequest.FileHash] = true
			fileHash = uploadRequest.FileHash
		}
		wg.Add(1)
		go func() {
			wg.Done()
			uc.uploadService.UploadFile(conn, uploadRequest)
		}()
	}
	wg.Wait()
}

func checkFile(fileHash string, conn *websocket.Conn) (bool, bool, bool, error) {
	path,err := filepath.Abs("upload")
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying get abs path for upload: %v", err))
		restErr := rest_err.NewInternalServerError("server error")
		upload_controller_util.SendWsRes(restErr, conn)
		conn.Close()
		return false, false, false, err
	}
	fp := filepath.Join(path, fileHash)
	var folder, tempFolder, file bool
	entries, err := os.ReadDir(fp)
	if err != nil {
		if os.IsNotExist(err) {
			folder = false
			return false, false, false, nil
		} else {
			logger.Error(fmt.Sprintf("Error trying ReadDir: %v", err))
			restErr := rest_err.NewInternalServerError("server error")
			upload_controller_util.SendWsRes(restErr, conn)
			conn.Close()
			return false, false, false, err
		}
	} else {
		folder = true
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".mp4") {
			file = true
		}
		if entry.IsDir() && entry.Name() == "temp" {
			tempFolder = true
		}
	}
	return folder, tempFolder, file, nil
}