package handlers

import "net/http"

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	return mux
}
