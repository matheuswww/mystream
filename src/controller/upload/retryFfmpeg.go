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

func (uc *uploadController) RetryFfmpeg(w http.ResponseWriter, r *http.Request) {
	logger.Log("Init RetryFfmpeg")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		restErr := rest_err.NewInternalServerError("server error")
		router.SendResponse(w, restErr, restErr.Code)
		return
	}
	defer r.Body.Close() 
	var retryFfmpeg upload_request.FileHash
	if err := json.Unmarshal(body, &retryFfmpeg); err != nil {
		restErr := rest_err.NewBadRequestError("invalid fields")
		router.SendResponse(w, restErr, restErr.Code)
		return
	}
	restErr := uc.uploadService.RetryFfmpeg(retryFfmpeg.FileHash)
	if restErr != nil {
		router.SendResponse(w, restErr, restErr.Code)
		return
	}
	router.SendResponse(w, struct{ Message string }{ Message: "success" }, http.StatusOK)
}