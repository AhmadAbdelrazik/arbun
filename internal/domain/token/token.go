package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

type Token struct {
	UserID     int64
	Plaintext  string
	Hash       []byte
	ExpiryTime time.Time
	Scope      string
}

func Generate(adminID int64, tokenScope string, ttl time.Duration) (Token, error) {
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
