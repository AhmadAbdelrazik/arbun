package models

import (
	"AhmadAbdelrazik/arbun/internal/assert"
	"testing"
	"time"
)

func TestTokens(t *testing.T) {
	tokenModel := NewTokenModel()

	token, err := GenerateToken(1, ScopeAuth, 3*time.Hour)
	assert.Nil(t, err)

	err = tokenModel.InsertToken(token)
	assert.Nil(t, err)

	tt, err := tokenModel.GetToken(token.Plaintext, ScopeAuth)
	assert.Nil(t, err)
	assert.Equal(t, tt.UserID, token.UserID)
}
