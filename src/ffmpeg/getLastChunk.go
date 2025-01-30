package ffmpeg

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/matheuswww/mystream/src/logger"
)

func getLastChunk(fileHash string) (map[int]int, error) {
	uploadPath, err := filepath.Abs("upload")
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying getLastChunk: %v", err))
		return nil, err
	}
	folderPath := filepath.Join(uploadPath, fileHash)

	chunks := make(map[int]int)
	for _, v := range resolutions {
		filePath := filepath.Join(folderPath, fmt.Sprintf("%dp", v.height))
		var num int
		err = filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				if os.IsNotExist(err) {
					return nil
				}
				return err
			}
			if strings.HasSuffix(info.Name(), ".ts") {
				n, err := getNumber(info.Name())
				if err != nil {
					return err
				}
				if n > num {
					num = n
				}
			}
			return nil
		})

		chunks[v.height] = num
	}
	if err != nil {
		return nil, err
	}
	return chunks, nil
}

func getNumber(filePath string) (int, error) {
	filePath = strings.Split(filePath, ".ts")[0]
	lastUnderscore := strings.LastIndex(filePath, "_")
	if lastUnderscore == -1 {
		err := fmt.Errorf("no underscore found in %s", filePath)
		logger.Error(err.Error())
		return 0, err
	}
	numberPart := filePath[lastUnderscore+1:]
	n, err := strconv.Atoi(numberPart)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying Atoi: %v", err))
		return 0, err
	}
	return n, nil
}