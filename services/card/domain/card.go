package domain

import (
	"errors"
	"time"
)

// Card represents a payment card entity
type Card struct {
	ID                string
	CardNumber        string
	Country           string
	AccountID         string
	Deleted           bool
	CreationTimestamp time.Time
}

// Card validation errors
var (
	ErrCardIDRequired     = errors.New("card ID is required")
	ErrCardNumberRequired = errors.New("card number is required")
	ErrCountryRequired    = errors.New("country is required")
	ErrAccountIDRequired  = errors.New("account ID is required")
	ErrAccountNotFound    = errors.New("account not found")
	ErrAccountDeleted     = errors.New("cannot create card for deleted account")
	ErrAccountInactive    = errors.New("cannot create card for inactive account")
	ErrCardNotFound       = errors.New("card not found")
	ErrCardAlreadyDeleted = errors.New("card is already deleted")
)

// NewCard creates a new Card with validation
func NewCard(id, cardNumber, country, accountID string, creationTimestamp time.Time) (*Card, error) {
	if id == "" {
		return nil, ErrCardIDRequired
	}
	if cardNumber == "" {
		return nil, ErrCardNumberRequired
	}
	if country == "" {
		return nil, ErrCountryRequired
	}
	if accountID == "" {
		return nil, ErrAccountIDRequired
	}

	return &Card{
		ID:                id,
		CardNumber:        cardNumber,
		Country:           country,
		AccountID:         accountID,
		Deleted:           false,
		CreationTimestamp: creationTimestamp,
	}, nil
}

// IsDeleted checks if the card is marked as deleted
func (c *Card) IsDeleted() bool {
	return c.Deleted
}

// Delete marks the card as deleted (soft delete)
func (c *Card) Delete() error {
	if c.Deleted {
		return ErrCardAlreadyDeleted
	}
	c.Deleted = true
	return nil
}
