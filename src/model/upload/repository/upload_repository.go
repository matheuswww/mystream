package upload_repository

import (
	"database/sql"

	admin_response "github.com/matheuswww/mystream/src/controller/model/admin/response"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func NewUploadRepository(db *sql.DB) UploadRepository {
	return &uploadRepository{
		db,
	}
}

type UploadRepository interface {
	InsertVideo(title, description, fileHash string) *rest_err.RestErr
	GetVideoByFileHash(fileHash string) (*Video, *rest_err.RestErr)
	UpdateVideo(id, title, description string, uploaded *bool) *rest_err.RestErr
	GetVideo(cursor string) ([]admin_response.GetVideo, *rest_err.RestErr)
}

type uploadRepository struct {
	sql *sql.DB
}