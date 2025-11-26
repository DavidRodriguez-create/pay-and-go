package domain

import "errors"

// AccountStatus represents the status of an account
type AccountStatus string

const (
	AccountStatusActive  AccountStatus = "ACTIVE"
	AccountStatusBlocked AccountStatus = "BLOCKED"
	AccountStatusDeleted AccountStatus = "DELETED"
)

// AccountCache represents a cached account for validation purposes
type AccountCache struct {
	ID     string
	Status AccountStatus
}

// Account cache validation errors
var (
	ErrAccountCacheNotFound = errors.New("account cache not found")
)

// NewAccountCache creates a new AccountCache
func NewAccountCache(id string, status AccountStatus) *AccountCache {
	return &AccountCache{
		ID:     id,
		Status: status,
	}
}

// IsActive checks if the account is active
func (a *AccountCache) IsActive() bool {
	return a.Status == AccountStatusActive
}

// IsDeleted checks if the account is deleted
func (a *AccountCache) IsDeleted() bool {
	return a.Status == AccountStatusDeleted
}

// IsBlocked checks if the account is blocked
func (a *AccountCache) IsBlocked() bool {
	return a.Status == AccountStatusBlocked
}
