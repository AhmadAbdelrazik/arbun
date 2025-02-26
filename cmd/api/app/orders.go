package app

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"net/http"
)

func (app *Application) postOrder(w http.ResponseWriter, r *http.Request) {
	customer := app.contextGetUser(r).(domain.User)

	var input struct {
		DeliveryAddress string
		MobilePhone     string
		PaymentType     string
	}

	readJSON(w, r, &input)

	app.services.Orders.CreateOrder(customer.ID, app.services.Carts)

}

func (app *Application) postOrderConfirmation(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) getOrder(w http.ResponseWriter, r *http.Request) {
	// customer := app.contextGetUser(r).(domain.Customer)
}

func (app *Application) getAllOrders(w http.ResponseWriter, r *http.Request) {
	// customer := app.contextGetUser(r).(domain.Customer)
}
