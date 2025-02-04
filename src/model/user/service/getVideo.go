package user_service

import (
	"fmt"
	"os"

	user_response "github.com/matheuswww/mystream/src/controller/model/user/response"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

func (us *userService) GetVideo(cursor string) ([]user_response.GetVideo, *rest_err.RestErr) {
	videos, restErr := us.user_repository.GetVideo(cursor)
	if restErr != nil {
		return nil, restErr
	}
	url := os.Getenv("FULL_URL")
	if url == "" {
		logger.Error("Error trying get env")
		return nil,rest_err.NewInternalServerError("server error")
	}
	for i := 0; i < len(videos); i++ {
		videos[i].Url.FullHd = fmt.Sprintf("%s/%s/1080p/video_1920x1080.m3u8", url, videos[i].FileHash)
		videos[i].Url.Hd = fmt.Sprintf("%s/%s/720p/video_1280x720.m3u8", url, videos[i].FileHash)
		videos[i].Url.Sd = fmt.Sprintf("%s/%s/360p/video_854x480.m3u8", url, videos[i].FileHash)
		videos[i].Url.Vd = fmt.Sprintf("%s/%s/180p/video_640x360.m3u8", url, videos[i].FileHash)
	}
	return videos,nil
}