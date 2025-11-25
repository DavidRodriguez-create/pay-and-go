package application

import (
	"errors"
	"time"
)

// UpdateAccount updates an existing account
func (s *AccountServiceImpl) UpdateAccount(req UpdateAccountRequest) error {
	// Verify account exists
	existingAccount, err := s.repository.GetByID(req.ID)
	if err != nil {
		return err
	}

	if existingAccount.IsDeleted() {
		return errors.New("cannot update deleted account")
	}

	if req.AccountNumber != "" {
		existingAccount.AccountNumber = req.AccountNumber
	}
	if req.BeholderName != "" {
		existingAccount.BeholderName = req.BeholderName
	}
	if req.CountryCode != "" {
		existingAccount.CountryCode = req.CountryCode
	}
	if req.Status != "" {
		existingAccount.Status = ToAccountStatus(req.Status)
	}
	existingAccount.UpdatedAt = time.Now()

	// Persist changes
	return s.repository.Update(existingAccount)
}
