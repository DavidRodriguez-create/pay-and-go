package application

import (
	"github.com/DavidRodriguez-create/pay-and-go/services/account/domain"
)

// AccountService defines the interface for account business operations
type AccountService interface {
	CreateAccount(req CreateAccountRequest) (*AccountResponse, error)
	GetAccountByID(id string) (*AccountResponse, error)
	GetAccountByAccountNumber(accountNumber string) (*AccountResponse, error)
	ListAccounts() (*AccountListResponse, error)
	UpdateAccount(req UpdateAccountRequest) error
	DeleteAccount(id string) error
}

// Ensure use cases implement the service interface
var (
	_ AccountService = (*AccountServiceImpl)(nil)
)

// AccountServiceImpl implements the AccountService interface
type AccountServiceImpl struct {
	repository     domain.AccountRepository
	eventPublisher domain.EventPublisher
}

// NewAccountService creates a new instance of AccountServiceImpl
func NewAccountService(repository domain.AccountRepository, eventPublisher domain.EventPublisher) *AccountServiceImpl {
	return &AccountServiceImpl{
		repository:     repository,
		eventPublisher: eventPublisher,
	}
}
