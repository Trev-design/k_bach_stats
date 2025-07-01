package sessionenclave_test

import (
	"auth_server/cmd/api/temporary/session/sessionenclave"
	"bytes"
	"testing"
)

func Test_SessionEnclave(t *testing.T) {
	sessionKeyInstance(t)
}

func sessionKeyInstance(t *testing.T) {
	t.Run("session_enclave_instance", func(t *testing.T) {
		key := sessionKeyInit(t)
		keyBytes := keyGet(t, key)
		keyChange(t, key, keyBytes)
		keyDestroy(t, key)
	})
}

func sessionKeyInit(t *testing.T) *sessionenclave.Key {
	var key *sessionenclave.Key

	t.Run("init_sessionenclave", func(t *testing.T) {
		enclave := sessionenclave.NewKey(1)
		key = enclave
	})

	return key
}

func keyGet(t *testing.T, key *sessionenclave.Key) []byte {
	var keyData []byte

	t.Run("get_key_from_enclave", func(t *testing.T) {
		keyBytes, err := key.GetKey(0)
		if err != nil {
			t.Fatal(err)
		}

		if keyBytes == nil {
			t.Fatal("invalid key")
		}

		keyData := make([]byte, len(keyBytes))
		copy(keyData, keyBytes)
	})

	return keyData
}

func keyChange(t *testing.T, key *sessionenclave.Key, keyData []byte) {
	t.Run("change_key_from_enclave", func(t *testing.T) {
		if err := key.ChangeKey(); err != nil {
			t.Fatal(err)
		}

		newKeyData, err := key.GetKey(0)
		if err != nil {
			t.Fatal(err)
		}

		if newKeyData == nil {
			t.Fatal("invalid key data")
		}

		if bytes.Equal(keyData, newKeyData) {
			t.Fatal("key should be different")
		}
	})
}

func keyDestroy(t *testing.T, key *sessionenclave.Key) {
	t.Run("destroy_key_from_enclave", func(t *testing.T) {
		if err := key.DestroyKey(); err != nil {
			t.Fatal(err)
		}
	})
}
