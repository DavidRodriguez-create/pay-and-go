package application

import "time"

// CreateCardRequest represents the input for creating a card
type CreateCardRequest struct {
	Country   string `json:"country"`
	AccountID string `json:"account_id"`
}

// CardResponse represents the output for card operations
type CardResponse struct {
	ID                string    `json:"id"`
	CardNumber        string    `json:"card_number"`
	Country           string    `json:"country"`
	AccountID         string    `json:"account_id"`
	Deleted           bool      `json:"deleted"`
	CreationTimestamp time.Time `json:"creation_timestamp"`
}

// DeleteCardRequest represents the input for deleting a card
type DeleteCardRequest struct {
	ID string `json:"id"`
}

// GetCardRequest represents the input for retrieving a card
type GetCardRequest struct {
	ID string `json:"id"`
}

// GetCardByNumberRequest represents the input for retrieving a card by number
type GetCardByNumberRequest struct {
	CardNumber string `json:"card_number"`
}

// GetCardsByAccountRequest represents the input for retrieving cards by account
type GetCardsByAccountRequest struct {
	AccountID string `json:"account_id"`
}

// CardListResponse represents a list of cards
type CardListResponse struct {
	Cards []*CardResponse `json:"cards"`
	Total int             `json:"total"`
}
