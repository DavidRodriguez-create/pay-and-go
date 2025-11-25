package application_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DavidRodriguez-create/pay-and-go/services/account/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/account/domain"
)

func TestDeleteAccount(t *testing.T) {
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
			name:      "Successful account deletion",
			accountID: "123",
			setupMock: func(m *MockAccountRepository) {
				m.GetByIDFunc = func(id string) (*domain.Account, error) {
					return existingAccount, nil
				}
				m.DeleteFunc = func(id string) error {
					return nil
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
		{
			name:      "Delete already deleted account",
			accountID: "123",
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
			errMessage: "account is already deleted",
		},
		{
			name:      "Repository delete error",
			accountID: "123",
			setupMock: func(m *MockAccountRepository) {
				m.GetByIDFunc = func(id string) (*domain.Account, error) {
					return existingAccount, nil
				}
				m.DeleteFunc = func(id string) error {
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
			service := application.NewAccountService(mockRepo)

			err := service.DeleteAccount(tt.accountID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.errMessage != "" && err != nil && err.Error() != tt.errMessage {
				t.Errorf("DeleteAccount() error message = %v, want %v", err.Error(), tt.errMessage)
			}
		})
	}
}
