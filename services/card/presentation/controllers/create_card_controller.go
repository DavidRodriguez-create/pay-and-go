package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/DavidRodriguez-create/pay-and-go/services/card/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/card/presentation/presenters"
)

// CreateCardController handles card creation requests
type CreateCardController struct {
	useCase   *application.CreateCard
	presenter *presenters.ResponsePresenter
}

// NewCreateCardController creates a new CreateCardController
func NewCreateCardController(
	useCase *application.CreateCard,
	presenter *presenters.ResponsePresenter,
) *CreateCardController {
	return &CreateCardController{
		useCase:   useCase,
		presenter: presenter,
	}
}

// Handle processes card creation requests
func (c *CreateCardController) Handle(w http.ResponseWriter, r *http.Request) {
	var req application.CreateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.presenter.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	resp, err := c.useCase.Execute(&req)
	if err != nil {
		c.presenter.HandleError(w, err)
		return
	}

	c.presenter.Success(w, resp, http.StatusCreated)
}
