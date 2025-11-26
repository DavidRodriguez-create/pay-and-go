package controllers

import (
	"net/http"

	"github.com/DavidRodriguez-create/pay-and-go/services/card/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/card/presentation/presenters"
)

// GetCardController handles card retrieval requests
type GetCardController struct {
	useCase   *application.ViewCard
	presenter *presenters.ResponsePresenter
}

// NewGetCardController creates a new GetCardController
func NewGetCardController(
	useCase *application.ViewCard,
	presenter *presenters.ResponsePresenter,
) *GetCardController {
	return &GetCardController{
		useCase:   useCase,
		presenter: presenter,
	}
}

// HandleByID retrieves a card by its ID
func (c *GetCardController) HandleByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	req := &application.GetCardRequest{
		ID: id,
	}

	resp, err := c.useCase.GetByID(req)
	if err != nil {
		c.presenter.HandleError(w, err)
		return
	}

	c.presenter.Success(w, resp, http.StatusOK)
}

// HandleByCardNumber retrieves a card by its card number
func (c *GetCardController) HandleByCardNumber(w http.ResponseWriter, r *http.Request) {
	cardNumber := r.URL.Query().Get("card_number")

	req := &application.GetCardByNumberRequest{
		CardNumber: cardNumber,
	}

	resp, err := c.useCase.GetByCardNumber(req)
	if err != nil {
		c.presenter.HandleError(w, err)
		return
	}

	c.presenter.Success(w, resp, http.StatusOK)
}

// HandleByAccountID retrieves all cards for an account
func (c *GetCardController) HandleByAccountID(w http.ResponseWriter, r *http.Request) {
	accountID := r.URL.Query().Get("account_id")

	req := &application.GetCardsByAccountRequest{
		AccountID: accountID,
	}

	resp, err := c.useCase.GetByAccountID(req)
	if err != nil {
		c.presenter.HandleError(w, err)
		return
	}

	c.presenter.Success(w, resp, http.StatusOK)
}
