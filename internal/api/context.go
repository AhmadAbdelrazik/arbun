package handlers

import (
	"AhmadAbdelrazik/arbun/internal/models"
	"context"
	"net/http"
)

type ContextKey string

const (
	UserContext ContextKey = "user"
)

func (app *Application) contextSetUser(r *http.Request, user models.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContext, user)
	return r.WithContext(ctx)
}

func (app *Application) contextGetUser(r *http.Request) models.User {
	user := r.Context().Value(UserContext)
	if user == nil {
		panic("missing value in request context")
	}

	return user.(models.User)
}
