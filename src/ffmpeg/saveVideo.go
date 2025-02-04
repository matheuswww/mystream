package ffmpeg

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
)

var conns map[string]*websocket.Conn = make(map[string]*websocket.Conn)
var beingProcessed map[string]bool = make(map[string]bool)

var chunkTime = 10

func GetBeingProcessed(fileHash string) bool {
	_,f := beingProcessed[fileHash]
	return f
}

func UpdateConn(fileHash string, conn *websocket.Conn) error {
	if _,f := conns[fileHash]; f {
		conns[fileHash] = conn
		return nil
	}
	return errors.New("file not found")
}

func removeBeingProcessed(filehash string) {
	delete(beingProcessed, filehash)
}

func SaveVideo(uploadPath, filePath, fileHash string, conn *websocket.Conn) error {
	defer removeBeingProcessed(fileHash)
	conns[fileHash] = conn
	beingProcessed[fileHash] = true
	chunks, frames, frameRate, duration, err := getInfos(fileHash, filePath, conn)
	if err != nil {
		return err
	}
	var totalFrames = frames
	var wg sync.WaitGroup
	var indexs map[int]int = make(map[int]int)
	var i int
	var processing int = numberResolutions
	var frame, speedCounter, fpsCounter int64
	var fps, speed, percentage float64
	var allFps, allSpeed float64
	var newFrame chan bool = make(chan bool)
	var mu sync.Mutex
	go func ()  {
		newFrame <- true
	}()
	for index, res := range resolutions {
		wg.Add(1)
		go func() {
			defer wg.Done()
			folder := fmt.Sprintf("%s/%s/%dp", uploadPath, fileHash, res.height)
			err := os.MkdirAll(folder, 0755)
			if err != nil {
				logger.Error(fmt.Sprintf("Error trying MkdirAll: %v", err))
				restErr := rest_err.NewInternalServerError("server error")
				sendWsRes(restErr, fileHash)
				close(newFrame)
				if _,f := conns[fileHash]; f {
					conns[fileHash].Close()
				}
				delete(conns, fileHash)
				return
			}
			var chunkStart int
			var f bool
			var timeStart float64
			if chunkStart,f = chunks[res.height]; f {
				if chunkStart > 0 && chunkStart >= ((int(duration) / numberResolutions) / chunkTime) {
					mu.Lock()
					processing = processing - 1
					frames = frames - (totalFrames/numberResolutions)
					if processing == 0 {
						percentage = 101
						newFrame <- false
					}
					mu.Unlock()
					return
				}
				mu.Lock()
				indexs[index] = i
				i++
				if chunkStart != 0 {
					timeStart,err = getManifestTime(fileHash, res.width, res.height)
					if err != nil {
						restErr := rest_err.NewInternalServerError("server error")
						sendWsRes(restErr, fileHash)
						close(newFrame)
						if _,f := conns[fileHash]; f {
							conns[fileHash].Close()
						}
						delete(conns, fileHash)
						return
					}
				}
				frames = int(float64(frames) - float64(frameRate) * timeStart)
				mu.Unlock()
			}
			cmd := exec.Command(
				"bash", "-c",
				fmt.Sprintf(
					"ffmpeg -loglevel error -ss %f -i \"%s\" -segment_start_number %d -c:v libx264 -b:v %s -vf scale=%d:%d -c:a aac -b:a %s -preset ultrafast -crf 28 -hls_time %d -hls_list_size 0 -hls_flags append_list -progress pipe:1 -hls_segment_filename \"%s/segment_%dx%d_%%03d.ts\" \"%s/video_%dx%d.m3u8\"",
					timeStart,
					filePath,
					chunkStart,
					res.videoBitrate,
					res.width,
					res.height,
					res.audioBitrate,
					chunkTime,
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
			sendWsRes(restErr, fileHash)
			close(newFrame)
			if _,f := conns[fileHash]; f {
				conns[fileHash].Close()
			}
			delete(conns, fileHash)
			return
		}
		
		if err := cmd.Start(); err != nil {
			logger.Error(fmt.Sprintf("Error trying Start cmd: %v", err))
			restErr := rest_err.NewInternalServerError("server error")
			sendWsRes(restErr, fileHash)
			close(newFrame)
			if _,f := conns[fileHash]; f {
				conns[fileHash].Close()
			}
			delete(conns, fileHash)
			return
		}

		scanner := bufio.NewScanner(outPipe)
		var lastFrame int64
		for scanner.Scan() {
			str := scanner.Text()
			if strings.HasPrefix(str, "frame=") {
				f := strings.Split(str, "=")[1]
				num, err := strconv.Atoi(f)
				if err != nil {
					logger.Error(fmt.Sprintf("Error trying Atoi: %v", err))
					continue
				}
				atomic.AddInt64(&frame, (int64(num) - lastFrame))
				lastFrame = int64(num)
				newFrame <- true
			}
			if strings.HasPrefix(str, "speed=") && strings.Contains(str, "x") {
				str = strings.ReplaceAll(str, " ", "")
				s := strings.Replace(strings.Split(str, "=")[1], "x", "", 1);
				num, err := strconv.ParseFloat(s, 64)
				if err != nil {
					logger.Error(fmt.Sprintf("Error trying ParseFloat: %v", err))
					continue
				}
				mu.Lock()
				if speedCounter == int64(indexs[index]) {
					speed += num
					if (processing - 1) > 0 {
						speedCounter++
					}
				}
				if speedCounter == int64(processing - 1) {
					allSpeed = speed
					speedCounter = 0
					speed = 0
				}
				mu.Unlock()
			}
			if strings.HasPrefix(str, "fps=") {
				f := strings.Split(str, "=")[1]
				num, err := strconv.ParseFloat(f, 64)
				if err != nil {
					logger.Error(fmt.Sprintf("Error trying ParseFloat: %v", err))
					continue
				}
				mu.Lock()
				if fpsCounter == int64(indexs[index]) {
					fps += num
					if (processing - 1) > 0 {
						fpsCounter++
					}
				}
				if fpsCounter == int64(processing - 1) {
					allFps = fps
					fpsCounter = 0
					fps = 0
				}
				mu.Unlock()
			}
		}
		
		if err := cmd.Wait(); err != nil {
			logger.Error(fmt.Sprintf("Error trying Wait cmd: %v", err))
			restErr := rest_err.NewInternalServerError("server error")
			sendWsRes(restErr, fileHash)
			close(newFrame)
			if _,f := conns[fileHash]; f {
				conns[fileHash].Close()
			}
			delete(conns, fileHash)
			return
		}
		}()
	}

	for percentage <= 100 {
		select {
			case v := <-newFrame:
				sendProgress(v, frames, frame, allSpeed, allFps, fileHash, &percentage)
		}
	}

	wg.Wait()
	close(newFrame)
	err = os.Remove(filePath)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying Remove: %v", err))
	}
	res := struct{ Message string }{
		Message: "sucesso",
	}
	sendWsRes(res, fileHash)
	conn.Close()
	delete(conns, fileHash)
	return nil
}

