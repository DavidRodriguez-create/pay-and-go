package application

import "github.com/DavidRodriguez-create/pay-and-go/services/card/domain"

// CardService orchestrates card-related use cases
type CardService struct {
	CreateCard *CreateCard
	DeleteCard *DeleteCard
	ViewCard   *ViewCard
	ListCards  *ListCards
}

// NewCardService creates a new CardService with all use cases
func NewCardService(
	cardRepo domain.CardRepository,
	accountRepo domain.AccountCacheRepository,
) *CardService {
	return &CardService{
		CreateCard: NewCreateCard(cardRepo, accountRepo),
		DeleteCard: NewDeleteCard(cardRepo),
		ViewCard:   NewViewCard(cardRepo),
		ListCards:  NewListCards(cardRepo),
	}
}
