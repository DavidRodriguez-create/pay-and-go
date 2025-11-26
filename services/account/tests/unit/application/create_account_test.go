package application_test

import (
	"errors"
	"testing"

	"github.com/DavidRodriguez-create/pay-and-go/services/account/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/account/domain"
)

// MockEventPublisher for testing
type MockEventPublisher struct {
	PublishAccountCreatedFunc       func(accountID string, status string) error
	PublishAccountStatusChangedFunc func(accountID string, status string) error
}

func (m *MockEventPublisher) PublishAccountCreated(accountID string, status string) error {
	if m.PublishAccountCreatedFunc != nil {
		return m.PublishAccountCreatedFunc(accountID, status)
	}
	return nil
}

func (m *MockEventPublisher) PublishAccountStatusChanged(accountID string, status string) error {
	if m.PublishAccountStatusChangedFunc != nil {
		return m.PublishAccountStatusChangedFunc(accountID, status)
	}
	return nil
}

// MockAccountRepository for testing
type MockAccountRepository struct {
	CreateFunc             func(account *domain.Account) error
	GetByIDFunc            func(id string) (*domain.Account, error)
	GetByAccountNumberFunc func(accountNumber string) (*domain.Account, error)
	UpdateFunc             func(account *domain.Account) error
	DeleteFunc             func(id string) error
	ListFunc               func() ([]*domain.Account, error)
}

func (m *MockAccountRepository) Create(account *domain.Account) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(account)
	}
	return nil
}

func (m *MockAccountRepository) GetByID(id string) (*domain.Account, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, errors.New("not found")
}

func (m *MockAccountRepository) GetByAccountNumber(accountNumber string) (*domain.Account, error) {
	if m.GetByAccountNumberFunc != nil {
		return m.GetByAccountNumberFunc(accountNumber)
	}
	return nil, errors.New("not found")
}

func (m *MockAccountRepository) Update(account *domain.Account) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(account)
	}
	return nil
}

func (m *MockAccountRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *MockAccountRepository) List() ([]*domain.Account, error) {
	if m.ListFunc != nil {
		return m.ListFunc()
	}
	return []*domain.Account{}, nil
}

func TestCreateAccount(t *testing.T) {
	tests := []struct {
		name       string
		request    application.CreateAccountRequest
		setupMock  func(*MockAccountRepository)
		wantErr    bool
		errMessage string
	}{
		{
			name: "Successful account creation",
			request: application.CreateAccountRequest{
				BeholderName: "John Doe",
				CountryCode:  "US",
			},
			setupMock: func(m *MockAccountRepository) {
				m.CreateFunc = func(account *domain.Account) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "Missing beholder name",
			request: application.CreateAccountRequest{
				BeholderName: "",
				CountryCode:  "US",
			},
			setupMock: func(m *MockAccountRepository) {},
			wantErr:   true,
		},
		{
			name: "Missing country code",
			request: application.CreateAccountRequest{
				BeholderName: "John Doe",
				CountryCode:  "",
			},
			setupMock: func(m *MockAccountRepository) {},
			wantErr:   true,
		},
		{
			name: "Repository creation error",
			request: application.CreateAccountRequest{
				BeholderName: "John Doe",
				CountryCode:  "US",
			},
			setupMock: func(m *MockAccountRepository) {
				m.CreateFunc = func(account *domain.Account) error {
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

			response, err := service.CreateAccount(tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.errMessage != "" && err != nil && err.Error() != tt.errMessage {
				t.Errorf("CreateAccount() error message = %v, want %v", err.Error(), tt.errMessage)
			}

			if !tt.wantErr {
				if response == nil {
					t.Error("CreateAccount() returned nil response")
					return
				}
				if response.BeholderName != tt.request.BeholderName {
					t.Errorf("CreateAccount() BeholderName = %v, want %v", response.BeholderName, tt.request.BeholderName)
				}
				if response.CountryCode != tt.request.CountryCode {
					t.Errorf("CreateAccount() CountryCode = %v, want %v", response.CountryCode, tt.request.CountryCode)
				}
				if response.Status != string(domain.StatusActive) {
					t.Errorf("CreateAccount() Status = %v, want %v", response.Status, domain.StatusActive)
				}
				if response.ID == "" {
					t.Error("CreateAccount() ID should not be empty")
				}
				if response.AccountNumber == "" {
					t.Error("CreateAccount() AccountNumber should not be empty")
				}
			}
		})
	}
}
