package handlers

import "net/http"

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /products", app.PostProduct)
	mux.HandleFunc("GET /products/{id}", app.GetProduct)
	mux.HandleFunc("GET /products", app.GetAllProducts)

	return mux
}
