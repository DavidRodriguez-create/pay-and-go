package domain_test

import (
	"testing"
	"time"

	"github.com/DavidRodriguez-create/pay-and-go/services/account/domain"
)

func TestNewAccount(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		accountNumber string
		beholderName  string
		countryCode   string
		wantErr       bool
	}{
		{
			name:          "Valid account creation",
			id:            "123",
			accountNumber: "ACC001",
			beholderName:  "John Doe",
			countryCode:   "US",
			wantErr:       false,
		},
		{
			name:          "Missing ID",
			id:            "",
			accountNumber: "ACC001",
			beholderName:  "John Doe",
			countryCode:   "US",
			wantErr:       true,
		},
		{
			name:          "Missing account number",
			id:            "123",
			accountNumber: "",
			beholderName:  "John Doe",
			countryCode:   "US",
			wantErr:       true,
		},
		{
			name:          "Missing beholder name",
			id:            "123",
			accountNumber: "ACC001",
			beholderName:  "",
			countryCode:   "US",
			wantErr:       true,
		},
		{
			name:          "Missing country code",
			id:            "123",
			accountNumber: "ACC001",
			beholderName:  "John Doe",
			countryCode:   "",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := time.Now()
			account, err := domain.NewAccount(tt.id, tt.accountNumber, tt.beholderName, tt.countryCode)
			after := time.Now()

			if (err != nil) != tt.wantErr {
				t.Errorf("NewAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if account == nil {
					t.Error("NewAccount() returned nil account")
					return
				}
				if account.ID != tt.id {
					t.Errorf("NewAccount() ID = %v, want %v", account.ID, tt.id)
				}
				if account.AccountNumber != tt.accountNumber {
					t.Errorf("NewAccount() AccountNumber = %v, want %v", account.AccountNumber, tt.accountNumber)
				}
				if account.BeholderName != tt.beholderName {
					t.Errorf("NewAccount() BeholderName = %v, want %v", account.BeholderName, tt.beholderName)
				}
				if account.CountryCode != tt.countryCode {
					t.Errorf("NewAccount() CountryCode = %v, want %v", account.CountryCode, tt.countryCode)
				}
				if account.Status != domain.StatusActive {
					t.Errorf("NewAccount() Status = %v, want %v", account.Status, domain.StatusActive)
				}
				if account.CreatedAt.Before(before) || account.CreatedAt.After(after) {
					t.Errorf("NewAccount() CreatedAt = %v, should be between %v and %v", account.CreatedAt, before, after)
				}
				if account.UpdatedAt.Before(before) || account.UpdatedAt.After(after) {
					t.Errorf("NewAccount() UpdatedAt = %v, should be between %v and %v", account.UpdatedAt, before, after)
				}
			}
		})
	}
}

func TestAccountStatusMethods(t *testing.T) {
	account := &domain.Account{
		ID:            "123",
		AccountNumber: "ACC001",
		BeholderName:  "John Doe",
		CountryCode:   "US",
		Status:        domain.StatusActive,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	t.Run("IsActive", func(t *testing.T) {
		if !account.IsActive() {
			t.Error("Account should be active")
		}
		account.Status = domain.StatusDeleted
		if account.IsActive() {
			t.Error("Account should not be active")
		}
	})

	t.Run("IsDeleted", func(t *testing.T) {
		account.Status = domain.StatusDeleted
		if !account.IsDeleted() {
			t.Error("Account should be deleted")
		}
		account.Status = domain.StatusActive
		if account.IsDeleted() {
			t.Error("Account should not be deleted")
		}
	})

	t.Run("IsBlocked", func(t *testing.T) {
		account.Status = domain.StatusBlocked
		if !account.IsBlocked() {
			t.Error("Account should be blocked")
		}
		account.Status = domain.StatusActive
		if account.IsBlocked() {
			t.Error("Account should not be blocked")
		}
	})
}
