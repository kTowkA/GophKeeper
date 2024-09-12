package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantError error
	}{
		{
			name:      "маленькая длина",
			value:     "123",
			wantError: ErrValidatePasswordMinLength,
		},
		{
			name:      "нет символов в нижнем регистре",
			value:     "123",
			wantError: ErrValidatePasswordLowerCase,
		},
		{
			name:      "нет символов в верхнем регистре",
			value:     "123",
			wantError: ErrValidatePasswordUpperCase,
		},
		{
			name:      "нет чисел",
			value:     "qwe",
			wantError: ErrValidatePasswordNumbers,
		},
		{
			name:      "нет спецсимволов",
			value:     "qwe",
			wantError: ErrValidatePasswordSymbols,
		},
		{
			name:      "все хорошо",
			value:     "123qweQWEQWE!",
			wantError: nil,
		},
		{
			name:      "все хорошо",
			value:     "123qweQWEQWE.",
			wantError: nil,
		},
		{
			name:      "все хорошо",
			value:     "123qweQWEQWE^",
			wantError: nil,
		},
		{
			name:      "все хорошо",
			value:     "123qweQWEQWE(",
			wantError: nil,
		},
		{
			name:      "все хорошо",
			value:     "123qweQWEQWE)",
			wantError: nil,
		},
		{
			name:      "все хорошо",
			value:     "123qweQWEQWE[",
			wantError: nil,
		},
		{
			name:      "все хорошо]",
			value:     "123qweQWEQWE]",
			wantError: nil,
		},
	}
	for _, tt := range tests {
		err := validatePassword(tt.value)
		assert.ErrorIs(t, err, tt.wantError, tt.name)
	}
}

func TestGenerateToken(t *testing.T) {
	login := "login"
	secret := "secret"
	_, err := generateToken(login, secret)
	assert.NoError(t, err)
}
