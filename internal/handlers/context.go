package handlers

import (
	"AhmadAbdelrazik/arbun/internal/repository"
	"context"
	"net/http"
)

type ContextKey string

const (
	UserContext ContextKey = "admin"
)

func (app *Application) contextSetUser(r *http.Request, user repository.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContext, user)
	return r.WithContext(ctx)
}

func (app *Application) contextGetUser(r *http.Request) repository.User {
	user, ok := r.Context().Value(UserContext).(repository.User)
	if !ok {
		panic("missing value in request context")
	}

	return user
}
