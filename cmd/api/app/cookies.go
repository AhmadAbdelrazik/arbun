package app

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"net/http"
)

const (
	AuthCookie = "AuthCookie"
)

func SetAuthTokenCookie(w http.ResponseWriter, token domain.Token) {
	cookie := &http.Cookie{
		Name:    AuthCookie,
		Value:   token.Plaintext,
		Expires: token.ExpiryTime,
	}

	http.SetCookie(w, cookie)
}

func GetAuthToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(AuthCookie)
	if err != nil {
		return "", err
	}

	err = cookie.Valid()
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}
