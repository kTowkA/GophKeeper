package crypto

// Crypter простой интерфейс для шифрования/дешифрования сообщений
type Crypter interface {
	// Encrypt шифрование сообщения value парольной фразой password
	Encrypt(value []byte, password string) ([]byte, error)
	// Decrypt дешифрование сообщения value парольной фразой password
	Decrypt(value []byte, password string) ([]byte, error)
}
