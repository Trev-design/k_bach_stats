package sessionkeys_test

import (
	"auth_server/cmd/api/temporary/session/sessionkeys"
	"bytes"
	"log"
	"testing"
	"time"
)

var keys *sessionkeys.KeyManager

func TestMain(m *testing.M) {
	newKeys := sessionkeys.NewKeyManager(12 * time.Second)
	keys = newKeys
	go keys.ComputeRotateInterval()

	m.Run()

	err := keys.StopKeyManager()
	if err != nil {
		log.Println(err)
	}
}

func Test_InitAndClose(t *testing.T) {
	keymanager := sessionkeys.NewKeyManager(1 * time.Second)
	if err := keymanager.StopKeyManager(); err != nil {
		t.Fatal(err)
	}
}

func Test_InitRotateAndClose(t *testing.T) {
	keymanager := sessionkeys.NewKeyManager(1 * time.Second)
	go keymanager.ComputeRotateInterval()
	timestamp := time.Now().UTC()
	key, err := keymanager.GetKey(timestamp)
	if err != nil {
		t.Fatal(err)
	}

	if key == nil {
		t.Fatal("should get a valid key but got nil")
	}

	time.Sleep(1125 * time.Millisecond)
	newTimestamp := time.Now().UTC()
	newKey, err := keymanager.GetKey(newTimestamp)
	if err != nil {
		t.Fatal(err)
	}

	if newKey == nil {
		t.Fatal("should get a valid key but got nil")
	}

	if bytes.Equal(key, newKey) {
		t.Fatal("keys should be different")
	}

	if err = keymanager.StopKeyManager(); err != nil {
		t.Fatal(err)
	}
}

func Test_GetKey(t *testing.T) {
	timestamp := time.Now().UTC()
	key, err := keys.GetKey(timestamp)
	if err != nil {
		t.Fatal(err)
	}

	if key == nil {
		t.Fatal("should get a valid key but got nil")
	}
}

func Test_GetAnotherKey(t *testing.T) {
	timestamp := time.Now().UTC()
	key, err := keys.GetKey(timestamp)
	if err != nil {
		t.Fatal(err)
	}

	if key == nil {
		t.Fatal("should get a valid key but got nil")
	}

	newTimeStamp := time.Now().UTC().Add(1 * time.Second)
	newKey, err := keys.GetKey(newTimeStamp)
	if err != nil {
		t.Fatal(err)
	}

	if newKey == nil {
		t.Fatal("should get a valid key but got nil")
	}

	if bytes.Equal(key, newKey) {
		t.Fatal("keys should be different")
	}
}

func Test_GetOldKey(t *testing.T) {
	time.Sleep(7 * time.Second)
	timestamp := time.Now().UTC()
	key, err := keys.GetKey(timestamp)
	if err != nil {
		t.Fatal(err)
	}
	if key == nil {
		t.Fatal("should get a valid key but got nil")
	}

	time.Sleep(7 * time.Second)
	newKey, err := keys.GetKey(timestamp)
	if err != nil {
		t.Fatal(err)
	}

	if newKey == nil {
		t.Fatal("should get a valid key but got nil")
	}

	if !bytes.Equal(key, newKey) {
		t.Fatal("keys should be equal")
	}
}

func Test_GetKeyFailedExpiredKey(t *testing.T) {
	timestamp := time.Now().UTC().Add(-13 * time.Second)
	if _, err := keys.GetKey(timestamp); err == nil {
		t.Fatal("should fail but got succeed")
	}
}
