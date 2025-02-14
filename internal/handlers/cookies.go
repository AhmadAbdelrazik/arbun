package handlers

import (
	"AhmadAbdelrazik/arbun/internal/services"
	"net/http"
)

const (
	AuthCookie = "AuthCookie"
)

func SetAuthTokenCookie(w http.ResponseWriter, token services.Token) {
	cookie := &http.Cookie{
		Name:    AuthCookie,
		Value:   token.Plaintext,
		Expires: token.ExpiryTime,
	}

	http.SetCookie(w, cookie)
}

func GetAuthToken(r *http.Request) (services.Token, error) {
	cookie, err := r.Cookie(AuthCookie)
	if err != nil {
		return services.Token{}, err
	}

	// TODO: validate token

	token := services.Token{
		Plaintext: cookie.Value,
	}

	return token, nil
}
