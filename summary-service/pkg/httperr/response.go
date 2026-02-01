package httperr

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func HandleError(w http.ResponseWriter, statusCode int, message string) {
	RespondJSON(w, statusCode, APIError{
		Message:    message,
		StatusCode: statusCode,
	})
}

func RespondJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
