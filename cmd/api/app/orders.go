package app

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"AhmadAbdelrazik/arbun/internal/services"
	"errors"
	"net/http"
)

func (app *Application) postOrder(w http.ResponseWriter, r *http.Request) {
	customer := app.contextGetCustomer(r)

	var input struct {
		DeliveryAddress domain.Address `json:"address"`
		MobilePhone     string         `json:"mobile_phone"`
	}

	err := readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	customer.Address = input.DeliveryAddress
	customer.MobilePhone = domain.MobilePhone(input.MobilePhone)

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

	err = writeJSON(w, http.StatusCreated, envelope{"order": order}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) getOrder(w http.ResponseWriter, r *http.Request) {
	customer := app.contextGetCustomer(r)
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

func (app *Application) patchOrder(w http.ResponseWriter, r *http.Request) {
	// read order ID parameter
	id, err := readIDParam(r, "id")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// parse status name from body
	var input struct {
		Status string `json:"status"`
	}
	err = readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// validate status
	status := domain.OrderStatus(input.Status)
	v := status.Validate()
	if v != nil {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// change order status
	err = app.services.Orders.ChangeOrderStatus(id, status)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrOrderNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// response
	err = writeJSON(w, http.StatusOK, envelope{"status": "order status updated"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) getAllOrders(w http.ResponseWriter, r *http.Request) {
	customer := app.contextGetCustomer(r)

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
