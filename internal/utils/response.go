package utils

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alvindashahrul/my-app/internal/api"
)

func JSONResponse(w http.ResponseWriter, status int, message string, data interface{}, meta *api.Metadata) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := api.Response{
		Status:   status,
		Success:  status < 400,
		Message:  message,
		Data:     data,
		Metadata: meta,
	}

	json.NewEncoder(w).Encode(res)
}

func SendSuccessResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	JSONResponse(w, status, message, data, nil)
}

func SendErrorResponse(w http.ResponseWriter, status int, message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg = message + ": " + err.Error()
	}
	JSONResponse(w, status, errorMsg, nil, nil)
}

func SendPaginatedResponse(w http.ResponseWriter, status int, message string, data interface{}, page, limit, total int) {
	meta := &api.Metadata{
		Page:  page,
		Limit: limit,
		Total: total,
	}
	JSONResponse(w, status, message, data, meta)
}

func ExtractID(r *http.Request, prefix string) string {
	return strings.TrimPrefix(r.URL.Path, prefix)
}
