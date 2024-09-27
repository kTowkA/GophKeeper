package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrypt(t *testing.T) {
	original := "original text test 123 456 789 original text test 123 456 789 original text test 123 456 789"
	password := "veryStrong01P@ssw0rd!"
	sc := SimpleCrypter{}
	ciphertext, err := sc.Encrypt([]byte(original), password)
	assert.NoError(t, err)
	haveText, err := sc.Decrypt(ciphertext, password)
	assert.NoError(t, err)
	assert.EqualValues(t, original, string(haveText))
}
