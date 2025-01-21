package file_controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/matheuswww/mystream/src/logger"
)

func (vc *fileController) ServeFile(w http.ResponseWriter, r *http.Request) {
	logger.Log("Init ServeFile")
	absPath, err := filepath.Abs("upload")
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying get absolute path for upload: %v", err))
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	pos := strings.Index(r.URL.Path, "/file/") + 6
	path := filepath.Join(absPath, filepath.Clean(r.URL.Path[pos:]))
	if !strings.HasPrefix(path, absPath+string(os.PathSeparator)) {
    http.Error(w, "File not found", http.StatusNotFound)
    return
	}
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		logger.Error(fmt.Sprintf("Error trying Open file: %v", err))
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	b, err := file.Stat()
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying get file stat: %v", err))
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, path, b.ModTime(), file)
}
