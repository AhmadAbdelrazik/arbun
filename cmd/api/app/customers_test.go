package app

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"AhmadAbdelrazik/arbun/internal/pkg/assert"
	"fmt"
	"net/http"
	"testing"
)

func TestCustomer(t *testing.T) {
	ts := NewTestClient()
	defer ts.Close()

	t.Run("signup", func(t *testing.T) {
		customerSignup(t, ts)
	})
	t.Run("login", func(t *testing.T) {
		customerLogin(t, ts)
	})
}

func customerSignup(t *testing.T, ts *TestClient) {
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
				UserType: domain.TypeCustomer,
			},
			{
				FullName: "customer2",
				Email:    "customer2@example.com",
				Password: "password2",
				UserType: domain.TypeCustomer,
			},
		}

		for i, tt := range tests {
			t.Run(fmt.Sprintf("customer%d", i+1), func(t *testing.T) {
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
				UserType: domain.TypeCustomer,
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
				UserType: domain.TypeCustomer,
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
				UserType: domain.TypeCustomer,
			}

			res, err := ts.Post("/signup", requestBody)
			assert.Nil(t, err)
			assert.Equal(t, res.StatusCode, http.StatusBadRequest)

			authCookie := ts.GetCookie(res, AuthCookie)
			assert.True(t, authCookie == nil)
		})
	})

}

func customerLogin(t *testing.T, ts *TestClient) {
	t.Run("valid logins", func(t *testing.T) {
		tests := []struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			UserType string `json:"type"`
		}{
			{
				Email:    "customer1@example.com",
				Password: "password1",
				UserType: domain.TypeCustomer,
			},
			{
				Email:    "customer2@example.com",
				Password: "password2",
				UserType: domain.TypeCustomer,
			},
		}

		for i, tt := range tests {
			t.Run(fmt.Sprintf("customer%d", i+1), func(t *testing.T) {
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
				UserType: domain.TypeCustomer,
			}

			res, err := ts.Post("/login", responseBody)
			assert.Nil(t, err)
			assert.Equal(t, res.StatusCode, http.StatusUnauthorized)

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
				UserType: domain.TypeCustomer,
			}

			res, err := ts.Post("/login", responseBody)
			assert.Nil(t, err)
			assert.Equal(t, res.StatusCode, http.StatusUnauthorized)

			authCookie := ts.GetCookie(res, AuthCookie)
			assert.True(t, authCookie == nil)
		})

	})
}
