package server

import (
	"context"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
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

func TestGeneratePassword(t *testing.T) {
	pwd, err := generatePassword(13)
	assert.NoError(t, err)
	assert.EqualValues(t, 13, utf8.RuneCountInString(pwd))
	err = validatePassword(pwd)
	assert.NoError(t, err)

	pwd, err = generatePassword(16)
	assert.NoError(t, err)
	assert.EqualValues(t, 16, utf8.RuneCountInString(pwd))
	err = validatePassword(pwd)
	assert.NoError(t, err)
}

func TestGetUsername(t *testing.T) {
	username := "test_username"
	secret := "secret"
	token, err := generateToken(username, secret)
	assert.NoError(t, err)

	// ok
	md := metadata.New(map[string]string{TokenTitle: token})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	haveUsername, err := usernameFromToken(ctx, secret)
	assert.NoError(t, err)
	assert.EqualValues(t, username, haveUsername)

	// не найден контекст
	_, err = usernameFromToken(context.Background(), secret)
	assert.Error(t, err)

	// не найден токен в контексте
	_, err = usernameFromToken(metadata.NewIncomingContext(context.Background(), metadata.MD{}), secret)
	assert.Error(t, err)

	// неверный секрет
	_, err = usernameFromToken(ctx, secret+"!")
	assert.Error(t, err)
}
