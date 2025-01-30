package upload_controller

import (
	"encoding/json"
	"io"
	"net/http"
	upload_request "github.com/matheuswww/mystream/src/controller/model/upload/request"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
	"github.com/matheuswww/mystream/src/router"
)

func (uc *uploadController) GetLastChunk(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		restErr := rest_err.NewInternalServerError("server error")
		router.SendResponse(w, restErr, restErr.Code)
		return
	}
	defer r.Body.Close() 
	var getLastChunkRequest upload_request.FileHash
	if err := json.Unmarshal(body, &getLastChunkRequest); err != nil {
		restErr := rest_err.NewBadRequestError("campos inválidos")
		router.SendResponse(w, restErr, restErr.Code)
		return
	}

	fileName, restErr := uc.uploadService.GetLastChunk(getLastChunkRequest)
	if restErr != nil {
		router.SendResponse(w, restErr, restErr.Code)
		return
	}
	router.SendResponse(w, struct{ Chunk string }{ Chunk: fileName }, http.StatusOK)
}

