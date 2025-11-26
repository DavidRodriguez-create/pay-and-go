package controllers

import (
	"net/http"

	"github.com/DavidRodriguez-create/pay-and-go/services/card/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/card/presentation/presenters"
)

// ListCardsController handles card listing requests
type ListCardsController struct {
	useCase   *application.ListCards
	presenter *presenters.ResponsePresenter
}

// NewListCardsController creates a new ListCardsController
func NewListCardsController(
	useCase *application.ListCards,
	presenter *presenters.ResponsePresenter,
) *ListCardsController {
	return &ListCardsController{
		useCase:   useCase,
		presenter: presenter,
	}
}

// Handle retrieves all cards
func (c *ListCardsController) Handle(w http.ResponseWriter, r *http.Request) {
	resp, err := c.useCase.Execute()
	if err != nil {
		c.presenter.HandleError(w, err)
		return
	}

	c.presenter.Success(w, resp, http.StatusOK)
}
