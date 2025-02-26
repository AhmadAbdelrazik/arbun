package app

import (
	"net/http"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /products", app.isAdmin(app.postProduct))
	mux.HandleFunc("GET /products/{id}", app.getProduct)
	mux.HandleFunc("GET /products", app.getAllProducts)
	mux.HandleFunc("PATCH /products/{id}", app.isAdmin(app.patchProduct))
	mux.HandleFunc("DELETE /products/{id}", app.isAdmin(app.deleteProduct))

	mux.HandleFunc("POST /signup", app.postSignup)
	mux.HandleFunc("POST /login", app.postLogin)
	mux.HandleFunc("POST /logout", app.postLogout)

	mux.HandleFunc("GET /cart", app.isCustomer(app.getCart))
	mux.HandleFunc("POST /cart", app.isCustomer(app.postCartItems))
	mux.HandleFunc("PATCH /cart", app.isCustomer(app.postCartItems))
	mux.HandleFunc("DELETE /cart", app.isCustomer(app.deleteCartItem))
	return app.recoverPanic(securityHeaders(mux))
}
