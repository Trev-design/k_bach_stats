package sessioncrypto

import (
	"time"
)

// encrypt the payload.
// on success you'll get a string that contains the encrypted payload.
// on failure you'll get an error.
func (crypto *Crypt) EncryptPayload(payload []byte, timestamp time.Time) (string, error) {
	key, err := crypto.keys.GetKey(timestamp)
	if err != nil {
		return "", err
	}

	gcm, err := getGCM(key)
	if err != nil {
		return "", err
	}

	return newCipher(payload, gcm)
}

// decrypts encrypted payload.
// on success you'll get the decrypted payload as a string.
// on failure you'll get an error.
func (crypto *Crypt) DecryptPayload(payload string, timestamp time.Time) (string, error) {
	key, err := crypto.keys.GetKey(timestamp)
	if err != nil {
		return "", err
	}

	gcm, err := getGCM(key)
	if err != nil {
		return "", err
	}

	return plainFromCipher(payload, gcm)
}
