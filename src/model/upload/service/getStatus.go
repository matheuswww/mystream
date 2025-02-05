package upload_service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/matheuswww/mystream/src/ffmpeg"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

var BeingProcessed = "video is being processed"
var NotSent = "the video was not sent, call getLastChunk and send all chunks"
var NotBeingProcessed = "video was not being processed, call retryFfmpeg"
var NotSavedInDb = "video was processed, but the status was not saved in the db, call updateVideo"
var Processed = "video processed"

func (us *uploadService) GetStatus(fileHash string, beingProcessed map[string]bool) (string, *rest_err.RestErr) {
	path,err := filepath.Abs("upload")
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying get absolute path for upload: %v", err))
		restErr := rest_err.NewInternalServerError("server error")
		return "", restErr
	}
	path = filepath.Join(path, fileHash) 
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", rest_err.NewNotFoundError("file not found")
		}
		logger.Error(fmt.Sprintf("Error trying get Stat: %v", err))
		return "", rest_err.NewInternalServerError("server error")
	}

	var hasTemp, hasMP4 bool
	filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "temp" {
			hasTemp = true
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".mp4" {
			hasMP4 = true
		}
		return nil
	})

	var msg string
	if hasTemp {
		val := beingProcessed[fileHash]
		if val {
			msg = BeingProcessed
		} else {			
			msg = NotSent
		}
	} else if hasMP4 {
		val := ffmpeg.GetBeingProcessed(fileHash)
		if val {
			msg = BeingProcessed
		} else {
			msg = NotBeingProcessed
		}
	} else {
		video,restErr := us.uploadRepository.GetVideoByFileHash(fileHash)
		if restErr != nil {
			return "", restErr
		}
		if (!video.Uploaded) {
			msg = NotSavedInDb
		} else {				
			msg = Processed
		}
	}

	return msg,nil
}