package presenters

import (
	"encoding/json"
	"net/http"

	"github.com/DavidRodriguez-create/pay-and-go/services/card/domain"
)

// ResponsePresenter handles HTTP response formatting
type ResponsePresenter struct{}

// NewResponsePresenter creates a new ResponsePresenter
func NewResponsePresenter() *ResponsePresenter {
	return &ResponsePresenter{}
}

// Success writes a successful JSON response
func (p *ResponsePresenter) Success(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Error writes an error JSON response
func (p *ResponsePresenter) Error(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// HandleError maps domain errors to HTTP responses
func (p *ResponsePresenter) HandleError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrCardIDRequired, domain.ErrCardNumberRequired,
		domain.ErrCountryRequired, domain.ErrAccountIDRequired:
		p.Error(w, err.Error(), http.StatusBadRequest)
	case domain.ErrCardNotFound, domain.ErrAccountNotFound, domain.ErrAccountCacheNotFound:
		p.Error(w, err.Error(), http.StatusNotFound)
	case domain.ErrCardAlreadyDeleted:
		p.Error(w, err.Error(), http.StatusConflict)
	case domain.ErrAccountDeleted, domain.ErrAccountInactive:
		p.Error(w, err.Error(), http.StatusForbidden)
	default:
		p.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
