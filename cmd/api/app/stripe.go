package app

import (
	"io"
	"net/http"
)

func (app *Application) PostOrderConfirmation(w http.ResponseWriter, r *http.Request) {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer r.Body.Close()

	header := r.Header.Get("Stripe-Signature")

	err = app.services.Stripe.ConfirmOrder(payload, header)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	writeJSON(w, http.StatusOK, envelope{}, nil)
}
