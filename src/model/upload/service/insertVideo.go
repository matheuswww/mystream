package upload_service

import rest_err "github.com/matheuswww/mystream/src/restErr"

func (us *uploadService) InsertVideo(title, description, fileHash string) *rest_err.RestErr {
	return us.uploadRepository.InsertVideo(title, description, fileHash)
}