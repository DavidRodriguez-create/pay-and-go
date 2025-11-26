package infrastructure

import (
	"sync"

	"github.com/DavidRodriguez-create/pay-and-go/services/card/domain"
)

// InMemoryCardRepository implements CardRepository with in-memory storage
type InMemoryCardRepository struct {
	cards map[string]*domain.Card
	mu    sync.RWMutex
}

// NewInMemoryCardRepository creates a new in-memory card repository
func NewInMemoryCardRepository() *InMemoryCardRepository {
	return &InMemoryCardRepository{
		cards: make(map[string]*domain.Card),
	}
}

// Create stores a new card
func (r *InMemoryCardRepository) Create(card *domain.Card) error {
	if card == nil {
		return domain.ErrCardNotFound
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.cards[card.ID] = card
	return nil
}

// GetByID retrieves a card by its ID
func (r *InMemoryCardRepository) GetByID(id string) (*domain.Card, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	card, exists := r.cards[id]
	if !exists {
		return nil, domain.ErrCardNotFound
	}

	return card, nil
}

// GetByCardNumber retrieves a card by its card number
func (r *InMemoryCardRepository) GetByCardNumber(cardNumber string) (*domain.Card, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, card := range r.cards {
		if card.CardNumber == cardNumber {
			return card, nil
		}
	}

	return nil, domain.ErrCardNotFound
}

// GetByAccountID retrieves all cards for a specific account
func (r *InMemoryCardRepository) GetByAccountID(accountID string) ([]*domain.Card, error) {
	if accountID == "" {
		return nil, domain.ErrAccountIDRequired
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var cards []*domain.Card
	for _, card := range r.cards {
		if card.AccountID == accountID {
			cards = append(cards, card)
		}
	}

	return cards, nil
}

// Delete marks a card as deleted (soft delete)
func (r *InMemoryCardRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	card, exists := r.cards[id]
	if !exists {
		return domain.ErrCardNotFound
	}

	card.Deleted = true
	return nil
}

// List retrieves all cards
func (r *InMemoryCardRepository) List() ([]*domain.Card, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cards := make([]*domain.Card, 0, len(r.cards))
	for _, card := range r.cards {
		cards = append(cards, card)
	}

	return cards, nil
}
