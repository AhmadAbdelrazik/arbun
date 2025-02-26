package app

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"net/http"
)

func (app *Application) PostOrder(w http.ResponseWriter, r *http.Request) {
	customer := app.contextGetUser(r).(domain.User)

	var input struct {
		DeliveryAddress string
		MobilePhone     string
		PaymentType     string
	}

	readJSON(w, r, &input)

	app.services.Orders.CreateOrder(customer.ID, app.services.Carts)

}

func (app *Application) PostOrderConfirmation(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) GetOrder(w http.ResponseWriter, r *http.Request) {
	// customer := app.contextGetUser(r).(domain.Customer)
}

func (app *Application) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	// customer := app.contextGetUser(r).(domain.Customer)
}
