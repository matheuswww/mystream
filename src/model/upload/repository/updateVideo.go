package upload_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (ur *uploadRepository) UpdateVideo(fileHash, title, description string, uploaded *bool) *rest_err.RestErr {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, restErr  := ur.GetVideoByFileHash(fileHash)
	if restErr != nil {
		return restErr
	}
	var args []any
	query := "UPDATE video SET "
	if title != "" {
		args = append(args, title)
		query += fmt.Sprintf("title = $%d, ", len(args))
	}
	if description != "" {
		args = append(args, description)
		query += fmt.Sprintf("description = $%d, ", len(args))
	}
	if uploaded != nil {
		args = append(args, *uploaded)
		query += fmt.Sprintf("uploaded = $%d, ", len(args))
	}
	if len(args) ==  0 {
		return rest_err.NewBadRequestError("no params")
	}
	args = append(args, fileHash)
	query = query[:len(query) - 2]
	query += fmt.Sprintf(" WHERE file_hash = $%d", len(args))
	_,err := ur.sql.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying ExecContext: %v", err))
		return rest_err.NewInternalServerError("server error")
	}
	return nil
}