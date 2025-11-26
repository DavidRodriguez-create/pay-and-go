package infrastructure

import (
	"sync"

	"github.com/DavidRodriguez-create/pay-and-go/services/card/domain"
)

// InMemoryAccountCacheRepository implements AccountCacheRepository with in-memory storage
type InMemoryAccountCacheRepository struct {
	accounts map[string]*domain.AccountCache
	mu       sync.RWMutex
}

// NewInMemoryAccountCacheRepository creates a new in-memory account cache repository
func NewInMemoryAccountCacheRepository() *InMemoryAccountCacheRepository {
	return &InMemoryAccountCacheRepository{
		accounts: make(map[string]*domain.AccountCache),
	}
}

// Upsert creates or updates an account cache entry
func (r *InMemoryAccountCacheRepository) Upsert(account *domain.AccountCache) error {
	if account == nil {
		return domain.ErrAccountNotFound
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.accounts[account.ID] = account
	return nil
}

// GetByID retrieves an account cache by its ID
func (r *InMemoryAccountCacheRepository) GetByID(id string) (*domain.AccountCache, error) {
	if id == "" {
		return nil, domain.ErrAccountCacheNotFound
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	account, exists := r.accounts[id]
	if !exists {
		return nil, domain.ErrAccountCacheNotFound
	}

	return account, nil
}

// Exists checks if an account exists in cache
func (r *InMemoryAccountCacheRepository) Exists(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.accounts[id]
	return exists
}

// Delete removes an account from cache
func (r *InMemoryAccountCacheRepository) Delete(id string) error {
	if id == "" {
		return domain.ErrAccountCacheNotFound
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.accounts[id]
	if !exists {
		return domain.ErrAccountCacheNotFound
	}

	delete(r.accounts, id)
	return nil
}

// List retrieves all cached accounts
func (r *InMemoryAccountCacheRepository) List() ([]*domain.AccountCache, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	accounts := make([]*domain.AccountCache, 0, len(r.accounts))
	for _, account := range r.accounts {
		accounts = append(accounts, account)
	}

	return accounts, nil
}
