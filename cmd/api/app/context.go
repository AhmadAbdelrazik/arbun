package app

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"context"
	"net/http"
)

type ContextKey string

const (
	UserContext ContextKey = "user"
)

func (app *Application) contextSetUser(r *http.Request, user domain.IUser) *http.Request {
	ctx := context.WithValue(r.Context(), UserContext, user)
	return r.WithContext(ctx)
}

func (app *Application) contextGetUser(r *http.Request) domain.IUser {
	u := r.Context().Value(UserContext)
	if u == nil {
		panic("missing value in request context")
	}

	return u.(domain.IUser)
}
