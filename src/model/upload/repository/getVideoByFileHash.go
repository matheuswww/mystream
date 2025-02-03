package upload_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (ur *uploadRepository) GetVideoByFileHash(fileHash string) *rest_err.RestErr {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	query := "SELECT COUNT(*) FROM video WHERE file_hash = $1"
	var count int
	err := ur.sql.QueryRowContext(ctx, query, fileHash).Scan(&count)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying QueryRowContext: %v", err))
		return nil
	}
	if count != 0 {
		return rest_err.NewNotFoundError("video not found")
	}
	return nil
}