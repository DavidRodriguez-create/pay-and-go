package application

import (
	"time"

	"github.com/DavidRodriguez-create/pay-and-go/services/card/domain"
	"github.com/google/uuid"
)

// CreateCard handles card creation use case
type CreateCard struct {
	cardRepo    domain.CardRepository
	accountRepo domain.AccountCacheRepository
}

// NewCreateCard creates a new CreateCard use case
func NewCreateCard(cardRepo domain.CardRepository, accountRepo domain.AccountCacheRepository) *CreateCard {
	return &CreateCard{
		cardRepo:    cardRepo,
		accountRepo: accountRepo,
	}
}

// Execute creates a new card after validating the account
func (uc *CreateCard) Execute(req *CreateCardRequest) (*CardResponse, error) {
	// Validate input
	if req.Country == "" {
		return nil, domain.ErrCountryRequired
	}
	if req.AccountID == "" {
		return nil, domain.ErrAccountIDRequired
	}

	// Check if account exists and is active
	account, err := uc.accountRepo.GetByID(req.AccountID)
	if err != nil {
		return nil, domain.ErrAccountNotFound
	}

	if account.IsDeleted() {
		return nil, domain.ErrAccountDeleted
	}

	if !account.IsActive() {
		return nil, domain.ErrAccountInactive
	}

	// Generate card number (simple format: COUNTRY-UUID)
	cardNumber := req.Country + "-" + uuid.New().String()[:8]

	// Create card entity
	card, err := domain.NewCard(
		uuid.New().String(),
		cardNumber,
		req.Country,
		req.AccountID,
		time.Now(),
	)
	if err != nil {
		return nil, err
	}

	// Persist card
	if err := uc.cardRepo.Create(card); err != nil {
		return nil, err
	}

	return CardToResponse(card), nil
}
