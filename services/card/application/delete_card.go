package application

import "github.com/DavidRodriguez-create/pay-and-go/services/card/domain"

// DeleteCard handles card deletion use case (soft delete)
type DeleteCard struct {
	cardRepo domain.CardRepository
}

// NewDeleteCard creates a new DeleteCard use case
func NewDeleteCard(cardRepo domain.CardRepository) *DeleteCard {
	return &DeleteCard{
		cardRepo: cardRepo,
	}
}

// Execute marks a card as deleted
func (uc *DeleteCard) Execute(req *DeleteCardRequest) error {
	if req.ID == "" {
		return domain.ErrCardIDRequired
	}

	// Retrieve the card
	card, err := uc.cardRepo.GetByID(req.ID)
	if err != nil {
		return domain.ErrCardNotFound
	}

	// Mark as deleted
	if err := card.Delete(); err != nil {
		return err
	}

	// Persist the change
	return uc.cardRepo.Delete(req.ID)
}
