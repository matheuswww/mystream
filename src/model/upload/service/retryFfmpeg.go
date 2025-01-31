package upload_service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/matheuswww/mystream/src/ffmpeg"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (us *uploadService) RetryFfmpeg(fileHash string) *rest_err.RestErr {
	path,err := filepath.Abs("upload")
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying get absolute path for upload: %v", err))
		restErr := rest_err.NewInternalServerError("server error")
		return restErr
	}
	fp := filepath.Join(path, fileHash)
	entries, err := os.ReadDir(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return rest_err.NewNotFoundError("path does not exist")
		}
		logger.Error(fmt.Sprintf("ReadDir: %v", err))
		return rest_err.NewInternalServerError("server error")
	}
	var found bool
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".mp4") {
			f := filepath.Join(fp, entry.Name())
			ffmpeg.SaveVideo(path, f, fileHash, nil)
			found = true
			break
		}
	}
	if !found {
		return rest_err.NewBadRequestError("video not found, either the file has already been processed or the chunks have not been sent")
	}
	return nil
}