package models

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"time"
)

const (
	ScopeAuth = "authorization"
)

var (
	ErrTokenNotFound = errors.New("token not found")
)

type Token struct {
	UserID     int64
	Plaintext  string
	Hash       []byte
	ExpiryTime time.Time
	Scope      string
}

func GenerateToken(adminID int64, tokenScope string, ttl time.Duration) (Token, error) {
	t := Token{
		Scope:      tokenScope,
		UserID:     adminID,
		ExpiryTime: time.Now().Add(ttl),
	}

	// allocates 16 bytes from memory
	randomBytes := make([]byte, 16)

	// fill it with cryptographically random numbers
	_, err := rand.Read(randomBytes)
	if err != nil {
		return Token{}, err
	}

	// encode the random bytes to be readable
	t.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	// hash the token
	hash := sha256.Sum256([]byte(t.Plaintext))
	t.Hash = hash[:]

	return t, nil
}

type TokenModel struct {
	tokens []Token
}

func NewTokenModel() *TokenModel {
	return &TokenModel{
		tokens: make([]Token, 0),
	}
}

func (m *TokenModel) InsertToken(token Token) error {
	token.Plaintext = ""
	m.tokens = append(m.tokens, token)
	return nil
}

func (m *TokenModel) GetToken(plaintext string, scope string) (Token, error) {
	hash := sha256.Sum256([]byte(plaintext))
	for _, token := range m.tokens {
		if bytes.Equal(hash[:], token.Hash) && token.Scope == scope && token.ExpiryTime.After(time.Now()) {
			return token, nil
		}
	}

	return Token{}, ErrTokenNotFound
}

func (m *TokenModel) DeleteTokensByID(id int64) error {
	tokens := make([]Token, 0, len(m.tokens))

	for _, tt := range m.tokens {
		if tt.UserID == id {
			continue
		}
		tokens = append(tokens, tt)
	}

	return nil
}
