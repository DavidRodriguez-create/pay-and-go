package application

import (
	"github.com/DavidRodriguez-create/pay-and-go/services/account/domain"
	"github.com/google/uuid"
)

// CreateAccount creates a new account
func (s *AccountServiceImpl) CreateAccount(req CreateAccountRequest) (*AccountResponse, error) {
	// TODO: Improve UUID generation - consider using a more robust ID generation service
	// or allow custom ID/AccountNumber providers for better testability and flexibility
	id := uuid.New().String()
	accountNumber := uuid.New().String()

	// Create the account entity with validation
	account, err := domain.NewAccount(id, accountNumber, req.BeholderName, req.CountryCode)
	if err != nil {
		return nil, err
	}

	// Persist the account
	err = s.repository.Create(account)
	if err != nil {
		return nil, err
	}

	// Publish account.created event
	if s.eventPublisher != nil {
		if err := s.eventPublisher.PublishAccountCreated(account.ID, string(account.Status)); err != nil {
			// Log error but don't fail the request - event publishing is best-effort
			// In production, you might want to implement retry logic or dead-letter queue
			// For now, we continue as the account was successfully created
		}
	}

	return ToAccountResponse(account), nil
}
