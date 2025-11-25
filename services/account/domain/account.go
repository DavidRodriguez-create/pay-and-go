package domain

import (
	"errors"
	"time"
)

type AccountStatus string

const (
	StatusActive  AccountStatus = "ACTIVE"
	StatusDeleted AccountStatus = "DELETED"
	StatusBlocked AccountStatus = "BLOCKED"
)

type Account struct {
	ID            string
	AccountNumber string
	BeholderName  string
	CountryCode   string
	Status        AccountStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewAccount(id, accountNumber, beholderName, countryCode string) (*Account, error) {
	if id == "" || accountNumber == "" || beholderName == "" || countryCode == "" {
		return nil, errors.New("all fields are required to create an account")
	}
	return &Account{
		ID:            id,
		AccountNumber: accountNumber,
		BeholderName:  beholderName,
		CountryCode:   countryCode,
		Status:        StatusActive,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

func (a *Account) IsActive() bool {
	return a.Status == StatusActive
}

func (a *Account) IsDeleted() bool {
	return a.Status == StatusDeleted
}

func (a *Account) IsBlocked() bool {
	return a.Status == StatusBlocked
}
