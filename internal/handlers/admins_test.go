package handlers

import (
	"AhmadAbdelrazik/arbun/internal/assert"
	"AhmadAbdelrazik/arbun/internal/services"
	"fmt"
	"net/http"
	"testing"
)

func TestAdminSignup(t *testing.T) {
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
				FullName: "admin1",
				Email:    "admin1@example.com",
				Password: "password1",
				UserType: services.TypeAdmin,
			},
			{
				FullName: "admin2",
				Email:    "admin2@example.com",
				Password: "password2",
				UserType: services.TypeAdmin,
			},
		}

		for i, tt := range tests {
			t.Run(fmt.Sprintf("admin%d", i), func(t *testing.T) {
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
		t.Run("duplicate email", func(t *testing.T) {
			responseBody := struct {
				FullName string `json:"full_name"`
				Email    string `json:"email"`
				Password string `json:"password"`
				UserType string `json:"type"`
			}{
				FullName: "admin1",
				Email:    "admin1@example.com",
				Password: "password1",
				UserType: services.TypeAdmin,
			}

			res, err := ts.Post("/signup", responseBody)
			assert.Nil(t, err)
			assert.Equal(t, res.StatusCode, http.StatusBadRequest)

			authCookie := ts.GetCookie(res, AuthCookie)
			assert.True(t, authCookie == nil)
		})
	})
}

func TestAdminLogin(t *testing.T) {
	ts := NewTestClient()
	defer ts.Close()

	loginTestSetup(ts)

	t.Run("valid logins", func(t *testing.T) {
		tests := []struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			UserType string `json:"type"`
		}{
			{
				Email:    "admin1@example.com",
				Password: "password1",
				UserType: services.TypeAdmin,
			},
			{
				Email:    "admin2@example.com",
				Password: "password2",
				UserType: services.TypeAdmin,
			},
		}

		for i, tt := range tests {
			t.Run(fmt.Sprintf("admin%d", i), func(t *testing.T) {
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
				Email:    "admin3@example.com",
				Password: "password1",
				UserType: services.TypeAdmin,
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
				Email:    "admin1@example.com",
				Password: "password2",
				UserType: services.TypeAdmin,
			}

			res, err := ts.Post("/login", responseBody)
			assert.Nil(t, err)
			assert.Equal(t, res.StatusCode, http.StatusForbidden)

			authCookie := ts.GetCookie(res, AuthCookie)
			assert.True(t, authCookie == nil)
		})

	})
}

func loginTestSetup(ts *TestClient) {
	bodies := []struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		UserType string `json:"type"`
	}{
		{
			FullName: "admin1",
			Email:    "admin1@example.com",
			Password: "password1",
			UserType: services.TypeAdmin,
		},
		{
			FullName: "admin2",
			Email:    "admin2@example.com",
			Password: "password2",
			UserType: services.TypeAdmin,
		},
	}

	for _, body := range bodies {
		ts.Post("/signup", body)
	}
}
