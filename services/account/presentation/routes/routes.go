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

// corsMiddleware adds CORS headers to allow browser requests
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// SetupRoutes configures all HTTP routes for the account service
func SetupRoutes(ctrls *Controllers) *http.ServeMux {
	mux := http.NewServeMux()

	// Collection endpoint (plural) - list all accounts
	// GET /accounts - List all accounts
	mux.HandleFunc("/accounts", corsMiddleware(handleAccountList(ctrls)))

	// Search endpoint - GET /accounts/by-number?account_number=xxx
	mux.HandleFunc("/accounts/by-number", corsMiddleware(handleAccountByNumber(ctrls)))

	// Single resource endpoint (singular) - operates on ONE account
	// POST /account - Create a new account
	// GET /account?id=xxx - Get account by ID
	// PUT /account?id=xxx - Update account by ID
	// DELETE /account?id=xxx - Delete account by ID
	mux.HandleFunc("/account", corsMiddleware(handleAccount(ctrls)))

	// Health check endpoint - GET /health
	mux.HandleFunc("/health", corsMiddleware(handleHealth()))

	return mux
}

// handleAccountList handles listing all accounts
func handleAccountList(ctrls *Controllers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		ctrls.ListAccounts.Handle(w, r)
	}
}

// handleAccount handles operations on a single account resource
func handleAccount(ctrls *Controllers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// POST /account - Create new account (no ID needed)
		if r.Method == http.MethodPost {
			ctrls.CreateAccount.Handle(w, r)
			return
		}

		// All other operations require an ID
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing required query parameter: id", http.StatusBadRequest)
			return
		}

		// Route based on HTTP method
		switch r.Method {
		case http.MethodGet:
			ctrls.GetAccount.HandleByID(w, r)
		case http.MethodPut, http.MethodPatch:
			ctrls.UpdateAccount.Handle(w, r)
		case http.MethodDelete:
			ctrls.DeleteAccount.Handle(w, r)
		default:
			w.Header().Set("Allow", "POST, GET, PUT, PATCH, DELETE")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// handleAccountByNumber handles searching for an account by account number
func handleAccountByNumber(ctrls *Controllers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		ctrls.GetAccount.HandleByAccountNumber(w, r)
	}
}

// handleHealth returns the health status of the service
func handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"account-service"}`))
	}
}
