package app

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"AhmadAbdelrazik/arbun/internal/services"
	"errors"
	"net/http"
)

func (app *Application) getCart(w http.ResponseWriter, r *http.Request) {
	customer := app.contextGetCustomer(r)

	cart, err := app.services.Carts.GetCart(customer.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"cart": cart}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) postCartItems(w http.ResponseWriter, r *http.Request) {
	customer := app.contextGetCustomer(r)

	var input struct {
		Items []domain.CartItem `json:"items"`
	}

	err := readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	cart, err := app.services.Carts.UpdateItems(customer.ID, input.Items)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrProductNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"cart": cart}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) deleteCartItem(w http.ResponseWriter, r *http.Request) {
	customer := app.contextGetCustomer(r)

	var input struct {
		ProductID int64 `json:"product_id"`
	}

	err := readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	cart, err := app.services.Carts.DeleteItem(customer.ID, input.ProductID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"cart": cart}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
