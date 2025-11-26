package application

import "github.com/DavidRodriguez-create/pay-and-go/services/card/domain"

// CardToResponse converts a Card domain entity to CardResponse DTO
func CardToResponse(card *domain.Card) *CardResponse {
	if card == nil {
		return nil
	}

	return &CardResponse{
		ID:                card.ID,
		CardNumber:        card.CardNumber,
		Country:           card.Country,
		AccountID:         card.AccountID,
		Deleted:           card.Deleted,
		CreationTimestamp: card.CreationTimestamp,
	}
}

// CardsToResponse converts a slice of Card entities to CardListResponse
func CardsToResponse(cards []*domain.Card) *CardListResponse {
	if cards == nil {
		return &CardListResponse{
			Cards: []*CardResponse{},
			Total: 0,
		}
	}

	responses := make([]*CardResponse, len(cards))
	for i, card := range cards {
		responses[i] = CardToResponse(card)
	}

	return &CardListResponse{
		Cards: responses,
		Total: len(responses),
	}
}
