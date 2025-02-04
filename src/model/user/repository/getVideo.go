package user_repository

import (
	"context"
	"fmt"
	"time"

	user_response "github.com/matheuswww/mystream/src/controller/model/user/response"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (ur *userRepository) GetVideo(cursor string) ([]user_response.GetVideo, *rest_err.RestErr) {
	ctx,cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var args []any
	query := "SELECT id,title,description,file_hash,created_at FROM video"
	if cursor != "" {
    query += " WHERE created_at < ?"
    args = append(args, cursor)
	}
	query += " ORDER BY created_at DESC LIMIT 10" 
	rows, err := ur.sql.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying QueryContext: %v", err))
		return nil,rest_err.NewInternalServerError("server error")
	}
	var video []user_response.GetVideo
	defer rows.Close()
	for rows.Next() {
		var id, title, description, file_hash, created_at string
		err := rows.Scan(&id, &title, &description, &file_hash, &created_at)
		if err != nil {
			logger.Error(fmt.Sprintf("Error trying Scan: %v", err))
			return nil,rest_err.NewInternalServerError("server error")
		}
		video = append(video, user_response.GetVideo{
			Id: id,
			Title: title,
			Description: description,
			Cursor: created_at,
			FileHash: file_hash,
		})
	}
	if len(video) == 0 {
		return nil, rest_err.NewNotFoundError("no videos found")
	}
	return video, nil
}