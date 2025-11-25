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
	return s.repository.Delete(id)
}
