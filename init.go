package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/matheuswww/mystream/src/ffmpeg"
	"github.com/matheuswww/mystream/src/logger"
)

func checkUploads() {
	dir := "upload"
	subdirs, err := os.ReadDir(dir)
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying ReadDir: %v", err))
		return
	}
	var wg sync.WaitGroup
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
				wg.Add(1)
				go func ()  {
					defer wg.Done()
					ffmpeg.SaveVideo(dir, filePath, fileHash, nil)
				}()
			}
		}
	}
	wg.Wait()
}