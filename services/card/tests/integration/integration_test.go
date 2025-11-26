package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DavidRodriguez-create/pay-and-go/services/card/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/card/domain"
	"github.com/DavidRodriguez-create/pay-and-go/services/card/infrastructure"
	"github.com/DavidRodriguez-create/pay-and-go/services/card/presentation/controllers"
	"github.com/DavidRodriguez-create/pay-and-go/services/card/presentation/presenters"
	"github.com/DavidRodriguez-create/pay-and-go/services/card/presentation/routes"
)

func setupTestServer() (*httptest.Server, *infrastructure.InMemoryCardRepository, *infrastructure.InMemoryAccountCacheRepository) {
	// Setup repositories
	cardRepo := infrastructure.NewInMemoryCardRepository()
	accountCacheRepo := infrastructure.NewInMemoryAccountCacheRepository()

	// Setup service
	service := application.NewCardService(cardRepo, accountCacheRepo)

	// Setup presenter
	presenter := presenters.NewResponsePresenter()

	// Setup controllers
	createController := controllers.NewCreateCardController(service.CreateCard, presenter)
	deleteController := controllers.NewDeleteCardController(service.DeleteCard, presenter)
	getController := controllers.NewGetCardController(service.ViewCard, presenter)
	listController := controllers.NewListCardsController(service.ListCards, presenter)

	// Setup router
	ctrls := &routes.Controllers{
		CreateCard: createController,
		DeleteCard: deleteController,
		GetCard:    getController,
		ListCards:  listController,
	}
	router := routes.SetupRoutes(ctrls)

	// Create test server
	server := httptest.NewServer(router)

	return server, cardRepo, accountCacheRepo
}

func TestCreateCardEndpoint(t *testing.T) {
	server, _, accountCacheRepo := setupTestServer()
	defer server.Close()

	// Setup: Add an active account to cache
	account := domain.NewAccountCache("acc-123", "ACTIVE")
	accountCacheRepo.Upsert(account)

	t.Run("Successful card creation", func(t *testing.T) {
		reqBody := map[string]string{
			"country":    "US",
			"account_id": "acc-123",
		}
		body, _ := json.Marshal(reqBody)

		resp, err := http.Post(server.URL+"/card", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		if response["id"] == nil {
			t.Error("Response should contain card ID")
		}
		if response["card_number"] == nil {
			t.Error("Response should contain card number")
		}
		if response["country"] != "US" {
			t.Errorf("Expected country US, got %v", response["country"])
		}
	})

	t.Run("Create card with missing country", func(t *testing.T) {
		reqBody := map[string]string{
			"account_id": "acc-123",
		}
		body, _ := json.Marshal(reqBody)

		resp, err := http.Post(server.URL+"/card", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("Create card for nonexistent account", func(t *testing.T) {
		reqBody := map[string]string{
			"country":    "US",
			"account_id": "acc-999",
		}
		body, _ := json.Marshal(reqBody)

		resp, err := http.Post(server.URL+"/card", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestGetCardEndpoint(t *testing.T) {
	server, cardRepo, accountCacheRepo := setupTestServer()
	defer server.Close()

	// Setup: Create a card
	card, _ := domain.NewCard("card-123", "US-12345", "US", "acc-123", time.Now())
	cardRepo.Create(card)

	account := domain.NewAccountCache("acc-123", "ACTIVE")
	accountCacheRepo.Upsert(account)

	t.Run("Get card by ID", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/card?id=card-123")
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		if response["id"] != "card-123" {
			t.Errorf("Expected ID card-123, got %v", response["id"])
		}
		if response["card_number"] != "US-12345" {
			t.Errorf("Expected card number US-12345, got %v", response["card_number"])
		}
	})

	t.Run("Get nonexistent card", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/card?id=card-999")
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestGetCardByCardNumberEndpoint(t *testing.T) {
	server, cardRepo, accountCacheRepo := setupTestServer()
	defer server.Close()

	// Setup: Create a card
	card, _ := domain.NewCard("card-123", "US-12345", "US", "acc-123", time.Now())
	cardRepo.Create(card)

	account := domain.NewAccountCache("acc-123", "ACTIVE")
	accountCacheRepo.Upsert(account)

	t.Run("Get card by card number", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/cards/by-number?card_number=US-12345")
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		if response["card_number"] != "US-12345" {
			t.Errorf("Expected card number US-12345, got %v", response["card_number"])
		}
	})

	t.Run("Get nonexistent card by number", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/cards/by-number?card_number=US-99999")
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestGetCardsByAccountEndpoint(t *testing.T) {
	server, cardRepo, accountCacheRepo := setupTestServer()
	defer server.Close()

	// Setup: Create multiple cards for same account
	card1, _ := domain.NewCard("card-1", "US-111", "US", "acc-123", time.Now())
	card2, _ := domain.NewCard("card-2", "US-222", "US", "acc-123", time.Now())
	cardRepo.Create(card1)
	cardRepo.Create(card2)

	account := domain.NewAccountCache("acc-123", "ACTIVE")
	accountCacheRepo.Upsert(account)

	t.Run("Get cards by account ID", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/cards/by-account?account_id=acc-123")
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		total := response["total"].(float64)
		if total != 2 {
			t.Errorf("Expected 2 cards, got %v", total)
		}
	})

	t.Run("Get cards for account with no cards", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/cards/by-account?account_id=acc-999")
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		total := response["total"].(float64)
		if total != 0 {
			t.Errorf("Expected 0 cards, got %v", total)
		}
	})
}

func TestListCardsEndpoint(t *testing.T) {
	server, cardRepo, accountCacheRepo := setupTestServer()
	defer server.Close()

	// Setup: Create multiple cards
	card1, _ := domain.NewCard("card-1", "US-111", "US", "acc-123", time.Now())
	card2, _ := domain.NewCard("card-2", "UK-222", "UK", "acc-456", time.Now())
	cardRepo.Create(card1)
	cardRepo.Create(card2)

	account1 := domain.NewAccountCache("acc-123", "ACTIVE")
	account2 := domain.NewAccountCache("acc-456", "ACTIVE")
	accountCacheRepo.Upsert(account1)
	accountCacheRepo.Upsert(account2)

	t.Run("List all cards", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/cards")
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		total := response["total"].(float64)
		if total != 2 {
			t.Errorf("Expected 2 cards, got %v", total)
		}
	})
}

func TestDeleteCardEndpoint(t *testing.T) {
	server, cardRepo, accountCacheRepo := setupTestServer()
	defer server.Close()

	t.Run("Delete card successfully", func(t *testing.T) {
		// Setup: Create a card
		card, _ := domain.NewCard("card-delete-1", "US-DELETE-1", "US", "acc-123", time.Now())
		cardRepo.Create(card)

		account := domain.NewAccountCache("acc-123", "ACTIVE")
		accountCacheRepo.Upsert(account)

		req, _ := http.NewRequest("DELETE", server.URL+"/card?id=card-delete-1", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		// Verify card is deleted
		deletedCard, _ := cardRepo.GetByID("card-delete-1")
		if !deletedCard.Deleted {
			t.Error("Card should be marked as deleted")
		}
	})

	t.Run("Delete nonexistent card", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", server.URL+"/card?id=card-999", nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

func TestHealthCheckEndpoint(t *testing.T) {
	server, _, _ := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)

	if response["status"] != "healthy" {
		t.Errorf("Expected status healthy, got %s", response["status"])
	}
}
