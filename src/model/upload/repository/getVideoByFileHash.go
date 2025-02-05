package upload_repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

type Video struct {
	Id          string
	Title 			string
	Description string
	FileHash    string
	Uploaded    bool
}

func (ur *uploadRepository) GetVideoByFileHash(fileHash string) (*Video, *rest_err.RestErr) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	
	query := "SELECT id, title, description, uploaded FROM video WHERE file_hash = $1"
	var id, title, description string
	var uploaded bool
	err := ur.sql.QueryRowContext(ctx, query, fileHash).Scan(&id, &title, &description, &uploaded)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, rest_err.NewNotFoundError("video not found")
		}
		logger.Error(fmt.Sprintf("Error trying QueryRowContext: %v", err))
		return nil, rest_err.NewInternalServerError("sever error")
	}
	return &Video{
		id,
		title,
		description,
		fileHash,
		uploaded,
	},nil
}