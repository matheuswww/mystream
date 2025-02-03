package upload_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (ur *uploadRepository) InsertVideo(title, description, fileHash string) *rest_err.RestErr {
	ctx,cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	id := uuid.NewString()
	query := "INSERT INTO video (id, title, description, file_hash) VALUES ($1, $2, $3, $4)"
	_,err := ur.sql.ExecContext(ctx, query, id, title, description, fileHash)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying ExecContext: %v", err))
		return rest_err.NewInternalServerError("server error")
	}
	return nil
}