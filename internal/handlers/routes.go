package handlers

import "net/http"

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /products", app.IsAdmin(app.PostProduct))
	mux.HandleFunc("GET /products/{id}", app.GetProduct)
	mux.HandleFunc("GET /products", app.GetAllProducts)
	mux.HandleFunc("PATCH /products/{id}", app.IsAdmin(app.PatchProduct))
	mux.HandleFunc("DELETE /products/{id}", app.IsAdmin(app.DeleteProduct))

	mux.HandleFunc("POST /signup", app.PostAdminSignup)
	mux.HandleFunc("POST /login", app.PostAdminLogin)
	mux.HandleFunc("POST /logout", app.PostAdminLogout)

	return app.recoverPanic(securityHeaders(mux))
}
