package upload_service

import upload_repository "github.com/matheuswww/mystream/src/model/upload/repository"

func NewUploadService(uploadRepository upload_repository.UploadRepository) UploadService {
	return &uploadService {
		uploadRepository,
	}
}

type uploadService struct {
	uploadRepository upload_repository.UploadRepository
}

type UploadService interface{}
