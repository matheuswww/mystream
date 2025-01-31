package upload_service

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	upload_request "github.com/matheuswww/mystream/src/controller/model/upload/request"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (us *uploadService) GetLastChunk(getLastChunkRequest upload_request.FileHash) (string, *rest_err.RestErr) {
	path,err := filepath.Abs("upload")
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying get absolute path for upload: %v", err))
		restErr := rest_err.NewInternalServerError("server error")
		return "", restErr
	}

	dir := fmt.Sprintf("%s/%s/temp", path, getLastChunkRequest.FileHash)
	fileName,err := deleteLastChunk(dir)
	if err != nil {
		if os.IsNotExist(err) {
			restErr := rest_err.NewNotFoundError("nenhum arquivo foi encontrado")
			return "", restErr
		}
		logger.Error(err)
		restErr := rest_err.NewInternalServerError("server error")
		return "", restErr
	}
	if fileName == "" {
		restErr := rest_err.NewNotFoundError("nenhum arquivo foi encontrado")
		return "", restErr
	}

	fileName,err = getLastModifiedFile(dir)
	if err != nil {
		if os.IsNotExist(err) {
			restErr := rest_err.NewNotFoundError("nenhum arquivo foi encontrado")
			return "", restErr
		}
		logger.Error(err)
		restErr := rest_err.NewInternalServerError("server error")
		return "", restErr
	}
	if fileName == "" {
		restErr := rest_err.NewNotFoundError("nenhum arquivo foi encontrado")
		return "", restErr
	}
	return fileName, nil
}

func getLastModifiedFile(dir string) (string, error) {
	var lastModifiedFile string
	var lastModifiedTime time.Time

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.Error(fmt.Sprintf("Error trying WalkDir: %v", err))
			return err
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				logger.Error(fmt.Sprintf("Error trying checking if is dir: %v", err))
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
		logger.Error(fmt.Sprintf("Error trying RemovePath: %v", err))
		return "",err
	}
	return fileName,nil
}