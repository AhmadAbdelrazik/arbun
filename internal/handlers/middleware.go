package handlers

import (
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

func (app *Application) IsUser(userType string, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := GetAuthToken(r)
		if err != nil {
			app.authenticationErrorResponse(w, r)
			return
		}

		admin, err := app.services.Users.GetAuthToken(token.Plaintext, userType)
		if err != nil {
			switch {
			case errors.Is(err, services.ErrInvalidAuthToken):
				app.authenticationErrorResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		_ = admin

		next.ServeHTTP(w, r)
	})
}

func (app *Application) IsAdmin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Check for token
		token, err := GetAuthToken(r)
		if err != nil {
			app.authenticationErrorResponse(w, r)
			return
		}
		// 2. register admin in the request context
		// TODO: Find a better way to deal with scope rather
		// than calling repository module directly
		admin, err := app.services.Admins.GetAdminbyAuthToken(token.Plaintext)
		if err != nil {
			switch {
			case errors.Is(err, services.ErrInvalidAuthToken):
				app.authenticationErrorResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetAdmin(r, admin)

		next.ServeHTTP(w, r)
	})
}
