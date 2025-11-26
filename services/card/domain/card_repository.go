package domain

// CardRepository defines the interface for card data persistence
type CardRepository interface {
	// Create stores a new card
	Create(card *Card) error

	// GetByID retrieves a card by its ID
	GetByID(id string) (*Card, error)

	// GetByCardNumber retrieves a card by its card number
	GetByCardNumber(cardNumber string) (*Card, error)

	// GetByAccountID retrieves all cards for a specific account
	GetByAccountID(accountID string) ([]*Card, error)

	// Delete marks a card as deleted (soft delete)
	Delete(id string) error

	// List retrieves all cards (including deleted ones)
	List() ([]*Card, error)
}
