package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/matheuswww/mystream/src/ffmpeg"
	"github.com/matheuswww/mystream/src/logger"
)

func init() {
	dir := "upload"
	subdirs, err := os.ReadDir(dir)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying ReadDir: %v", err))
		return
	}

	for _, entry := range subdirs {
		if entry.IsDir() {
			subdirPath := filepath.Join(dir, entry.Name())
			var fileHash, filePath string
			filepath.WalkDir(subdirPath, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					logger.Error(fmt.Sprintf("Error trying WalkDir: %v", err))
					return err
				}
				if !d.IsDir() && filepath.Ext(d.Name()) == ".mp4" {
					filePath = path
					fileHash = filepath.Base(filepath.Dir(path))
					return fs.SkipDir
				}
				return nil
			})
			if fileHash != "" && filePath != "" {
				ffmpeg.SaveVideo(dir, filePath, fileHash, nil)
			}
		}
	}
}