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

func (app *Application) contextSetUser(r *http.Request, user domain.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContext, user)
	return r.WithContext(ctx)
}

func (app *Application) contextGetUser(r *http.Request) domain.User {
	u := r.Context().Value(UserContext)
	if u == nil {
		panic("missing value in request context")
	}

	return u.(domain.User)
}

func (app *Application) contextGetAdmin(r *http.Request) domain.Admin {
	user := app.contextGetUser(r)

	if user.Type != domain.TypeAdmin {
		panic("expeceted user with admin type")
	}

	var admin domain.Admin
	admin.User = user

	return admin
}

func (app *Application) contextGetCustomer(r *http.Request) domain.Customer {
	user := app.contextGetUser(r)

	if user.Type != domain.TypeCustomer {
		panic("expeceted user with customer type")
	}

	var customer domain.Customer
	customer.User = user

	return customer
}
