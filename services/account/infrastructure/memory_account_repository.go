package infrastructure

import (
	"errors"
	"sync"
	"time"

	"github.com/DavidRodriguez-create/pay-and-go/services/account/domain"
)

// InMemoryAccountRepository implements the AccountRepository interface using in-memory storage
// TODO: Replace with database connection in the future
type InMemoryAccountRepository struct {
	accounts map[string]*domain.Account
	mu       sync.RWMutex
}

// NewInMemoryAccountRepository creates a new instance of InMemoryAccountRepository
func NewInMemoryAccountRepository() *InMemoryAccountRepository {
	return &InMemoryAccountRepository{
		accounts: make(map[string]*domain.Account),
	}
}

// ------- Implementing AccountRepository interface -------

// Create adds a new account to the repository
func (r *InMemoryAccountRepository) Create(account *domain.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.accounts[account.ID]; exists {
		return errors.New("account with this ID already exists")
	}

	if _, err := r.findByAccountNumber(account.AccountNumber); err == nil {
		return errors.New("account with this account number already exists")
	}

	r.accounts[account.ID] = account
	return nil
}

// GetByID retrieves an account by its ID
func (r *InMemoryAccountRepository) GetByID(id string) (*domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	account, exists := r.accounts[id]
	if !exists {
		return nil, errors.New("account not found")
	}

	return account, nil
}

// GetByAccountNumber retrieves an account by its account number
func (r *InMemoryAccountRepository) GetByAccountNumber(accountNumber string) (*domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.findByAccountNumber(accountNumber)
}

// findByAccountNumber is an internal helper method (no lock needed, caller must lock)
func (r *InMemoryAccountRepository) findByAccountNumber(accountNumber string) (*domain.Account, error) {
	for _, account := range r.accounts {
		if account.AccountNumber == accountNumber {
			return account, nil
		}
	}
	return nil, errors.New("account not found")
}

// Update updates an existing account
func (r *InMemoryAccountRepository) Update(account *domain.Account) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.accounts[account.ID]; !exists {
		return errors.New("account not found")
	}

	account.UpdatedAt = time.Now()
	r.accounts[account.ID] = account
	return nil
}

// Delete removes an account from the repository (soft delete by updating status)
func (r *InMemoryAccountRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	account, exists := r.accounts[id]
	if !exists {
		return errors.New("account not found")
	}

	account.Status = domain.StatusDeleted
	account.UpdatedAt = time.Now()
	return nil
}

// List returns all accounts in the repository
func (r *InMemoryAccountRepository) List() ([]*domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	accounts := make([]*domain.Account, 0, len(r.accounts))
	for _, account := range r.accounts {
		accounts = append(accounts, account)
	}

	return accounts, nil
}
