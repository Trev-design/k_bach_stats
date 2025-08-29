package sessioncrypto_test

import (
	"auth_server/cmd/api/temporary/session/sessioncrypto"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	random "math/rand"
	"testing"
	"time"
)

var crypto *sessioncrypto.Crypt

func TestMain(m *testing.M) {
	newCrypt, _ := sessioncrypto.NewCrypt(2 * time.Second)
	crypto = newCrypt

	m.Run()
}

func Test_InitCrypt(t *testing.T) {
	_, err := sessioncrypto.NewCrypt(2 * time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_InitAndCloseCrypt(t *testing.T) {
	newCrypt, err := sessioncrypto.NewCrypt(2 * time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	if err = newCrypt.CloseCrypto(); err != nil {
		t.Fatal(err)
	}
}

func Test_GetEncrypred(t *testing.T) {
	payload, err := crypto.EncryptPayload([]byte("halli hallo"), time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("the encrypted payload is %s", payload)
}

func Test_GetDecrypted(t *testing.T) {
	message := []byte("hallo halli")
	timestamp := time.Now().UTC()

	encrypted, err := crypto.EncryptPayload(message, timestamp)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("the encrypted data of %s is %s", string(message), encrypted)

	decryted, err := crypto.DecryptPayload(encrypted, timestamp)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte(decryted), message) {
		t.Fatalf("false decryption %s is not equal to %s", decryted, string(message))
	}
}

func Test_GetDecryptedWithOldKey(t *testing.T) {
	time.Sleep(1200 * time.Millisecond)
	timestamp := time.Now().UTC()
	encrypted, err := crypto.EncryptPayload([]byte("hdjhfjskashdsh"), timestamp)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1200 * time.Millisecond)
	_, err = crypto.DecryptPayload(encrypted, timestamp)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_GetdecryptedFailedWithExpiredKey(t *testing.T) {
	timestamp := time.Now().UTC()
	encrypted, err := crypto.EncryptPayload([]byte("jhdjhdsdjdsjhasghshg"), timestamp)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(2345 * time.Millisecond)
	_, err = crypto.DecryptPayload(encrypted, timestamp)
	if err == nil {
		t.Fatal("should fail but got succeed")
	}
	t.Log(err.Error())
}

func Test_GetDecryptedFailureFalseKey(t *testing.T) {
	message := []byte("hallo halli")
	timestamp := time.Now().UTC()

	encrypted, err := crypto.EncryptPayload(message, timestamp)
	if err != nil {
		t.Fatal(err)
	}

	_, err = crypto.DecryptPayload(encrypted, timestamp.Add(1*time.Second))
	if err == nil {
		t.Fatal("should fail but succeed")
	}
}

func Test_GetDecryptedFailureFalsePayloadRandomByte(T *testing.T) {

}

func Test_GetDecryptedFailureCompletelyFalsePayload(t *testing.T) {
	timestamp := time.Now().UTC()
	encrypted, err := crypto.EncryptPayload([]byte("Holla Halli"), timestamp)
	if err != nil {
		t.Fatal(err)
	}

	corupted, err := getCompleteCoruptedPayload(encrypted)
	if err != nil {
		t.Fatal(err)
	}

	_, err = crypto.DecryptPayload(corupted, timestamp)
	if err == nil {
		t.Fatal("should fail but got succeed")
	}
}

func Test_GetDecryptedFailureFalseCipherRandomByte(t *testing.T) {
	timestamp := time.Now().UTC()
	encrypted, err := crypto.EncryptPayload([]byte("hilla holli"), timestamp)
	if err != nil {
		t.Fatal(err)
	}

	corupted, err := getCoruptedRandomByteCipher(encrypted, 12)
	if err != nil {
		t.Fatal(err)
	}

	_, err = crypto.DecryptPayload(corupted, timestamp)
	if err == nil {
		t.Fatal("should fail but got succeed")
	}
}

func Test_GetDecryptedFailureCompletelyWrongCipher(t *testing.T) {
	timestamp := time.Now().UTC()
	encrypted, err := crypto.EncryptPayload([]byte("Hillo Halla"), timestamp)
	if err != nil {
		t.Fatal(err)
	}

	corupted, err := getCompleteCoruptedCipher(encrypted, 12)
	if err != nil {
		t.Fatal(err)
	}

	_, err = crypto.DecryptPayload(corupted, timestamp)
	if err == nil {
		t.Fatal("should fail but got succeed")
	}
}

func Test_GetDecryptedFailureFalseNonceRandomByte(t *testing.T) {
	timestamp := time.Now().UTC()
	encrypted, err := crypto.EncryptPayload([]byte("ciaoi"), timestamp)
	if err != nil {
		t.Fatal(err)
	}

	corupted, err := getCoruptedRandomByteNonce(encrypted, 12)
	if err != nil {
		t.Fatal(err)
	}

	_, err = crypto.DecryptPayload(corupted, timestamp)
	if err == nil {
		t.Fatal("should fail but got succeed")
	}
}

func Test_GetDecryptedFailureCompletelyfalseNonce(t *testing.T) {
	timestamp := time.Now().UTC()
	encrypted, err := crypto.EncryptPayload([]byte("djhgkdrhgj"), timestamp)
	if err != nil {
		t.Fatal(err)
	}

	corupted, err := getCompleteCoruptedNonce(encrypted, 12)
	if err != nil {
		t.Fatal(err)
	}

	_, err = crypto.DecryptPayload(corupted, timestamp)
	if err == nil {
		t.Fatal("should fail but got succeed")
	}
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
