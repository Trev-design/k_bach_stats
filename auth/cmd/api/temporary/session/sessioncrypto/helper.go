package sessioncrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

func getGCM(key []byte) (cipher.AEAD, error) {
	aes, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewGCM(aes)
}

func newCipher(payload []byte, gcm cipher.AEAD) (string, error) {
	nonce, err := newNonce(gcm)
	if err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, payload, nil)
	base64CipherText := base64.RawURLEncoding.EncodeToString(cipherText)

	return base64CipherText, nil
}

func plainFromCipher(payload string, gcm cipher.AEAD) (string, error) {
	nonce, cipher, err := getDecryptCredentials(gcm, payload)
	if err != nil {
		return "", err
	}

	return getPlainFromCipher(gcm, nonce, cipher)
}

func newNonce(gcm cipher.AEAD) ([]byte, error) {
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	return nonce, nil
}

func getDecryptCredentials(gcm cipher.AEAD, cipher string) ([]byte, []byte, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(cipher)
	if err != nil {
		return nil, nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce := decoded[:nonceSize]
	decodedCipher := decoded[nonceSize:]

	return nonce, decodedCipher, nil
}

func getPlainFromCipher(gcm cipher.AEAD, nonce, cipher []byte) (string, error) {
	plain, err := gcm.Open(nil, nonce, cipher, nil)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}
