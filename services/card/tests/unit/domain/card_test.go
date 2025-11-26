package domain_test

import (
	"testing"
	"time"

	"github.com/DavidRodriguez-create/pay-and-go/services/card/domain"
)

func TestNewCard(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		id          string
		cardNumber  string
		country     string
		accountID   string
		timestamp   time.Time
		expectError error
	}{
		{
			name:        "Valid card creation",
			id:          "card-123",
			cardNumber:  "US-12345678",
			country:     "US",
			accountID:   "acc-456",
			timestamp:   now,
			expectError: nil,
		},
		{
			name:        "Missing ID",
			id:          "",
			cardNumber:  "US-12345678",
			country:     "US",
			accountID:   "acc-456",
			timestamp:   now,
			expectError: domain.ErrCardIDRequired,
		},
		{
			name:        "Missing card number",
			id:          "card-123",
			cardNumber:  "",
			country:     "US",
			accountID:   "acc-456",
			timestamp:   now,
			expectError: domain.ErrCardNumberRequired,
		},
		{
			name:        "Missing country",
			id:          "card-123",
			cardNumber:  "US-12345678",
			country:     "",
			accountID:   "acc-456",
			timestamp:   now,
			expectError: domain.ErrCountryRequired,
		},
		{
			name:        "Missing account ID",
			id:          "card-123",
			cardNumber:  "US-12345678",
			country:     "US",
			accountID:   "",
			timestamp:   now,
			expectError: domain.ErrAccountIDRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card, err := domain.NewCard(tt.id, tt.cardNumber, tt.country, tt.accountID, tt.timestamp)

			if tt.expectError != nil {
				if err != tt.expectError {
					t.Errorf("Expected error %v, got %v", tt.expectError, err)
				}
				if card != nil {
					t.Error("Expected nil card when error occurs")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if card.ID != tt.id {
				t.Errorf("Expected ID %s, got %s", tt.id, card.ID)
			}
			if card.CardNumber != tt.cardNumber {
				t.Errorf("Expected CardNumber %s, got %s", tt.cardNumber, card.CardNumber)
			}
			if card.Country != tt.country {
				t.Errorf("Expected Country %s, got %s", tt.country, card.Country)
			}
			if card.AccountID != tt.accountID {
				t.Errorf("Expected AccountID %s, got %s", tt.accountID, card.AccountID)
			}
			if card.Deleted {
				t.Error("New card should not be deleted")
			}
			if !card.CreationTimestamp.Equal(tt.timestamp) {
				t.Errorf("Expected timestamp %v, got %v", tt.timestamp, card.CreationTimestamp)
			}
		})
	}
}

func TestCardMethods(t *testing.T) {
	t.Run("IsDeleted", func(t *testing.T) {
		card, _ := domain.NewCard("card-1", "US-123", "US", "acc-1", time.Now())

		if card.IsDeleted() {
			t.Error("New card should not be deleted")
		}

		card.Deleted = true
		if !card.IsDeleted() {
			t.Error("Card marked as deleted should return true")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		card, _ := domain.NewCard("card-1", "US-123", "US", "acc-1", time.Now())

		err := card.Delete()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !card.Deleted {
			t.Error("Card should be marked as deleted")
		}
	})

	t.Run("Delete already deleted card", func(t *testing.T) {
		card, _ := domain.NewCard("card-1", "US-123", "US", "acc-1", time.Now())
		card.Delete()

		err := card.Delete()
		if err != domain.ErrCardAlreadyDeleted {
			t.Errorf("Expected error %v, got %v", domain.ErrCardAlreadyDeleted, err)
		}
	})
}
