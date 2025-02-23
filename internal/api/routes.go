package handlers

import (
	"AhmadAbdelrazik/arbun/internal/domain/user"
	"net/http"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /products", app.IsUser(user.TypeAdmin, app.PostProduct))
	mux.HandleFunc("GET /products/{id}", app.GetProduct)
	mux.HandleFunc("GET /products", app.GetAllProducts)
	mux.HandleFunc("PATCH /products/{id}", app.IsUser(user.TypeAdmin, app.PatchProduct))
	mux.HandleFunc("DELETE /products/{id}", app.IsUser(user.TypeAdmin, app.DeleteProduct))

	mux.HandleFunc("POST /signup", app.PostSignup)
	mux.HandleFunc("POST /login", app.PostLogin)
	mux.HandleFunc("POST /logout", app.PostLogout)

	mux.HandleFunc("GET /cart", app.IsUser(user.TypeCustomer, app.GetCart))
	mux.HandleFunc("POST /cart", app.IsUser(user.TypeCustomer, app.PostCartItems))
	mux.HandleFunc("PATCH /cart", app.IsUser(user.TypeCustomer, app.PostCartItems))
	mux.HandleFunc("DELETE /cart", app.IsUser(user.TypeCustomer, app.DeleteCartItem))
	return app.recoverPanic(securityHeaders(mux))
}
