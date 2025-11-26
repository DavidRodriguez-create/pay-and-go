package application

import "github.com/DavidRodriguez-create/pay-and-go/services/card/domain"

// ViewCard handles card retrieval use cases
type ViewCard struct {
	cardRepo domain.CardRepository
}

// NewViewCard creates a new ViewCard use case
func NewViewCard(cardRepo domain.CardRepository) *ViewCard {
	return &ViewCard{
		cardRepo: cardRepo,
	}
}

// GetByID retrieves a card by its ID
func (uc *ViewCard) GetByID(req *GetCardRequest) (*CardResponse, error) {
	if req.ID == "" {
		return nil, domain.ErrCardIDRequired
	}

	card, err := uc.cardRepo.GetByID(req.ID)
	if err != nil {
		return nil, domain.ErrCardNotFound
	}

	return CardToResponse(card), nil
}

// GetByCardNumber retrieves a card by its card number
func (uc *ViewCard) GetByCardNumber(req *GetCardByNumberRequest) (*CardResponse, error) {
	if req.CardNumber == "" {
		return nil, domain.ErrCardNumberRequired
	}

	card, err := uc.cardRepo.GetByCardNumber(req.CardNumber)
	if err != nil {
		return nil, domain.ErrCardNotFound
	}

	return CardToResponse(card), nil
}

// GetByAccountID retrieves all cards for an account
func (uc *ViewCard) GetByAccountID(req *GetCardsByAccountRequest) (*CardListResponse, error) {
	if req.AccountID == "" {
		return nil, domain.ErrAccountIDRequired
	}

	cards, err := uc.cardRepo.GetByAccountID(req.AccountID)
	if err != nil {
		return nil, err
	}

	return CardsToResponse(cards), nil
}

// ListCards retrieves all cards
type ListCards struct {
	cardRepo domain.CardRepository
}

// NewListCards creates a new ListCards use case
func NewListCards(cardRepo domain.CardRepository) *ListCards {
	return &ListCards{
		cardRepo: cardRepo,
	}
}

// Execute retrieves all cards
func (uc *ListCards) Execute() (*CardListResponse, error) {
	cards, err := uc.cardRepo.List()
	if err != nil {
		return nil, err
	}

	return CardsToResponse(cards), nil
}
