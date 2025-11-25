package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DavidRodriguez-create/pay-and-go/services/account/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/account/infrastructure"
	"github.com/DavidRodriguez-create/pay-and-go/services/account/presentation/controllers"
	"github.com/DavidRodriguez-create/pay-and-go/services/account/presentation/routes"
)

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Initialize repository (in-memory for now)
	repo := infrastructure.NewInMemoryAccountRepository()

	// Initialize service
	service := application.NewAccountService(repo)

	// Initialize controllers
	ctrls := &routes.Controllers{
		CreateAccount: controllers.NewCreateAccountController(service),
		GetAccount:    controllers.NewGetAccountController(service),
		ListAccounts:  controllers.NewListAccountsController(service),
		UpdateAccount: controllers.NewUpdateAccountController(service),
		DeleteAccount: controllers.NewDeleteAccountController(service),
	}

	// Setup routes
	mux := routes.SetupRoutes(ctrls)

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("🚀 Account service starting on port %s...", port)
	log.Printf("Health check: http://localhost:%s/health", port)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
