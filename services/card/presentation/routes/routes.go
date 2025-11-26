package routes

import (
	"net/http"

	"github.com/DavidRodriguez-create/pay-and-go/services/card/presentation/controllers"
)

// Controllers holds all controller instances
type Controllers struct {
	CreateCard *controllers.CreateCardController
	GetCard    *controllers.GetCardController
	ListCards  *controllers.ListCardsController
	DeleteCard *controllers.DeleteCardController
}

// SetupRoutes configures all HTTP routes for the card service
func SetupRoutes(ctrls *Controllers) *http.ServeMux {
	mux := http.NewServeMux()

	// Collection endpoint (plural) - list all cards
	// GET /cards - List all cards
	mux.HandleFunc("/cards", handleCardList(ctrls))

	// Search endpoints
	// GET /cards/by-number?card_number=xxx - Get card by card number
	mux.HandleFunc("/cards/by-number", handleCardByNumber(ctrls))
	// GET /cards/by-account?account_id=xxx - Get cards by account ID
	mux.HandleFunc("/cards/by-account", handleCardsByAccount(ctrls))

	// Single resource endpoint (singular) - operates on ONE card
	// POST /card - Create a new card
	// GET /card?id=xxx - Get card by ID
	// DELETE /card?id=xxx - Delete card by ID
	mux.HandleFunc("/card", handleCard(ctrls))

	// Health check endpoint - GET /health
	mux.HandleFunc("/health", handleHealth())

	return mux
}

// handleCardList handles listing all cards
func handleCardList(ctrls *Controllers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		ctrls.ListCards.Handle(w, r)
	}
}

// handleCard handles operations on a single card resource
func handleCard(ctrls *Controllers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// POST /card - Create new card (no ID needed)
		if r.Method == http.MethodPost {
			ctrls.CreateCard.Handle(w, r)
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
			ctrls.GetCard.HandleByID(w, r)
		case http.MethodDelete:
			ctrls.DeleteCard.Handle(w, r)
		default:
			w.Header().Set("Allow", "POST, GET, DELETE")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// handleCardByNumber handles searching for a card by card number
func handleCardByNumber(ctrls *Controllers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		ctrls.GetCard.HandleByCardNumber(w, r)
	}
}

// handleCardsByAccount handles retrieving cards by account ID
func handleCardsByAccount(ctrls *Controllers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		ctrls.GetCard.HandleByAccountID(w, r)
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
		w.Write([]byte(`{"status":"healthy","service":"card-service"}`))
	}
}
