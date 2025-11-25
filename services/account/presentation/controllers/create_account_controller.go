package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/DavidRodriguez-create/pay-and-go/services/account/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/account/presentation/presenters"
)

// CreateAccountController handles account creation requests
type CreateAccountController struct {
	service application.AccountService
}

// NewCreateAccountController creates a new instance
func NewCreateAccountController(service application.AccountService) *CreateAccountController {
	return &CreateAccountController{
		service: service,
	}
}

// Handle processes POST /accounts
func (c *CreateAccountController) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		presenters.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req application.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		presenters.RespondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := c.service.CreateAccount(req)
	if err != nil {
		presenters.RespondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	presenters.RespondSuccess(w, response, http.StatusCreated)
}
