package ffmpeg

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"os"
	"os/exec"
	"sync"

	"github.com/gorilla/websocket"
	upload_controller_util "github.com/matheuswww/mystream/src/controller/upload/util"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

var resolutions = []struct {
	width, height  int
	outputFile     string
	audioBitrate   string
	videoBitrate   string
}{
	{1920, 1080, "video_1080p.m3u8", "5000k", "5000k"},
	{1280, 720, "video_720p.m3u8", "2500k", "2500k"},
	{854, 480, "video_480p.m3u8", "1200k", "1200k"},
	{640, 360, "video_360p.m3u8", "800k", "800k"},
}

var numberResolutions = len(resolutions)

func SaveVideo(uploadPath, filePath, fileHash string, conn *websocket.Conn) error {
	var wg sync.WaitGroup
	cmd := exec.Command(
		"bash", "-c",
		fmt.Sprintf(
			"ffprobe -v error -select_streams v:0 -show_entries stream=r_frame_rate,duration -of default=noprint_wrappers=1 '%s'",
			filePath,
		),
	)
			
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying StdoutPipe: %v", err))
		return err
	}

	if err := cmd.Start(); err != nil {
		logger.Error(fmt.Sprintf("Error trying Start cmd: %v", err))
		return err
	}

	var frameRateStr, durationStr string
	scanner := bufio.NewScanner(outPipe)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "r_frame_rate") {
			frameRateStr = strings.Split(line, "=")[1]
		}
		if strings.HasPrefix(line, "duration") {
			durationStr = strings.Split(line, "=")[1]
		}
	}

	if frameRateStr == "" || durationStr == "" {
		err := "Error trying get frameRateStr or durationStr"
		logger.Error(err)
		return errors.New(err)
	}

	frameRateParts := strings.Split(frameRateStr, "/")
	if len(frameRateParts) != 2 {
		err := "Error trying Split frameRateStr"
		logger.Error(err)
		return errors.New(err)
	}

	numerator, err := strconv.Atoi(frameRateParts[0])
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying Atoi: %v", err))
		return err
	}

	denominator, err := strconv.Atoi(frameRateParts[1])
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying Atoi: %v", err))
		return err
	}
	frameRate := float64(numerator) / float64(denominator)
	
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying ParseFloat: %v", err))
		return err
	}		

	totalFrames := int(frameRate * duration)

	if err := cmd.Wait(); err != nil {
		logger.Error(fmt.Sprintf("Error trying Wait cmd: %v", err))
		return err
	}
	var frame int64
	var newFrame chan bool = make(chan bool)
	go func ()  {
		newFrame <- true
	}()
	for _, res := range resolutions {
		wg.Add(1)
		go func() {
			defer wg.Done()
			folder := fmt.Sprintf("%s/%s/%dp", uploadPath, fileHash, res.height)
			err := os.MkdirAll(folder, 0755)
			if err != nil {
				logger.Error(fmt.Sprintf("Error trying MkdirAll: %v", err))
				restErr := rest_err.NewInternalServerError("server error")
				upload_controller_util.SendWsRes(restErr, conn)
				conn.Close()
				return
			}

			cmd := exec.Command(
				"bash", "-c",
				fmt.Sprintf(
					"ffmpeg -loglevel error -i \"%s\" -c:v libx264 -b:v %s -vf scale=%d:%d -c:a aac -b:a %s -preset ultrafast -crf 28 -hls_time 10 -hls_list_size 0 -progress pipe:1 -hls_segment_filename \"%s/segment_%dx%d_%%03d.ts\" \"%s/video_%dx%d.m3u8\" | grep 'frame='",
					filePath,
					res.videoBitrate,
					res.width,
					res.height,
					res.audioBitrate,
					folder,
					res.width,
					res.height,
					folder,
					res.width,
					res.height,
				),
		)
		
		outPipe, err := cmd.StdoutPipe()
		if err != nil {
			logger.Error(fmt.Sprintf("Error trying StdoutPipe: %v", err))
			restErr := rest_err.NewInternalServerError("server error")
			upload_controller_util.SendWsRes(restErr, conn)
			conn.Close()
			return
		}
		
		if err := cmd.Start(); err != nil {
			logger.Error(fmt.Sprintf("Error trying Start cmd: %v", err))
			restErr := rest_err.NewInternalServerError("server error")
			upload_controller_util.SendWsRes(restErr, conn)
			conn.Close()
			return
		}

		scanner := bufio.NewScanner(outPipe)
		var count int64
		for scanner.Scan() {
			f := scanner.Text()
			f = strings.Split(f, "=")[1]
			num, err := strconv.Atoi(f)
			if err != nil {
				logger.Error(fmt.Sprintf("Error trying Atoi: %v", err))
				restErr := rest_err.NewInternalServerError("server error")
				upload_controller_util.SendWsRes(restErr, conn)
				conn.Close()
				return
			}
			atomic.AddInt64(&frame, (int64(num) - count))
			count = int64(num)
			newFrame <- true
		}
		
		if err := cmd.Wait(); err != nil {
			logger.Error(fmt.Sprintf("Error trying Wait cmd: %v", err))
			restErr := rest_err.NewInternalServerError("server error")
			upload_controller_util.SendWsRes(restErr, conn)
			conn.Close()
			return
		}
		}()
	}
	var percentage float64
	for percentage <= 100 {
		select {
			case v := <-newFrame:
			if v {
				if frame != 0 {
					percentage = (float64(frame)/(float64(int64(totalFrames))*float64(int64(numberResolutions))))*100
					upload_controller_util.SendWsRes(fmt.Sprintf("%.2f%%", (float64(frame)/(float64(int64(totalFrames))*float64(int64(numberResolutions))))*100), conn)
					if percentage >= 100 {
						percentage++
					}
				}
			}
		}
	}
	wg.Wait()
	close(newFrame)
	err = generateM3U8(fmt.Sprintf("%s/%s", uploadPath, fileHash), fileHash)
	if err != nil {
		return err
	}
	err = os.Remove(filePath)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying Remove: %v", err))
	}
	return nil
}


func generateM3U8(path, fileHash string) error {
	filePath := fmt.Sprintf("%s/master.m3u8", path)
	_, err := os.Create(filePath)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying Create file: %v", err))
		return err
	}
	host := os.Getenv("URL")
	if host == "" {
		err := errors.New("Error trying get env")
		logger.Error("Error trying get env")
		return err
	}
	urlFilePath := fmt.Sprintf("%s/file/%s",host, fileHash)
	m3u8 := `
	#EXTM3U
	#EXT-X-STREAM-INF:BANDWIDTH=800000,RESOLUTION=640x360
	`+urlFilePath+`/360p/video_640x360.m3u8
	#EXT-X-STREAM-INF:BANDWIDTH=500000,RESOLUTION=854x480
	`+urlFilePath+`/480p/video_854x480.m3u8
	#EXT-X-STREAM-INF:BANDWIDTH=1500000,RESOLUTION=1280x720
	`+urlFilePath+`/720p/video_1280x720.m3u8
	#EXT-X-STREAM-INF:BANDWIDTH=3000000,RESOLUTION=1920x1080
	`+urlFilePath+`/1080p/video_1920x1080.m3u8
	`
	err = os.WriteFile(filePath, []byte(m3u8), 0755)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying WriteFile: %v", err))
		return err
	}
	return nil
}