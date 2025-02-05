package upload_service

import (
	"fmt"

	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (us *uploadService) UpdateVideo(fileHash, title, description string, uploaded *bool, beingProcessed map[string]bool) *rest_err.RestErr {
	video, restErr := us.uploadRepository.GetVideoByFileHash(fileHash)
	if restErr != nil {
		return restErr
	}
	status,restErr := us.GetStatus(fileHash, beingProcessed)
	if restErr != nil {
		return restErr
	}
	if status != Processed && status != NotSavedInDb {
		return rest_err.NewBadRequestError(fmt.Sprintf("it is not possible to change the video in this status: %s", status))
	} 
	if video.Uploaded && uploaded != nil {
		return rest_err.NewBadRequestError("you can't change uploaded")
	}
	return us.uploadRepository.UpdateVideo(fileHash, title, description, uploaded)
}