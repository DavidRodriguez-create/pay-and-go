package application_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DavidRodriguez-create/pay-and-go/services/account/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/account/domain"
)

func TestGetAccountByID(t *testing.T) {
	existingAccount := &domain.Account{
		ID:            "123",
		AccountNumber: "ACC001",
		BeholderName:  "John Doe",
		CountryCode:   "US",
		Status:        domain.StatusActive,
		CreatedAt:     time.Now().Add(-24 * time.Hour),
		UpdatedAt:     time.Now().Add(-24 * time.Hour),
	}

	tests := []struct {
		name       string
		accountID  string
		setupMock  func(*MockAccountRepository)
		wantErr    bool
		errMessage string
	}{
		{
			name:      "Successful retrieval",
			accountID: "123",
			setupMock: func(m *MockAccountRepository) {
				m.GetByIDFunc = func(id string) (*domain.Account, error) {
					return existingAccount, nil
				}
			},
			wantErr: false,
		},
		{
			name:      "Account not found",
			accountID: "999",
			setupMock: func(m *MockAccountRepository) {
				m.GetByIDFunc = func(id string) (*domain.Account, error) {
					return nil, errors.New("account not found")
				}
			},
			wantErr:    true,
			errMessage: "account not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockAccountRepository{}
			tt.setupMock(mockRepo)
			service := application.NewAccountService(mockRepo)

			response, err := service.GetAccountByID(tt.accountID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.errMessage != "" && err != nil && err.Error() != tt.errMessage {
				t.Errorf("GetAccountByID() error message = %v, want %v", err.Error(), tt.errMessage)
			}

			if !tt.wantErr {
				if response == nil {
					t.Error("GetAccountByID() returned nil response")
					return
				}
				if response.ID != existingAccount.ID {
					t.Errorf("GetAccountByID() ID = %v, want %v", response.ID, existingAccount.ID)
				}
			}
		})
	}
}

func TestGetAccountByAccountNumber(t *testing.T) {
	existingAccount := &domain.Account{
		ID:            "123",
		AccountNumber: "ACC001",
		BeholderName:  "John Doe",
		CountryCode:   "US",
		Status:        domain.StatusActive,
		CreatedAt:     time.Now().Add(-24 * time.Hour),
		UpdatedAt:     time.Now().Add(-24 * time.Hour),
	}

	tests := []struct {
		name          string
		accountNumber string
		setupMock     func(*MockAccountRepository)
		wantErr       bool
		errMessage    string
	}{
		{
			name:          "Successful retrieval",
			accountNumber: "ACC001",
			setupMock: func(m *MockAccountRepository) {
				m.GetByAccountNumberFunc = func(accountNumber string) (*domain.Account, error) {
					return existingAccount, nil
				}
			},
			wantErr: false,
		},
		{
			name:          "Account not found",
			accountNumber: "ACC999",
			setupMock: func(m *MockAccountRepository) {
				m.GetByAccountNumberFunc = func(accountNumber string) (*domain.Account, error) {
					return nil, errors.New("account not found")
				}
			},
			wantErr:    true,
			errMessage: "account not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockAccountRepository{}
			tt.setupMock(mockRepo)
			service := application.NewAccountService(mockRepo)

			response, err := service.GetAccountByAccountNumber(tt.accountNumber)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountByAccountNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.errMessage != "" && err != nil && err.Error() != tt.errMessage {
				t.Errorf("GetAccountByAccountNumber() error message = %v, want %v", err.Error(), tt.errMessage)
			}

			if !tt.wantErr {
				if response == nil {
					t.Error("GetAccountByAccountNumber() returned nil response")
					return
				}
				if response.AccountNumber != existingAccount.AccountNumber {
					t.Errorf("GetAccountByAccountNumber() AccountNumber = %v, want %v", response.AccountNumber, existingAccount.AccountNumber)
				}
			}
		})
	}
}

func TestListAccounts(t *testing.T) {
	accounts := []*domain.Account{
		{
			ID:            "123",
			AccountNumber: "ACC001",
			BeholderName:  "John Doe",
			CountryCode:   "US",
			Status:        domain.StatusActive,
			CreatedAt:     time.Now().Add(-24 * time.Hour),
			UpdatedAt:     time.Now().Add(-24 * time.Hour),
		},
		{
			ID:            "456",
			AccountNumber: "ACC002",
			BeholderName:  "Jane Smith",
			CountryCode:   "UK",
			Status:        domain.StatusActive,
			CreatedAt:     time.Now().Add(-12 * time.Hour),
			UpdatedAt:     time.Now().Add(-12 * time.Hour),
		},
	}

	tests := []struct {
		name       string
		setupMock  func(*MockAccountRepository)
		wantCount  int
		wantErr    bool
		errMessage string
	}{
		{
			name: "Successful list with accounts",
			setupMock: func(m *MockAccountRepository) {
				m.ListFunc = func() ([]*domain.Account, error) {
					return accounts, nil
				}
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "Successful list with no accounts",
			setupMock: func(m *MockAccountRepository) {
				m.ListFunc = func() ([]*domain.Account, error) {
					return []*domain.Account{}, nil
				}
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "Repository error",
			setupMock: func(m *MockAccountRepository) {
				m.ListFunc = func() ([]*domain.Account, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr:    true,
			errMessage: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockAccountRepository{}
			tt.setupMock(mockRepo)
			service := application.NewAccountService(mockRepo)

			response, err := service.ListAccounts()

			if (err != nil) != tt.wantErr {
				t.Errorf("ListAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.errMessage != "" && err != nil && err.Error() != tt.errMessage {
				t.Errorf("ListAccounts() error message = %v, want %v", err.Error(), tt.errMessage)
			}

			if !tt.wantErr {
				if response == nil {
					t.Error("ListAccounts() returned nil response")
					return
				}
				if response.Total != tt.wantCount {
					t.Errorf("ListAccounts() Total = %v, want %v", response.Total, tt.wantCount)
				}
				if len(response.Accounts) != tt.wantCount {
					t.Errorf("ListAccounts() Accounts length = %v, want %v", len(response.Accounts), tt.wantCount)
				}
			}
		})
	}
}
