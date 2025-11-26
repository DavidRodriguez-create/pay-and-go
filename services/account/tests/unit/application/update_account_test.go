package application_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DavidRodriguez-create/pay-and-go/services/account/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/account/domain"
)

func TestUpdateAccount(t *testing.T) {
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
		request    application.UpdateAccountRequest
		setupMock  func(*MockAccountRepository)
		wantErr    bool
		errMessage string
	}{
		{
			name: "Successful account update",
			request: application.UpdateAccountRequest{
				ID:           "123",
				BeholderName: "Jane Doe",
				CountryCode:  "UK",
			},
			setupMock: func(m *MockAccountRepository) {
				m.GetByIDFunc = func(id string) (*domain.Account, error) {
					return existingAccount, nil
				}
				m.UpdateFunc = func(account *domain.Account) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "Update all fields",
			request: application.UpdateAccountRequest{
				ID:            "123",
				AccountNumber: "ACC002",
				BeholderName:  "Jane Smith",
				CountryCode:   "CA",
				Status:        string(domain.StatusBlocked),
			},
			setupMock: func(m *MockAccountRepository) {
				m.GetByIDFunc = func(id string) (*domain.Account, error) {
					return existingAccount, nil
				}
				m.UpdateFunc = func(account *domain.Account) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "Account not found",
			request: application.UpdateAccountRequest{
				ID:           "999",
				BeholderName: "Jane Doe",
			},
			setupMock: func(m *MockAccountRepository) {
				m.GetByIDFunc = func(id string) (*domain.Account, error) {
					return nil, errors.New("account not found")
				}
			},
			wantErr:    true,
			errMessage: "account not found",
		},
		{
			name: "Update deleted account",
			request: application.UpdateAccountRequest{
				ID:           "123",
				BeholderName: "Jane Doe",
			},
			setupMock: func(m *MockAccountRepository) {
				deletedAccount := &domain.Account{
					ID:            "123",
					AccountNumber: "ACC001",
					BeholderName:  "John Doe",
					CountryCode:   "US",
					Status:        domain.StatusDeleted,
					CreatedAt:     time.Now().Add(-24 * time.Hour),
					UpdatedAt:     time.Now().Add(-24 * time.Hour),
				}
				m.GetByIDFunc = func(id string) (*domain.Account, error) {
					return deletedAccount, nil
				}
			},
			wantErr:    true,
			errMessage: "cannot update deleted account",
		},
		{
			name: "Repository update error",
			request: application.UpdateAccountRequest{
				ID:           "123",
				BeholderName: "Jane Doe",
			},
			setupMock: func(m *MockAccountRepository) {
				m.GetByIDFunc = func(id string) (*domain.Account, error) {
					return existingAccount, nil
				}
				m.UpdateFunc = func(account *domain.Account) error {
					return errors.New("database error")
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
			service := application.NewAccountService(mockRepo, &MockEventPublisher{})

			err := service.UpdateAccount(tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.errMessage != "" && err != nil && err.Error() != tt.errMessage {
				t.Errorf("UpdateAccount() error message = %v, want %v", err.Error(), tt.errMessage)
			}
		})
	}
}
