package controllers

import (
	"net/http"

	"github.com/DavidRodriguez-create/pay-and-go/services/card/application"
	"github.com/DavidRodriguez-create/pay-and-go/services/card/presentation/presenters"
)

// DeleteCardController handles card deletion requests
type DeleteCardController struct {
	useCase   *application.DeleteCard
	presenter *presenters.ResponsePresenter
}

// NewDeleteCardController creates a new DeleteCardController
func NewDeleteCardController(
	useCase *application.DeleteCard,
	presenter *presenters.ResponsePresenter,
) *DeleteCardController {
	return &DeleteCardController{
		useCase:   useCase,
		presenter: presenter,
	}
}

// Handle processes card deletion requests
func (c *DeleteCardController) Handle(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	req := &application.DeleteCardRequest{
		ID: id,
	}

	if err := c.useCase.Execute(req); err != nil {
		c.presenter.HandleError(w, err)
		return
	}

	c.presenter.Success(w, map[string]string{"message": "Card deleted successfully"}, http.StatusOK)
}
