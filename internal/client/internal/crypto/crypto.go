package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

type Crypto struct {
	gcm cipher.AEAD
}

func NewCrypto(pass string) (*Crypto, error) {
	key := HashPassword(pass)
	c, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("error on creating new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("error on creating new gcm: %w", err)
	}

	return &Crypto{
		gcm: gcm,
	}, nil
}

func HashPassword(pass string) [32]byte {
	return sha256.Sum256([]byte(pass))
}

func (c *Crypto) EncryptData(data string) ([]byte, error) {
	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("error reading rand nonce: %w", err)
	}

	return c.gcm.Seal(nonce, nonce, []byte(data), nil), nil
}

func (c *Crypto) DecryptData(data []byte) (string, error) {
	if data == nil {
		return "", nil
	}

	nonceSize := c.gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("wrong size of crypted data")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := c.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("error on decrypting data: %w", err)
	}

	return string(plaintext), nil
}
