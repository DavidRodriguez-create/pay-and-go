package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/DavidRodriguez-create/pay-and-go/services/account/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/account/presentation/presenters"
)

// UpdateAccountController handles account update requests
type UpdateAccountController struct {
	service application.AccountService
}

// NewUpdateAccountController creates a new instance
func NewUpdateAccountController(service application.AccountService) *UpdateAccountController {
	return &UpdateAccountController{
		service: service,
	}
}

// Handle processes PUT/PATCH /accounts
func (c *UpdateAccountController) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		presenters.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req application.UpdateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		presenters.RespondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		presenters.RespondError(w, "ID is required", http.StatusBadRequest)
		return
	}

	err := c.service.UpdateAccount(req)
	if err != nil {
		presenters.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	presenters.RespondSuccess(w, map[string]string{
		"message": "Account updated successfully",
	}, http.StatusOK)
}
