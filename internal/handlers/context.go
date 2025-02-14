package handlers

import (
	"AhmadAbdelrazik/arbun/internal/repository"
	"context"
	"net/http"
)

type ContextKey string

const (
	AdminContext ContextKey = "admin"
)

func (app *Application) contextSetAdmin(r *http.Request, admin repository.Admin) *http.Request {
	ctx := context.WithValue(r.Context(), AdminContext, admin)
	return r.WithContext(ctx)
}

func (app *Application) contextGetUser(r *http.Request) repository.Admin {
	admin, ok := r.Context().Value(AdminContext).(repository.Admin)
	if !ok {
		panic("missing value in request context")
	}

	return admin
}
