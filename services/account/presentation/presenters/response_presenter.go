package presenters

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// RespondError sends an error response
func RespondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	
	json.NewEncoder(w).Encode(response)
}

// RespondSuccess sends a success response
func RespondSuccess(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
