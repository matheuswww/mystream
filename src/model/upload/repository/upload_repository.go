package upload_repository

import (
	"database/sql"
)

func NewUploadRepository(db *sql.DB) UploadRepository {
	return &uploadRepository{
		db,
	}
}

type UploadRepository interface{}

type uploadRepository struct {
	db *sql.DB
}