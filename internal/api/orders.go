package handlers

import (
	"AhmadAbdelrazik/arbun/internal/domain/customer"
	"net/http"
)

func (app *Application) PostOrder(w http.ResponseWriter, r *http.Request) {
	customer := app.contextGetUser(r).(customer.Customer)
	var input struct {
		DeliveryAddress string
		MobilePhone     string
		PaymentType     string
	}

	readJSON(w, r, &input)

	app.services.Orders.CreateOrder(customer.ID)

}

func (app *Application) GetOrder(w http.ResponseWriter, r *http.Request) {
	customer := app.contextGetUser(r).(customer.Customer)
}

func (app *Application) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	customer := app.contextGetUser(r).(customer.Customer)

}
