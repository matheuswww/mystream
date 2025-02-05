package upload_service

import (
	"github.com/gorilla/websocket"
	admin_response "github.com/matheuswww/mystream/src/controller/model/admin/response"
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
	UploadFile(conn *websocket.Conn, uploadFile upload_request.UploadFile, file_hash string)
	GetLastChunk(getLastChunkRequest upload_request.FileHash) (string, *rest_err.RestErr)
	GetFfmpegProgress(fileHash string, conn *websocket.Conn)
	RetryFfmpeg(fileHash string) *rest_err.RestErr
	GetStatus(fileHash string, beingProcessed map[string]bool) (string, *rest_err.RestErr)
	CheckToken(token string) bool
	InsertVideo(title, description, fileHash string) *rest_err.RestErr
	GetVideoByFileHash(fileHash string) (*upload_repository.Video ,*rest_err.RestErr)
	UpdateVideo(id, title, description string, uploaded *bool, beingProcessed map[string]bool) *rest_err.RestErr
	GetVideo(cursor string) ([]admin_response.GetVideo, *rest_err.RestErr)
}
