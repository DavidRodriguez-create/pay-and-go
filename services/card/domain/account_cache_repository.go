package domain

// AccountCacheRepository defines the interface for account cache persistence
type AccountCacheRepository interface {
	// Upsert creates or updates an account cache entry
	Upsert(account *AccountCache) error

	// GetByID retrieves an account cache by its ID
	GetByID(id string) (*AccountCache, error)

	// Exists checks if an account exists in cache
	Exists(id string) bool

	// Delete removes an account from cache
	Delete(id string) error

	// List retrieves all cached accounts
	List() ([]*AccountCache, error)
}
