package ffmpeg

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/matheuswww/mystream/src/logger"
)

func getManifestTime(fileHash string, width, height int) (float64, error) {
	uploadPath,err := filepath.Abs("upload")
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying getLastChunk: %v", err))
		return 0,err
	}
	fp := filepath.Join(uploadPath, fileHash, fmt.Sprintf("%dp", height), fmt.Sprintf("video_%dx%d.m3u8", width, height))
	file, err := os.Open(fp)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying Open: %v", err))
		return 0,err
	}
	defer file.Close()

	var totalDuration float64
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#EXTINF:") {
			durationStr := strings.TrimSuffix(strings.TrimPrefix(line, "#EXTINF:"), ",")
			duration, err := strconv.ParseFloat(durationStr, 64)
			if err != nil {
				logger.Error(fmt.Sprintf("Error trying ParseFloat: %v", err))
				return 0,err
			}
			totalDuration += duration
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Error(fmt.Sprintf("Error reading manifest file: %v", err))
		return 0, err
	}
	return totalDuration, nil
}