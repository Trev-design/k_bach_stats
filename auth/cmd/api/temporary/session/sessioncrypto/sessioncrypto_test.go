package sessioncrypto_test

import (
	"auth_server/cmd/api/temporary/session/sessioncrypto"
	"crypto/rand"
	"encoding/base64"
	random "math/rand"
	"testing"
	"time"
)

type cryptoCreds struct {
	payload   string
	timeStamp time.Time
}

func Test_SessionCrypto(t *testing.T) {
	sessionCryptoInstance(t)
}

func sessionCryptoInstance(t *testing.T) {
	t.Run("sessioncrypto_instance", func(t *testing.T) {
		crypto := initializeCryptoInstance(t)
		go crypto.ComputeRotateInterval()
		creds := getEncrypted(t, crypto)
		getDecryptedSuccess(t, crypto, creds)
		getDecryptedFailure(t, crypto, creds)
		tryGetWithExpiredKey(t, crypto, creds)
		tryGetWithOldKey(t, crypto)
		closeSessionCrypto(t, crypto)
	})
}

func tryGetWithOldKey(t *testing.T, crypto *sessioncrypto.Crypt) {
	t.Run("try_get_with_old_key", func(t *testing.T) {
		time.Sleep(700 * time.Millisecond)
		currentTime := time.Now().UTC()
		creds, err := crypto.EncryptPayload([]byte("wuzzup"), currentTime)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(600 * time.Millisecond)
		_, err = crypto.DecryptPayload(creds, currentTime)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func tryGetWithExpiredKey(t *testing.T, crypto *sessioncrypto.Crypt, creds *cryptoCreds) {
	t.Run("try_get_with_expired_key", func(t *testing.T) {
		time.Sleep(1100 * time.Millisecond)
		_, err := crypto.DecryptPayload(creds.payload, creds.timeStamp)
		if err == nil {
			t.Fatal("should fail because expired key but got succeed")
		}

		t.Logf("the error is %s", err.Error())
	})
}

func initializeCryptoInstance(t *testing.T) *sessioncrypto.Crypt {
	var sessionCrypto *sessioncrypto.Crypt

	t.Run("initialize_session_crypto_instance", func(t *testing.T) {
		crypto, err := sessioncrypto.NewCrypt(1 * time.Second)
		if err != nil {
			t.Fatal(err)
		}

		sessionCrypto = crypto
	})

	return sessionCrypto
}

func getEncrypted(t *testing.T, crypto *sessioncrypto.Crypt) *cryptoCreds {
	var encryptedPayload string
	var currentTimeStamp time.Time

	t.Run("get_encrypted_from_sessioncrypto", func(t *testing.T) {
		timeStamp := time.Now().UTC()
		payload, err := crypto.EncryptPayload([]byte("halli hallo halloechen"), timeStamp)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(payload)

		currentTimeStamp = timeStamp
		encryptedPayload = payload
	})

	return &cryptoCreds{
		timeStamp: currentTimeStamp,
		payload:   encryptedPayload,
	}
}

func getDecryptedSuccess(t *testing.T, crypto *sessioncrypto.Crypt, creds *cryptoCreds) {
	t.Run("get_decrypted_from_sessioncrypto_success", func(t *testing.T) {
		payload, err := crypto.DecryptPayload(creds.payload, creds.timeStamp)
		if err != nil {
			t.Fatal(err)
		}

		if payload != "halli hallo halloechen" {
			t.Fatal("not expected payload: payload should be halli hallo halloechen")
		}
	})
}

func getDecryptedFailure(t *testing.T, crypto *sessioncrypto.Crypt, creds *cryptoCreds) {
	t.Run("get_decrypted_from_sessioncrypto_failure", func(t *testing.T) {
		falseCipher(t, crypto, creds)
		falseNonce(t, crypto, creds)
		falseKey(t, crypto, creds)
		falsePayload(t, crypto, creds)
	})
}

func falseCipher(t *testing.T, crypto *sessioncrypto.Crypt, creds *cryptoCreds) {
	t.Run("false_cipher_payload", func(t *testing.T) {
		completeFalseCipher(t, crypto, creds)
		falseRandomByteCipher(t, crypto, creds)
	})
}

func falseRandomByteCipher(t *testing.T, crypto *sessioncrypto.Crypt, creds *cryptoCreds) {
	t.Run("false_random_byte", func(t *testing.T) {
		encodedNewSecret, err := getCoruptedRandomByteCipher(creds.payload, 12)
		if err != nil {
			t.Fatal(err)
		}

		_, err = crypto.DecryptPayload(encodedNewSecret, creds.timeStamp)
		if err == nil {
			t.Fatal("expected failure got succeed")
		}

		t.Logf("the reason of the failure for decryption is: %s", err.Error())
	})
}

func completeFalseCipher(t *testing.T, crypto *sessioncrypto.Crypt, creds *cryptoCreds) {
	t.Run("complete_false_cipher", func(t *testing.T) {
		encodedNewSecret, err := getCompleteCoruptedCipher(creds.payload, 12)
		if err != nil {
			t.Fatal(err)
		}

		_, err = crypto.DecryptPayload(encodedNewSecret, creds.timeStamp)
		if err == nil {
			t.Fatal("expected failure got succeed")
		}

		t.Logf("the reason of the failure for decryption is: %s", err.Error())
	})
}

func falseNonce(t *testing.T, crypto *sessioncrypto.Crypt, creds *cryptoCreds) {
	t.Run("false_nonce_payload", func(t *testing.T) {
		falseRandomByteNonce(t, crypto, creds)
		completelyFalseNonce(t, crypto, creds)
	})
}

func falseRandomByteNonce(t *testing.T, crypto *sessioncrypto.Crypt, creds *cryptoCreds) {
	t.Run("false_random_byte", func(t *testing.T) {
		encodedNewSecret, err := getCoruptedRandomByteNonce(creds.payload, 12)
		if err != nil {
			t.Fatal(err)
		}

		_, err = crypto.DecryptPayload(encodedNewSecret, creds.timeStamp)
		if err == nil {
			t.Fatal("expected failure got succeed")
		}

		t.Logf("the reason of the failure for decryption is: %s", err.Error())
	})
}

func completelyFalseNonce(t *testing.T, crypto *sessioncrypto.Crypt, creds *cryptoCreds) {
	encodedNewSecret, err := getCompleteCoruptedNonce(creds.payload, 12)
	if err != nil {
		t.Fatal(err)
	}

	_, err = crypto.DecryptPayload(encodedNewSecret, creds.timeStamp)
	if err == nil {
		t.Fatal("expected failure got succeed")
	}

	t.Logf("the reason of the failure for decryption is: %s", err.Error())
}

func falseKey(t *testing.T, crypto *sessioncrypto.Crypt, creds *cryptoCreds) {
	t.Run("false_key", func(t *testing.T) {
		_, err := crypto.DecryptPayload(creds.payload, time.Now().Add(1*time.Second).UTC())
		if err == nil {
			t.Fatal("expected failure got succeed")
		}

		t.Logf("the reason of the failure for decryption is: %s", err.Error())
	})
}

func falsePayload(t *testing.T, crypto *sessioncrypto.Crypt, creds *cryptoCreds) {
	t.Run("false_payload_in_general", func(t *testing.T) {
		encodedNewSecret, err := getCompleteCoruptedPayload(creds.payload)
		if err != nil {
			t.Fatal(err)
		}

		_, err = crypto.DecryptPayload(encodedNewSecret, creds.timeStamp)
		if err == nil {
			t.Fatal("expected failure got succeed")
		}

		t.Logf("the reason of the failure for decryption is: %s", err.Error())
	})
}

func closeSessionCrypto(t *testing.T, crypto *sessioncrypto.Crypt) {
	t.Run("close_sessioncrypto", func(t *testing.T) {
		if err := crypto.CloseCrypto(); err != nil {
			t.Fatal(err)
		}
	})
}

func getCoruptedRandomByteCipher(payload string, nonceSize int) (string, error) {
	secret, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return "", err
	}

	nonce := secret[:nonceSize]
	cipher := secret[nonceSize:]

	manipulatedCipher := cipher
	randomIndex := random.Intn(200) % len(manipulatedCipher)

	manipulatedCipher[randomIndex] ^= 0xff

	newSecret := append(nonce, manipulatedCipher...)

	return base64.RawURLEncoding.EncodeToString(newSecret), nil
}

func getCompleteCoruptedCipher(payload string, nonceSize int) (string, error) {
	secret, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return "", err
	}

	nonce := secret[:nonceSize]
	cipher := secret[nonceSize:]

	manipulatedCipher := make([]byte, len(cipher))
	if _, err := rand.Read(manipulatedCipher); err != nil {
		return "", err
	}

	newSecret := append(nonce, manipulatedCipher...)

	return base64.RawURLEncoding.EncodeToString(newSecret), nil
}

func getCoruptedRandomByteNonce(payload string, nonceSize int) (string, error) {
	secret, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return "", err
	}

	nonce := secret[:nonceSize]
	cipher := secret[nonceSize:]

	manipulatedNonce := nonce
	randomIndex := random.Intn(100) % len(manipulatedNonce)
	manipulatedNonce[randomIndex] ^= 0xff

	newSecret := append(manipulatedNonce, cipher...)

	return base64.RawURLEncoding.EncodeToString(newSecret), nil
}

func getCompleteCoruptedNonce(payload string, nonceSize int) (string, error) {
	secret, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return "", err
	}

	nonce := secret[:nonceSize]
	cipher := secret[nonceSize:]

	manipulatedNonce := make([]byte, len(nonce))
	if _, err := rand.Read(manipulatedNonce); err != nil {
		return "", err
	}

	newSecret := append(manipulatedNonce, cipher...)

	return base64.RawURLEncoding.EncodeToString(newSecret), nil
}

func getCompleteCoruptedPayload(payload string) (string, error) {
	secret, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return "", err
	}

	newSecret := make([]byte, len(secret))
	if _, err := rand.Read(newSecret); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(newSecret), nil
}
