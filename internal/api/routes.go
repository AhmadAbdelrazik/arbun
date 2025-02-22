package handlers

import (
	"AhmadAbdelrazik/arbun/internal/services"
	"net/http"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /products", app.IsUser(services.TypeAdmin, app.PostProduct))
	mux.HandleFunc("GET /products/{id}", app.GetProduct)
	mux.HandleFunc("GET /products", app.GetAllProducts)
	mux.HandleFunc("PATCH /products/{id}", app.IsUser(services.TypeAdmin, app.PatchProduct))
	mux.HandleFunc("DELETE /products/{id}", app.IsUser(services.TypeAdmin, app.DeleteProduct))

	mux.HandleFunc("POST /signup", app.PostSignup)
	mux.HandleFunc("POST /login", app.PostLogin)
	mux.HandleFunc("POST /logout", app.PostLogout)

	return app.recoverPanic(securityHeaders(mux))
}
