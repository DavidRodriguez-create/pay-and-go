package controllers

import (
	"net/http"

	"github.com/DavidRodriguez-create/pay-and-go/services/account/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/account/presentation/presenters"
)

// GetAccountController handles account retrieval requests
type GetAccountController struct {
	service application.AccountService
}

// NewGetAccountController creates a new instance
func NewGetAccountController(service application.AccountService) *GetAccountController {
	return &GetAccountController{
		service: service,
	}
}

// HandleByID processes GET /accounts?id=xxx
func (c *GetAccountController) HandleByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		presenters.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		presenters.RespondError(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	response, err := c.service.GetAccountByID(id)
	if err != nil {
		presenters.RespondError(w, err.Error(), http.StatusNotFound)
		return
	}

	presenters.RespondSuccess(w, response, http.StatusOK)
}

// HandleByAccountNumber processes GET /accounts/by-number?account_number=xxx
func (c *GetAccountController) HandleByAccountNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		presenters.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accountNumber := r.URL.Query().Get("account_number")
	if accountNumber == "" {
		presenters.RespondError(w, "Account number parameter is required", http.StatusBadRequest)
		return
	}

	response, err := c.service.GetAccountByAccountNumber(accountNumber)
	if err != nil {
		presenters.RespondError(w, err.Error(), http.StatusNotFound)
		return
	}

	presenters.RespondSuccess(w, response, http.StatusOK)
}
