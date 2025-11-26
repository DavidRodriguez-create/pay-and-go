package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/DavidRodriguez-create/pay-and-go/services/account/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/account/domain"
	"github.com/DavidRodriguez-create/pay-and-go/services/account/infrastructure"
	"github.com/DavidRodriguez-create/pay-and-go/services/account/presentation/controllers"
	"github.com/DavidRodriguez-create/pay-and-go/services/account/presentation/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (optional)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Initialize repository (in-memory for now)
	repo := infrastructure.NewInMemoryAccountRepository()

	// Initialize Kafka producer (optional)
	var eventPublisher domain.EventPublisher
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")

	if kafkaBrokers != "" && kafkaTopic != "" {
		brokers := strings.Split(kafkaBrokers, ",")
		kafkaProducer := infrastructure.NewKafkaProducer(brokers, kafkaTopic)
		eventPublisher = kafkaProducer
		log.Printf("✅ Kafka producer initialized (brokers: %s, topic: %s)", kafkaBrokers, kafkaTopic)

		// Ensure graceful shutdown of Kafka producer
		defer func() {
			if err := kafkaProducer.Close(); err != nil {
				log.Printf("Error closing Kafka producer: %v", err)
			}
		}()
	} else {
		log.Println("⚠️  Kafka not configured - events will not be published")
		log.Println("   Set KAFKA_BROKERS and KAFKA_TOPIC environment variables to enable event publishing")
	}

	// Initialize service
	service := application.NewAccountService(repo, eventPublisher)

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
