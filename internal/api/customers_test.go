package handlers

import (
	"AhmadAbdelrazik/arbun/internal/assert"
	"AhmadAbdelrazik/arbun/internal/services"
	"fmt"
	"net/http"
	"testing"
)

func TestCustomerSignup(t *testing.T) {
	ts := NewTestClient()
	defer ts.Close()

	t.Run("valid signups", func(t *testing.T) {
		tests := []struct {
			FullName string `json:"full_name"`
			Email    string `json:"email"`
			Password string `json:"password"`
			UserType string `json:"type"`
		}{
			{
				FullName: "customer1",
				Email:    "customer1@example.com",
				Password: "password1",
				UserType: services.TypeCustomer,
			},
			{
				FullName: "customer2",
				Email:    "customer2@example.com",
				Password: "password2",
				UserType: services.TypeCustomer,
			},
		}

		for i, tt := range tests {
			t.Run(fmt.Sprintf("customer%d", i), func(t *testing.T) {
				res, err := ts.Post("/signup", tt)
				assert.Nil(t, err)
				assert.Equal(t, res.StatusCode, http.StatusCreated)

				authCookie := ts.GetCookie(res, AuthCookie)

				assert.Nil(t, authCookie.Valid())
				assert.True(t, len(authCookie.Value) == 26)
			})
		}
	})

	t.Run("invalid signups", func(t *testing.T) {
		t.Run("missing full name", func(t *testing.T) {
			requestBody := struct {
				FullName string `json:"full_name"`
				Email    string `json:"email"`
				Password string `json:"password"`
				UserType string `json:"type"`
			}{
				FullName: "",
				Email:    "customer4@example.com",
				Password: "password4",
				UserType: services.TypeCustomer,
			}

			res, err := ts.Post("/signup", requestBody)
			assert.Nil(t, err)
			var responseBody struct {
				Error struct {
					FullName string `json:"full_name"`
				} `json:"error"`
			}

			err = ts.ReadResponseBody(res, &responseBody)
			assert.Nil(t, err)

			assert.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)
			assert.Equal(t, responseBody.Error.FullName, "must not be empty")

			authCookie := ts.GetCookie(res, AuthCookie)
			assert.True(t, authCookie == nil)
		})
		t.Run("missing email and password", func(t *testing.T) {
			requestBody := struct {
				FullName string `json:"full_name"`
				Email    string `json:"email"`
				Password string `json:"password"`
				UserType string `json:"type"`
			}{
				FullName: "customer1",
				Email:    "",
				Password: "",
				UserType: services.TypeCustomer,
			}

			res, err := ts.Post("/signup", requestBody)
			assert.Nil(t, err)
			var responseBody struct {
				Error struct {
					Email    string `json:"email"`
					Password string `json:"password"`
				} `json:"error"`
			}

			err = ts.ReadResponseBody(res, &responseBody)
			assert.Nil(t, err)

			assert.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)
			assert.Equal(t, responseBody.Error.Email, "must not be empty")
			assert.Equal(t, responseBody.Error.Password, "must not be empty")

			authCookie := ts.GetCookie(res, AuthCookie)
			assert.True(t, authCookie == nil)
		})
		t.Run("duplicate email", func(t *testing.T) {
			requestBody := struct {
				FullName string `json:"full_name"`
				Email    string `json:"email"`
				Password string `json:"password"`
				UserType string `json:"type"`
			}{
				FullName: "customer1",
				Email:    "customer1@example.com",
				Password: "password1",
				UserType: services.TypeCustomer,
			}

			res, err := ts.Post("/signup", requestBody)
			assert.Nil(t, err)
			assert.Equal(t, res.StatusCode, http.StatusBadRequest)

			authCookie := ts.GetCookie(res, AuthCookie)
			assert.True(t, authCookie == nil)
		})
	})
}

func TestCustomerLogin(t *testing.T) {
	ts := NewTestClient()
	defer ts.Close()

	customerLoginTestSetup(ts)

	t.Run("valid logins", func(t *testing.T) {
		tests := []struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			UserType string `json:"type"`
		}{
			{
				Email:    "customer1@example.com",
				Password: "password1",
				UserType: services.TypeCustomer,
			},
			{
				Email:    "customer2@example.com",
				Password: "password2",
				UserType: services.TypeCustomer,
			},
		}

		for i, tt := range tests {
			t.Run(fmt.Sprintf("customer%d", i), func(t *testing.T) {
				res, err := ts.Post("/login", tt)
				assert.Nil(t, err)
				assert.Equal(t, res.StatusCode, http.StatusOK)

				authCookie := ts.GetCookie(res, AuthCookie)

				assert.Nil(t, authCookie.Valid())
				assert.True(t, len(authCookie.Value) == 26)
			})
		}
	})

	t.Run("invalid logins", func(t *testing.T) {
		t.Run("wrong email", func(t *testing.T) {
			responseBody := struct {
				Email    string `json:"email"`
				Password string `json:"password"`
				UserType string `json:"type"`
			}{
				Email:    "customer3@example.com",
				Password: "password1",
				UserType: services.TypeCustomer,
			}

			res, err := ts.Post("/login", responseBody)
			assert.Nil(t, err)
			assert.Equal(t, res.StatusCode, http.StatusForbidden)

			authCookie := ts.GetCookie(res, AuthCookie)
			assert.True(t, authCookie == nil)
		})
		t.Run("wrong password", func(t *testing.T) {
			responseBody := struct {
				Email    string `json:"email"`
				Password string `json:"password"`
				UserType string `json:"type"`
			}{
				Email:    "customer1@example.com",
				Password: "password2",
				UserType: services.TypeCustomer,
			}

			res, err := ts.Post("/login", responseBody)
			assert.Nil(t, err)
			assert.Equal(t, res.StatusCode, http.StatusForbidden)

			authCookie := ts.GetCookie(res, AuthCookie)
			assert.True(t, authCookie == nil)
		})

	})
}

func customerLoginTestSetup(ts *TestClient) {
	bodies := []struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		UserType string `json:"type"`
	}{
		{
			FullName: "customer1",
			Email:    "customer1@example.com",
			Password: "password1",
			UserType: services.TypeCustomer,
		},
		{
			FullName: "customer2",
			Email:    "customer2@example.com",
			Password: "password2",
			UserType: services.TypeCustomer,
		},
	}

	for _, body := range bodies {
		ts.Post("/signup", body)
	}
}
