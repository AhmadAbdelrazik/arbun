package app

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"AhmadAbdelrazik/arbun/internal/services"
	"errors"
	"fmt"
	"net/http"
)

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")

				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *Application) IsUser(userType string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := GetAuthToken(r)
		if err != nil {
			app.authenticationErrorResponse(w, r)
			return
		}

		user, err := app.services.Users.GetUserByToken(token.Plaintext)
		if err != nil {
			switch {
			case errors.Is(err, services.ErrInvalidAuthToken):
				app.authenticationErrorResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		if user.Type != userType {

		}

		r = app.contextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}

func (app *Application) IsAdmin(next http.HandlerFunc) http.HandlerFunc {
	return app.IsUser(domain.TypeAdmin, next)
}

func (app *Application) IsCustomer(next http.HandlerFunc) http.HandlerFunc {
	return app.IsUser(domain.TypeCustomer, next)
}
