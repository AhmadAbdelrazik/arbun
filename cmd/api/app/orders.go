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
		DeliveryAddress domain.Address       `json:"address"`
		MobilePhone     domain.MobilePhone   `json:"mobile_phone"`
		PaymentMethod   domain.PaymentMethod `json:"payment_method"`
		Properties      map[string]string    `json:"properties"`
	}

	err := readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := input.PaymentMethod.Validate()
	if v != nil {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	customer.Address = input.DeliveryAddress
	customer.MobilePhone = input.MobilePhone

	v = customer.Validate()
	if v != nil {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	switch input.PaymentMethod {
	case domain.PaymentCash:
		app.postCashOrder(w, r, customer)
	case domain.PaymentCard:
		app.postCardOrder(w, r, customer)
	default:
		panic("unimplemented payment method")
	}
}

func (app *Application) postCashOrder(w http.ResponseWriter, r *http.Request, customer domain.Customer) {
	order, err := app.services.Orders.CreateCashOrder(customer)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusCreated, envelope{"order": order}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) postCardOrder(w http.ResponseWriter, r *http.Request, customer domain.Customer) {
	redirectURL, err := app.services.Orders.CreateCardOrder(customer)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
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
