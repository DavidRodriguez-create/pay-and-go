package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DavidRodriguez-create/pay-and-go/services/account/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/account/infrastructure"
	"github.com/DavidRodriguez-create/pay-and-go/services/account/presentation/controllers"
	"github.com/DavidRodriguez-create/pay-and-go/services/account/presentation/routes"
)

// setupTestServer creates a test HTTP server with all dependencies
func setupTestServer() *http.ServeMux {
	repo := infrastructure.NewInMemoryAccountRepository()
	// Use nil event publisher for tests (events not needed in test environment)
	service := application.NewAccountService(repo, nil)

	ctrls := &routes.Controllers{
		CreateAccount: controllers.NewCreateAccountController(service),
		GetAccount:    controllers.NewGetAccountController(service),
		ListAccounts:  controllers.NewListAccountsController(service),
		UpdateAccount: controllers.NewUpdateAccountController(service),
		DeleteAccount: controllers.NewDeleteAccountController(service),
	}

	return routes.SetupRoutes(ctrls)
}

func TestAccountAPIIntegration(t *testing.T) {
	mux := setupTestServer()

	t.Run("Complete account lifecycle", func(t *testing.T) {
		// 1. Create an account
		createReq := map[string]interface{}{
			"beholder_name": "Integration Test User",
			"country_code":  "US",
		}
		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/account", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected status 201, got %d: %s", w.Code, w.Body.String())
		}

		var createResp map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&createResp); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		accountID, ok := createResp["id"].(string)
		if !ok || accountID == "" {
			t.Fatal("Expected account ID in response")
		}

		accountNumber, ok := createResp["account_number"].(string)
		if !ok || accountNumber == "" {
			t.Fatal("Expected account number in response")
		}

		// 2. Get account by ID
		req = httptest.NewRequest(http.MethodGet, "/account?id="+accountID, nil)
		w = httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", w.Code)
		}

		var getResp map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&getResp); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if getResp["beholder_name"] != "Integration Test User" {
			t.Errorf("Expected beholder_name 'Integration Test User', got %v", getResp["beholder_name"])
		}

		// 3. Get account by account number
		req = httptest.NewRequest(http.MethodGet, "/accounts/by-number?account_number="+accountNumber, nil)
		w = httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", w.Code)
		}

		// 4. Update account
		updateReq := map[string]interface{}{
			"id":            accountID,
			"beholder_name": "Updated User",
			"country_code":  "UK",
			"status":        "BLOCKED",
		}
		body2, _ := json.Marshal(updateReq)
		req = httptest.NewRequest(http.MethodPut, "/account?id="+accountID, bytes.NewReader(body2))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		// 5. Verify update
		req = httptest.NewRequest(http.MethodGet, "/account?id="+accountID, nil)
		w = httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if err := json.NewDecoder(w.Body).Decode(&getResp); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if getResp["beholder_name"] != "Updated User" {
			t.Errorf("Expected beholder_name 'Updated User', got %v", getResp["beholder_name"])
		}
		if getResp["country_code"] != "UK" {
			t.Errorf("Expected country_code 'UK', got %v", getResp["country_code"])
		}
		if getResp["status"] != "BLOCKED" {
			t.Errorf("Expected status 'BLOCKED', got %v", getResp["status"])
		}

		// 6. List accounts (should have 1)
		req = httptest.NewRequest(http.MethodGet, "/accounts", nil)
		w = httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", w.Code)
		}

		var listResp map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&listResp); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		total, ok := listResp["total"].(float64)
		if !ok || total != 1 {
			t.Errorf("Expected total 1, got %v", listResp["total"])
		}

		// 6. Delete account
		req = httptest.NewRequest(http.MethodDelete, "/account?id="+accountID, nil)
		w = httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", w.Code)
		}

		// 8. Verify account is marked as deleted
		req = httptest.NewRequest(http.MethodGet, "/account?id="+accountID, nil)
		w = httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if err := json.NewDecoder(w.Body).Decode(&getResp); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if getResp["status"] != "DELETED" {
			t.Errorf("Expected status 'DELETED', got %v", getResp["status"])
		}
	})

	t.Run("Multiple accounts", func(t *testing.T) {
		mux := setupTestServer()

		// Create multiple accounts
		for i := 0; i < 3; i++ {
			createReq := map[string]interface{}{
				"beholder_name": "User" + string(rune('A'+i)),
				"country_code":  "US",
			}
			body, _ := json.Marshal(createReq)
			req := httptest.NewRequest(http.MethodPost, "/account", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != http.StatusCreated {
				t.Fatalf("Failed to create account %d: %d", i, w.Code)
			}
		}

		// List all accounts
		req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", w.Code)
		}

		var listResp map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&listResp); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		total, ok := listResp["total"].(float64)
		if !ok || total != 3 {
			t.Errorf("Expected total 3, got %v", listResp["total"])
		}
	})
}

func TestHealthEndpoint(t *testing.T) {
	mux := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", resp["status"])
	}
	if resp["service"] != "account-service" {
		t.Errorf("Expected service 'account-service', got %v", resp["service"])
	}
}

func TestErrorHandling(t *testing.T) {
	mux := setupTestServer()

	t.Run("Invalid JSON payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/account", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("Missing required fields", func(t *testing.T) {
		createReq := map[string]interface{}{
			"beholder_name": "",
			"country_code":  "US",
		}
		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/account", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("Get non-existent account", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/account?id=nonexistent", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Expected status 404, got %d", w.Code)
		}
	})

	t.Run("Update non-existent account", func(t *testing.T) {
		updateReq := map[string]interface{}{
			"id":            "nonexistent",
			"beholder_name": "Updated User",
		}
		body, _ := json.Marshal(updateReq)
		req := httptest.NewRequest(http.MethodPut, "/account?id=nonexistent", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("Delete non-existent account", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/account?id=nonexistent", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("Method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/health", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Fatalf("Expected status 405, got %d", w.Code)
		}
	})
}
