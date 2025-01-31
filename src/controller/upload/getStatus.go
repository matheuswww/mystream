package upload_controller

import (
	"net/http"

	upload_request "github.com/matheuswww/mystream/src/controller/model/upload/request"
	"github.com/matheuswww/mystream/src/logger"
	rest_err "github.com/matheuswww/mystream/src/restErr"
	"github.com/matheuswww/mystream/src/router"
)

func (uc *uploadController) GetStatus(w http.ResponseWriter, r *http.Request) {
	logger.Log("Init GetStatus")
	var getStatus upload_request.FileHash
	if err := router.BindJson(r.Body, &getStatus); err != nil {
		restErr := rest_err.NewBadRequestError("invalid fields")
		router.SendResponse(w, restErr, restErr.Code)
		return
	}
	status,restErr := uc.uploadService.GetStatus(getStatus.FileHash, beingProcessed)
	if restErr != nil {
		router.SendResponse(w, restErr, restErr.Code)
		return
	}
	router.SendResponse(w, struct{Status string}{ Status: status }, http.StatusOK)
}