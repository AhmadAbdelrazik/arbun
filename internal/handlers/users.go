package handlers

import (
	"AhmadAbdelrazik/arbun/internal/services"
	"errors"
	"net/http"
)

func (app *Application) PostSignup(w http.ResponseWriter, r *http.Request) {
	// receive and validate JSON input
	var input struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		UserType string `json:"type"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	token, err := app.services.Users.Signup(
		input.FullName,
		input.Email,
		input.Password,
		input.UserType,
	)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmailAlreadyTaken):
			app.badRequestResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	SetAuthTokenCookie(w, token)

	err = app.writeJSON(w, http.StatusCreated, envelope{"message": "created successfully"}, nil)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
}

func (app *Application) PostLogin(w http.ResponseWriter, r *http.Request) {
	// receive and validate JSON input
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		UserType string `json:"type"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	token, err := app.services.Users.Login(
		input.Email,
		input.Password,
		input.UserType,
	)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			app.authenticationErrorResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	SetAuthTokenCookie(w, token)

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "logged in successfully"}, nil)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
}

func (app *Application) PostLogout(w http.ResponseWriter, r *http.Request) {
	token, err := GetAuthToken(r)
	if err != nil {
		err = app.writeJSON(w, http.StatusOK, envelope{"message": "logged out successfully"}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}

	err = app.services.Users.Logout(token, services.TypeAdmin)
	if err != nil {
		err = app.writeJSON(w, http.StatusOK, envelope{"message": "logged out successfully"}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "logged out successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
