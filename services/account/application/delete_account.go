package application

import (
	"errors"
)

// DeleteAccount deletes an account (soft delete)
func (s *AccountServiceImpl) DeleteAccount(id string) error {
	// Verify account exists
	account, err := s.repository.GetByID(id)
	if err != nil {
		return err
	}

	if account.IsDeleted() {
		return errors.New("account is already deleted")
	}

	// Perform soft delete
	if err := s.repository.Delete(id); err != nil {
		return err
	}

	// Publish account deletion event
	if s.eventPublisher != nil {
		if err := s.eventPublisher.PublishAccountStatusChanged(id, "DELETED"); err != nil {
			// Log error but don't fail the operation
			return err
		}
	}

	return nil
}
