package models

import (
	"AhmadAbdelrazik/arbun/internal/assert"
	"AhmadAbdelrazik/arbun/internal/domain/token"
	"testing"
	"time"
)

func TestTokens(t *testing.T) {
	tokenModel := newTokenModel()

	token, err := token.Generate(1, ScopeAuth, 3*time.Hour)
	assert.Nil(t, err)

	err = tokenModel.InsertToken(token)
	assert.Nil(t, err)

	tt, err := tokenModel.GetToken(token.Plaintext, ScopeAuth)
	assert.Nil(t, err)
	assert.Equal(t, tt.UserID, token.UserID)
}
