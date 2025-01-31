package upload_controller

import (
	"net/http"

	upload_request "github.com/matheuswww/mystream/src/controller/model/upload/request"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
	"github.com/matheuswww/mystream/src/router"
)

func (uc *uploadController) RetryFfmpeg(w http.ResponseWriter, r *http.Request) {
	logger.Log("Init RetryFfmpeg")
	var retryFfmpeg upload_request.FileHash
	if err := router.BindJson(r.Body, &retryFfmpeg); err != nil {
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