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

func ExtractID(r *http.Request, prefix string) string {
	return strings.TrimPrefix(r.URL.Path, prefix)
}
