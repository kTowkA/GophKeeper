// реализация интерфейса Crypter
package simplecrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
)

var SimpleCrypter = simpleCrypter{}

type simpleCrypter struct{}

func (sc simpleCrypter) Encrypt(value []byte, password string) ([]byte, error) {
	key := sha256.Sum256([]byte(password))

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}
	nonce := key[len(key)-aesgcm.NonceSize():]

	ciphertext := aesgcm.Seal(nil, nonce, value, nil)
	return ciphertext, nil
}

func (sc simpleCrypter) Decrypt(value []byte, password string) ([]byte, error) {
	key := sha256.Sum256([]byte(password))

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}
	nonce := key[len(key)-aesgcm.NonceSize():]

	// расшифровываем
	decrypted, err := aesgcm.Open(nil, nonce, value, nil)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}