func sendProgress(v bool, totalFrames int, frame int64, allSpeed, allFps float64, fileHash string, percentage *float64) {
	if v {
		if frame != 0 {
			*percentage = (float64(frame)/(float64(totalFrames)))*100
			var timeExpected float64
			if allSpeed != 0 && allFps != 0 {
				timeExpected = ((float64(totalFrames) - float64(frame))/allFps) * 1 / allSpeed
			} else {
				timeExpected = 0
			}
			var formatedTimeExpected string
			if timeExpected > 60 {
				formatedTimeExpected = fmt.Sprintf("%.2f min", timeExpected / 60)
			} else {
				formatedTimeExpected = fmt.Sprintf("%.2f seg", timeExpected)
			}
			res := struct { Percentage string; TimeExpected string }{Percentage: fmt.Sprintf("%.2f%%", *percentage), TimeExpected: formatedTimeExpected}
			sendWsRes(res, fileHash)
			if *percentage >= 100 {
				*percentage++
			}
		}
	}
}

func getInfos(fileHash, filePath string, conn *websocket.Conn) (map[int]int, int, int, float64, error) {
	conns[fileHash] = conn
	chunks,err := getLastChunk(fileHash)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying getLastChunk: %v", err))
		return nil, 0, 0, 0,  err
	}
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
		return nil, 0, 0, 0, err
	}

	if err := cmd.Start(); err != nil {
		logger.Error(fmt.Sprintf("Error trying Start cmd: %v", err))
		return nil, 0, 0, 0, err
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
		return nil, 0, 0, 0, errors.New(err)
	}

	frameRateParts := strings.Split(frameRateStr, "/")
	if len(frameRateParts) != 2 {
		err := "Error trying Split frameRateStr"
		logger.Error(err)
		return nil, 0, 0, 0, errors.New(err)
	}

	numerator, err := strconv.Atoi(frameRateParts[0])
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying Atoi: %v", err))
		return nil, 0, 0, 0, err
	}

	denominator, err := strconv.Atoi(frameRateParts[1])
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying Atoi: %v", err))
		return nil, 0, 0, 0, err
	}
	frameRate := float64(numerator) / float64(denominator)
	
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying ParseFloat: %v", err))
		return nil, 0, 0, 0, err
	}		

	totalFrames := int(frameRate * duration)

	if err := cmd.Wait(); err != nil {
		logger.Error(fmt.Sprintf("Error trying Wait cmd: %v", err))
		return nil, 0, 0, 0, err
	}
	return chunks, totalFrames*numberResolutions, int(frameRate), duration*float64(numberResolutions), nil
}