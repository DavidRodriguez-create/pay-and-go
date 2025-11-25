package controllers

import (
	"net/http"

	"github.com/DavidRodriguez-create/pay-and-go/services/account/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/account/presentation/presenters"
)

// DeleteAccountController handles account deletion requests
type DeleteAccountController struct {
	service application.AccountService
}

// NewDeleteAccountController creates a new instance
func NewDeleteAccountController(service application.AccountService) *DeleteAccountController {
	return &DeleteAccountController{
		service: service,
	}
}

// Handle processes DELETE /accounts?id=xxx
func (c *DeleteAccountController) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		presenters.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		presenters.RespondError(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	err := c.service.DeleteAccount(id)
	if err != nil {
		presenters.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	presenters.RespondSuccess(w, map[string]string{
		"message": "Account deleted successfully",
	}, http.StatusOK)
}
