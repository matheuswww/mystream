package upload_service

import (
	upload_repository "github.com/matheuswww/mystream/src/model/upload/repository"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (us *uploadService) GetVideoByFileHash(fileHash string) (*upload_repository.Video ,*rest_err.RestErr) {
	return us.uploadRepository.GetVideoByFileHash(fileHash)
}