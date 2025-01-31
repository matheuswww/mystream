package upload_controller

import (
	"net/http"
	upload_request "github.com/matheuswww/mystream/src/controller/model/upload/request"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
	"github.com/matheuswww/mystream/src/router"
)

func (uc *uploadController) GetLastChunk(w http.ResponseWriter, r *http.Request) {
	logger.Log("Init GetLastChunk")
	var getLastChunkRequest upload_request.FileHash
	if err := router.BindJson(r.Body, &getLastChunkRequest); err != nil {
		restErr := rest_err.NewBadRequestError("invalid fields")
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

