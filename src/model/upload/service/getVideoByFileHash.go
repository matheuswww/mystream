package upload_service

import rest_err "github.com/matheuswww/mystream/src/restErr"

func (us *uploadService) GetVideoByFileHash(fileHash string) *rest_err.RestErr {
	return us.uploadRepository.GetVideoByFileHash(fileHash)
}