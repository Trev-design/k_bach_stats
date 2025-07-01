package sessionkeys_test

import (
	"auth_server/cmd/api/temporary/session/sessionkeys"
	"bytes"
	"testing"
	"time"
)

func Test_SessionKeys(t *testing.T) {
	sessionKeys(t)
}

func sessionKeys(t *testing.T) {
	t.Run("sessio_keys_instance", func(t *testing.T) {
		manager := sessionKeysInit(t)
		// hier soll ein background service eingebunden werden
		go manager.ComputeRotateInterval()

		key, timeStramp := getSessionKey(t, manager)
		getAnotherSessionKey(t, manager, key)
		oldKey(t, manager, timeStramp, key)
		expiredKey(t, manager, timeStramp)
		closeKeyManager(t, manager)
	})
}

func sessionKeysInit(t *testing.T) *sessionkeys.KeyManager {
	var manager *sessionkeys.KeyManager

	t.Run("init_session_keys", func(t *testing.T) {
		keyManager := sessionkeys.NewKeyManager(1 * time.Second)
		manager = keyManager
	})

	return manager
}

func oldKey(
	t *testing.T,
	manager *sessionkeys.KeyManager,
	timeStamp time.Time,
	key []byte) {
	t.Run("old_key", func(t *testing.T) {
		time.Sleep(700 * time.Millisecond)
		otherKey, err := manager.GetKey(timeStamp)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(otherKey, key) {
			t.Fatal("should be the same key")
		}
	})
}

func expiredKey(
	t *testing.T,
	manager *sessionkeys.KeyManager,
	timeStamp time.Time,
) {
	t.Run("expired_key", func(t *testing.T) {
		time.Sleep(1 * time.Second)

		_, err := manager.GetKey(timeStamp)
		if err == nil {
			t.Fatal("should fail because key is expired but got succeed")
		}

		t.Logf("the error is %s", err.Error())
	})
}

func getSessionKey(t *testing.T, manager *sessionkeys.KeyManager) ([]byte, time.Time) {
	var keyBytes []byte
	var userTimeStamp time.Time

	t.Run("get_session_key_from_keys", func(t *testing.T) {
		timeStamp := time.Now().UTC()
		key, err := manager.GetKey(timeStamp)
		if err != nil {
			t.Fatal(err)
		}

		userTimeStamp = timeStamp

		keyBytes = make([]byte, len(key))
		copy(keyBytes, key)
	})

	return keyBytes, userTimeStamp
}

func getAnotherSessionKey(t *testing.T, manager *sessionkeys.KeyManager, key []byte) {
	t.Run("get_another_session_key", func(t *testing.T) {
		anotherKey, err := manager.GetKey(time.Now().Add(1 * time.Second).UTC())
		if err != nil {
			t.Fatal(err)
		}

		if bytes.Equal(key, anotherKey) {
			t.Fatal(err)
		}
	})
}

func closeKeyManager(t *testing.T, manager *sessionkeys.KeyManager) {
	t.Run("close_store", func(t *testing.T) {
		if err := manager.StopKeyManager(); err != nil {
			t.Fatal(err)
		}
	})
}
