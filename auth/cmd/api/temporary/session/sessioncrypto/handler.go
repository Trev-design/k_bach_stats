package sessioncrypto

import (
	"time"
)

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
