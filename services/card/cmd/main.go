package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/DavidRodriguez-create/pay-and-go/services/card/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/card/infrastructure"
	"github.com/DavidRodriguez-create/pay-and-go/services/card/presentation/controllers"
	"github.com/DavidRodriguez-create/pay-and-go/services/card/presentation/presenters"
	"github.com/DavidRodriguez-create/pay-and-go/services/card/presentation/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()

	// Get configuration from environment variables
	port := getEnv("PORT", "8082")
	kafkaBrokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ",")
	kafkaTopic := getEnv("KAFKA_TOPIC", "account-events")
	kafkaGroupID := getEnv("KAFKA_GROUP_ID", "card-service")

	// Initialize repositories
	cardRepo := infrastructure.NewInMemoryCardRepository()
	accountRepo := infrastructure.NewInMemoryAccountCacheRepository()

	// Initialize application services
	cardService := application.NewCardService(cardRepo, accountRepo)

	// Initialize presenter
	presenter := presenters.NewResponsePresenter()

	// Initialize controllers
	ctrls := &routes.Controllers{
		CreateCard: controllers.NewCreateCardController(cardService.CreateCard, presenter),
		GetCard:    controllers.NewGetCardController(cardService.ViewCard, presenter),
		ListCards:  controllers.NewListCardsController(cardService.ListCards, presenter),
		DeleteCard: controllers.NewDeleteCardController(cardService.DeleteCard, presenter),
	}

	// Setup routes
	mux := routes.SetupRoutes(ctrls)

	// Initialize Kafka consumer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kafkaConsumer := infrastructure.NewKafkaAccountConsumer(
		kafkaBrokers,
		kafkaTopic,
		kafkaGroupID,
		accountRepo,
	)

	// Start Kafka consumer
	if err := kafkaConsumer.Start(ctx); err != nil {
		log.Printf("Warning: Failed to start Kafka consumer: %v\n", err)
		log.Println("Service will continue without Kafka event consumption")
	} else {
		log.Println("Kafka consumer started successfully")
	}

	// Setup HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Card service starting on port %s...\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Stop Kafka consumer
	if err := kafkaConsumer.Stop(); err != nil {
		log.Printf("Error stopping Kafka consumer: %v\n", err)
	}

	// Shutdown HTTP server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exited")
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
