package app

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"AhmadAbdelrazik/arbun/internal/services"
	"errors"
	"net/http"
)

func (app *Application) postOrder(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r).(domain.User)

	var input struct {
		DeliveryAddress domain.Address `json:"address"`
		MobilePhone     string         `json:"mobile_phone"`
	}

	readJSON(w, r, &input)

	customer := domain.Customer{
		User:        user,
		Address:     input.DeliveryAddress,
		MobilePhone: domain.MobilePhone(input.MobilePhone),
	}

	v := customer.Validate()
	if v != nil {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	order, err := app.services.Orders.CreateOrder(customer, app.services.Carts)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"order": order}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) getOrder(w http.ResponseWriter, r *http.Request) {
	customer := app.contextGetUser(r).(domain.Customer)
	orderID, err := readIDParam(r, "id")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	order, err := app.services.Orders.GetOrder(customer, orderID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrOrderNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"order": order}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) getAllOrders(w http.ResponseWriter, r *http.Request) {
	customer := app.contextGetUser(r).(domain.Customer)

	orders, err := app.services.Orders.GetAllUserOrders(customer)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, envelope{"orders": orders}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
