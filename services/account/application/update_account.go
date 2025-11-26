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

	// Track if status changed
	statusChanged := false

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
		newStatus := ToAccountStatus(req.Status)
		if existingAccount.Status != newStatus {
			existingAccount.Status = newStatus
			statusChanged = true
		}
	}
	existingAccount.UpdatedAt = time.Now()

	// Persist changes
	if err := s.repository.Update(existingAccount); err != nil {
		return err
	}

	// Publish account.status_changed event if status was modified
	if statusChanged && s.eventPublisher != nil {
		if err := s.eventPublisher.PublishAccountStatusChanged(existingAccount.ID, string(existingAccount.Status)); err != nil {
			// Log error but don't fail the request - event publishing is best-effort
			// The update was successful, event delivery failure shouldn't rollback the operation
		}
	}

	return nil
}
