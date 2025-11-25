package controllers

import (
	"net/http"

	"github.com/DavidRodriguez-create/pay-and-go/services/account/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/account/presentation/presenters"
)

// ListAccountsController handles account listing requests
type ListAccountsController struct {
	service application.AccountService
}

// NewListAccountsController creates a new instance
func NewListAccountsController(service application.AccountService) *ListAccountsController {
	return &ListAccountsController{
		service: service,
	}
}

// Handle processes GET /accounts/list
func (c *ListAccountsController) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		presenters.RespondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response, err := c.service.ListAccounts()
	if err != nil {
		presenters.RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	presenters.RespondSuccess(w, response, http.StatusOK)
}
