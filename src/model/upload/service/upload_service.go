package upload_service

import (
	"github.com/gorilla/websocket"
	upload_request "github.com/matheuswww/mystream/src/controller/model/upload/request"
	upload_repository "github.com/matheuswww/mystream/src/model/upload/repository"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func NewUploadService(uploadRepository upload_repository.UploadRepository) UploadService {
	return &uploadService {
		uploadRepository,
	}
}

type uploadService struct {
	uploadRepository upload_repository.UploadRepository
}

type UploadService interface {
	UploadFile(conn *websocket.Conn, uploadFile upload_request.UploadFile)
	GetLastChunk(getLastChunkRequest upload_request.FileHash) (string, *rest_err.RestErr)
	GetFfmpegProgress(fileHash string, conn *websocket.Conn)
}
