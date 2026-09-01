package api

import (
	"encoding/json"
	"net/http"

	"github.com/dawsonyoung/linden/errs"
)

type errorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func writeError(w http.ResponseWriter, err error) {
	var statusCode int
	var codeStr string

	code := errs.CodeOf(err)
	switch code {
	case errs.InvalidArgument:
		statusCode = http.StatusBadRequest
		codeStr = "invalid_argument"
	case errs.NotFound:
		statusCode = http.StatusNotFound
		codeStr = "not_found"
	case errs.Unavailable:
		statusCode = http.StatusServiceUnavailable
		codeStr = "unavailable"
	case errs.DeadlineExceeded:
		statusCode = http.StatusServiceUnavailable // API maps deadline to 503
		codeStr = "deadline_exceeded"
	default:
		statusCode = http.StatusInternalServerError
		codeStr = "internal"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse{
		Error: err.Error(),
		Code:  codeStr,
	})
}
