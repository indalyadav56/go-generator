package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
)

func GenerateRandomPassphrase() string {
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_-+=[]{};:,.<>?/~"

	passphraseLength := 32

	passphrase := make([]byte, passphraseLength)
	charSetLength := big.NewInt(int64(len(charSet)))

	for i := range passphrase {
		randomIndex, err := rand.Int(rand.Reader, charSetLength)
		if err != nil {
			fmt.Println("Error generating random index:", err)
			return ""
		}
		passphrase[i] = charSet[randomIndex.Int64()]
	}

	return string(passphrase)
}

func Encrypt(input []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, input, nil)

	return hex.EncodeToString(nonce) + hex.EncodeToString(ciphertext), nil
}

func Decrypt(cipherText string, key []byte) (string, error) {
	decoded, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	if len(decoded) < 24 {
		return "", errors.New("invalid cipher text")
	}

	nonce := decoded[:12]
	cipherTextBytes := decoded[12:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, cipherTextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
