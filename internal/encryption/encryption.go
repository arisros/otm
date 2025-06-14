package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
)

// GenerateRandomKey generates a secure 32-byte key
func GenerateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	return key, err
}

// GenerateNonce generates a 12-byte nonce for AES-GCM
func GenerateNonce() ([]byte, error) {
	nonce := make([]byte, 12)
	_, err := rand.Read(nonce)
	return nonce, err
}

// Encrypt encrypts plaintext using AES-GCM
func Encrypt(plaintext, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Seal(nil, nonce, plaintext, nil), nil
}

// Decrypt decrypts ciphertext using AES-GCM
func Decrypt(ciphertext, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Open(nil, nonce, ciphertext, nil)
}

// EncryptMessageLayer performs dual encryption (msg + key)
func EncryptMessageLayer(plaintext []byte) (cipherText, nonceMsg, encryptedKey, nonceKey []byte, err error) {
	// Generate a random per-message key
	msgKey, err := GenerateRandomKey()
	if err != nil {
		return
	}

	// Encrypt message with msgKey
	nonceMsg, err = GenerateNonce()
	if err != nil {
		return
	}

	cipherText, err = Encrypt(plaintext, msgKey, nonceMsg)
	if err != nil {
		return
	}

	// Encrypt msgKey using SECRET_KEY from env
	masterKeyB64 := os.Getenv("SECRET_KEY")
	masterKey, err := base64.StdEncoding.DecodeString(masterKeyB64)
	if err != nil || len(masterKey) != 32 {
		return nil, nil, nil, nil, errors.New("invalid master SECRET_KEY")
	}

	nonceKey, err = GenerateNonce()
	if err != nil {
		return
	}

	encryptedKey, err = Encrypt(msgKey, masterKey, nonceKey)
	return
}

// DecryptMessageLayer decrypts the stored message using encrypted msgKey
func DecryptMessageLayer(cipherText, nonceMsg, encryptedKey, nonceKey []byte) ([]byte, error) {
	masterKeyB64 := os.Getenv("SECRET_KEY")
	masterKey, err := base64.StdEncoding.DecodeString(masterKeyB64)
	if err != nil || len(masterKey) != 32 {
		return nil, errors.New("invalid master SECRET_KEY")
	}

	msgKey, err := Decrypt(encryptedKey, masterKey, nonceKey)
	if err != nil {
		return nil, err
	}

	return Decrypt(cipherText, msgKey, nonceMsg)
}
