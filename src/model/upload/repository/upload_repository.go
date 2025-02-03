package upload_repository

import (
	"database/sql"

	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func NewUploadRepository(db *sql.DB) UploadRepository {
	return &uploadRepository{
		db,
	}
}

type UploadRepository interface {
	InsertVideo(title, description, fileHash string) *rest_err.RestErr
	GetVideoByFileHash(fileHash string) *rest_err.RestErr
}

type uploadRepository struct {
	sql *sql.DB
}