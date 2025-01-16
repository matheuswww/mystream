package upload_controller

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"

	upload_request "github.com/matheuswww/mystream/src/controller/model/upload/request"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
	"github.com/matheuswww/mystream/src/router"
)

func (uc *uploadController) GetLastChunk(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		restErr := rest_err.NewInternalServerError("server error")
		router.SendResponse(w, restErr, restErr.Code)
		return
	}
	defer r.Body.Close() 
	var getLastChunkRequest upload_request.GetLastChunk
	if err := json.Unmarshal(body, &getLastChunkRequest); err != nil {
		restErr := rest_err.NewBadRequestError("campos inválidos")
		router.SendResponse(w, restErr, restErr.Code)
		return
	}
	path,err := filepath.Abs("upload")
	if err != nil {
		restErr := rest_err.NewInternalServerError("server error")
		router.SendResponse(w, restErr, restErr.Code)
		return
	}

	dir := fmt.Sprintf("%s/%s/temp", path, getLastChunkRequest.FileHash)
	fileName,err := deleteLastChunk(dir)
	if err != nil {
		if os.IsNotExist(err) {
			restErr := rest_err.NewNotFoundError("nenhum arquivo foi encontrado")
			router.SendResponse(w, restErr, restErr.Code)
			return
		}
		logger.Error(err)
		restErr := rest_err.NewInternalServerError("server error")
		router.SendResponse(w, restErr, restErr.Code)
		return
	}
	if fileName == "" {
		restErr := rest_err.NewNotFoundError("nenhum arquivo foi encontrado")
		router.SendResponse(w, restErr, restErr.Code)
		return
	}

	fileName,err = getLastModifiedFile(dir)
	if err != nil {
		if os.IsNotExist(err) {
			restErr := rest_err.NewNotFoundError("nenhum arquivo foi encontrado")
			router.SendResponse(w, restErr, restErr.Code)
		}
		logger.Error(err)
		restErr := rest_err.NewInternalServerError("server error")
		router.SendResponse(w, restErr, restErr.Code)
		return
	}
	if fileName == "" {
		restErr := rest_err.NewNotFoundError("nenhum arquivo foi encontrado")
		router.SendResponse(w, restErr, restErr.Code)
		return
	}

	router.SendResponse(w, struct{ Chunk string }{ Chunk: fileName }, http.StatusOK)
}

func getLastModifiedFile(dir string) (string, error) {
	var lastModifiedFile string
	var lastModifiedTime time.Time

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			modTime := info.ModTime()
			if modTime.After(lastModifiedTime) {
				lastModifiedTime = modTime
				lastModifiedFile = filepath.Base(path)
			}
		}
		return nil
	})
	return lastModifiedFile, err
}

func deleteLastChunk(dir string) (string,error) {
	fileName, err := getLastModifiedFile(dir)
	if err != nil {
		return "",err
	}
	path := fmt.Sprintf("%s/%s", dir, fileName)
	err = os.Remove(path)
	if err != nil {
		return "",err
	}
	return fileName,nil
}