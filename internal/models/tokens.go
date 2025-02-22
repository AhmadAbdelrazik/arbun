package models

import (
	"AhmadAbdelrazik/arbun/internal/domain/token"
	"bytes"
	"crypto/sha256"
	"errors"
	"time"
)

const (
	ScopeAuth = "authorization"
)

var (
	ErrTokenNotFound = errors.New("token not found")
)

type TokenModel struct {
	tokens []token.Token
}

func NewTokenModel() *TokenModel {
	return &TokenModel{
		tokens: make([]token.Token, 0),
	}
}

func (m *TokenModel) InsertToken(token token.Token) error {
	token.Plaintext = ""
	m.tokens = append(m.tokens, token)
	return nil
}

func (m *TokenModel) GetToken(plaintext string, scope string) (token.Token, error) {
	hash := sha256.Sum256([]byte(plaintext))
	for _, token := range m.tokens {
		if bytes.Equal(hash[:], token.Hash) && token.Scope == scope && token.ExpiryTime.After(time.Now()) {
			return token, nil
		}
	}

	return token.Token{}, ErrTokenNotFound
}

func (m *TokenModel) DeleteTokensByID(id int64) error {
	tokens := make([]token.Token, 0, len(m.tokens))

	for _, tt := range m.tokens {
		if tt.UserID == id {
			continue
		}
		tokens = append(tokens, tt)
	}

	return nil
}
