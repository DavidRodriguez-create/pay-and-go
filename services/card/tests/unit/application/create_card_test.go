package application_test

import (
	"testing"

	"github.com/DavidRodriguez-create/pay-and-go/services/card/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/card/domain"
)

// MockCardRepository implements domain.CardRepository for testing
type MockCardRepository struct {
	cards     map[string]*domain.Card
	createErr error
	getErr    error
	deleteErr error
	listErr   error
}

func NewMockCardRepository() *MockCardRepository {
	return &MockCardRepository{
		cards: make(map[string]*domain.Card),
	}
}

func (m *MockCardRepository) Create(card *domain.Card) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.cards[card.ID] = card
	return nil
}

func (m *MockCardRepository) GetByID(id string) (*domain.Card, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	card, exists := m.cards[id]
	if !exists {
		return nil, domain.ErrCardNotFound
	}
	return card, nil
}

func (m *MockCardRepository) GetByCardNumber(cardNumber string) (*domain.Card, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	for _, card := range m.cards {
		if card.CardNumber == cardNumber {
			return card, nil
		}
	}
	return nil, domain.ErrCardNotFound
}

func (m *MockCardRepository) GetByAccountID(accountID string) ([]*domain.Card, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	var cards []*domain.Card
	for _, card := range m.cards {
		if card.AccountID == accountID {
			cards = append(cards, card)
		}
	}
	return cards, nil
}

func (m *MockCardRepository) Delete(id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	card, exists := m.cards[id]
	if !exists {
		return domain.ErrCardNotFound
	}
	card.Deleted = true
	return nil
}

func (m *MockCardRepository) List() ([]*domain.Card, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var cards []*domain.Card
	for _, card := range m.cards {
		cards = append(cards, card)
	}
	return cards, nil
}

// MockAccountCacheRepository implements domain.AccountCacheRepository for testing
type MockAccountCacheRepository struct {
	accounts  map[string]*domain.AccountCache
	upsertErr error
	getErr    error
}

func NewMockAccountCacheRepository() *MockAccountCacheRepository {
	return &MockAccountCacheRepository{
		accounts: make(map[string]*domain.AccountCache),
	}
}

func (m *MockAccountCacheRepository) Upsert(account *domain.AccountCache) error {
	if m.upsertErr != nil {
		return m.upsertErr
	}
	m.accounts[account.ID] = account
	return nil
}

func (m *MockAccountCacheRepository) GetByID(id string) (*domain.AccountCache, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	account, exists := m.accounts[id]
	if !exists {
		return nil, domain.ErrAccountCacheNotFound
	}
	return account, nil
}

func (m *MockAccountCacheRepository) Exists(id string) bool {
	_, exists := m.accounts[id]
	return exists
}

func (m *MockAccountCacheRepository) Delete(id string) error {
	delete(m.accounts, id)
	return nil
}

func (m *MockAccountCacheRepository) List() ([]*domain.AccountCache, error) {
	var accounts []*domain.AccountCache
	for _, account := range m.accounts {
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func TestCreateCard(t *testing.T) {
	t.Run("Successful card creation", func(t *testing.T) {
		cardRepo := NewMockCardRepository()
		accountRepo := NewMockAccountCacheRepository()

		// Setup: Add active account to cache
		accountCache := domain.NewAccountCache("acc-123", domain.AccountStatusActive)
		accountRepo.Upsert(accountCache)

		useCase := application.NewCreateCard(cardRepo, accountRepo)

		req := &application.CreateCardRequest{
			Country:   "US",
			AccountID: "acc-123",
		}

		resp, err := useCase.Execute(req)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if resp.Country != "US" {
			t.Errorf("Expected Country US, got %s", resp.Country)
		}
		if resp.AccountID != "acc-123" {
			t.Errorf("Expected AccountID acc-123, got %s", resp.AccountID)
		}
		if resp.Deleted {
			t.Error("New card should not be deleted")
		}
		if resp.CardNumber == "" {
			t.Error("CardNumber should be generated")
		}
	})

	t.Run("Missing country", func(t *testing.T) {
		cardRepo := NewMockCardRepository()
		accountRepo := NewMockAccountCacheRepository()

		useCase := application.NewCreateCard(cardRepo, accountRepo)

		req := &application.CreateCardRequest{
			Country:   "",
			AccountID: "acc-123",
		}

		_, err := useCase.Execute(req)

		if err != domain.ErrCountryRequired {
			t.Errorf("Expected error %v, got %v", domain.ErrCountryRequired, err)
		}
	})

	t.Run("Missing account ID", func(t *testing.T) {
		cardRepo := NewMockCardRepository()
		accountRepo := NewMockAccountCacheRepository()

		useCase := application.NewCreateCard(cardRepo, accountRepo)

		req := &application.CreateCardRequest{
			Country:   "US",
			AccountID: "",
		}

		_, err := useCase.Execute(req)

		if err != domain.ErrAccountIDRequired {
			t.Errorf("Expected error %v, got %v", domain.ErrAccountIDRequired, err)
		}
	})

	t.Run("Account not found", func(t *testing.T) {
		cardRepo := NewMockCardRepository()
		accountRepo := NewMockAccountCacheRepository()

		useCase := application.NewCreateCard(cardRepo, accountRepo)

		req := &application.CreateCardRequest{
			Country:   "US",
			AccountID: "nonexistent",
		}

		_, err := useCase.Execute(req)

		if err != domain.ErrAccountNotFound {
			t.Errorf("Expected error %v, got %v", domain.ErrAccountNotFound, err)
		}
	})

	t.Run("Account is deleted", func(t *testing.T) {
		cardRepo := NewMockCardRepository()
		accountRepo := NewMockAccountCacheRepository()

		// Setup: Add deleted account to cache
		accountCache := domain.NewAccountCache("acc-123", domain.AccountStatusDeleted)
		accountRepo.Upsert(accountCache)

		useCase := application.NewCreateCard(cardRepo, accountRepo)

		req := &application.CreateCardRequest{
			Country:   "US",
			AccountID: "acc-123",
		}

		_, err := useCase.Execute(req)

		if err != domain.ErrAccountDeleted {
			t.Errorf("Expected error %v, got %v", domain.ErrAccountDeleted, err)
		}
	})

	t.Run("Account is blocked", func(t *testing.T) {
		cardRepo := NewMockCardRepository()
		accountRepo := NewMockAccountCacheRepository()

		// Setup: Add blocked account to cache
		accountCache := domain.NewAccountCache("acc-123", domain.AccountStatusBlocked)
		accountRepo.Upsert(accountCache)

		useCase := application.NewCreateCard(cardRepo, accountRepo)

		req := &application.CreateCardRequest{
			Country:   "US",
			AccountID: "acc-123",
		}

		_, err := useCase.Execute(req)

		if err != domain.ErrAccountInactive {
			t.Errorf("Expected error %v, got %v", domain.ErrAccountInactive, err)
		}
	})

	t.Run("Repository creation error", func(t *testing.T) {
		cardRepo := NewMockCardRepository()
		cardRepo.createErr = domain.ErrCardNotFound // Any error

		accountRepo := NewMockAccountCacheRepository()
		accountCache := domain.NewAccountCache("acc-123", domain.AccountStatusActive)
		accountRepo.Upsert(accountCache)

		useCase := application.NewCreateCard(cardRepo, accountRepo)

		req := &application.CreateCardRequest{
			Country:   "US",
			AccountID: "acc-123",
		}

		_, err := useCase.Execute(req)

		if err == nil {
			t.Error("Expected repository error, got nil")
		}
	})
}
