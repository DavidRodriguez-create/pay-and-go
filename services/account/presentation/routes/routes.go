package routes

import (
	"net/http"

	"github.com/DavidRodriguez-create/pay-and-go/services/account/presentation/controllers"
)

// Controllers holds all controller instances
type Controllers struct {
	CreateAccount *controllers.CreateAccountController
	GetAccount    *controllers.GetAccountController
	ListAccounts  *controllers.ListAccountsController
	UpdateAccount *controllers.UpdateAccountController
	DeleteAccount *controllers.DeleteAccountController
}

// SetupRoutes configures all HTTP routes for the account service
func SetupRoutes(ctrls *Controllers) *http.ServeMux {
	mux := http.NewServeMux()

	// POST /accounts - Create a new account
	mux.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			ctrls.CreateAccount.Handle(w, r)
			return
		}
		if r.Method == http.MethodGet {
			ctrls.ListAccounts.Handle(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// Operations on specific accounts via query params
	// GET /accounts/?id=xxx - Get account by ID
	// PUT /accounts/?id=xxx - Update account by ID
	// DELETE /accounts/?id=xxx - Delete account by ID
	mux.HandleFunc("/accounts/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		if r.Method == http.MethodGet && id != "" {
			ctrls.GetAccount.HandleByID(w, r)
			return
		}
		if (r.Method == http.MethodPut || r.Method == http.MethodPatch) && id != "" {
			ctrls.UpdateAccount.Handle(w, r)
			return
		}
		if r.Method == http.MethodDelete && id != "" {
			ctrls.DeleteAccount.Handle(w, r)
			return
		}
		http.Error(w, "Method not allowed or missing ID", http.StatusBadRequest)
	})

	// GET /accounts/by-number?account_number=xxx - Search by account number
	mux.HandleFunc("/accounts/by-number", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			ctrls.GetAccount.HandleByAccountNumber(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// GET /health - Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"account-service"}`))
	})

	return mux
}
