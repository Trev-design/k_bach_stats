package sessionenclave_test

import (
	"auth_server/cmd/api/temporary/session/sessionenclave"
	"bytes"
	"testing"
)

var enclave *sessionenclave.Key

func TestMain(m *testing.M) {
	newEnclave := sessionenclave.NewKey(4)
	enclave = newEnclave

	m.Run()
}

func Test_GetKey(t *testing.T) {
	_, err := enclave.GetKey(2)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_GetKeyfailedDestroyedKey(t *testing.T) {
	newEnclave := sessionenclave.NewKey(2)
	if err := newEnclave.DestroyKey(); err != nil {
		t.Fatal(err)
	}

	if _, err := newEnclave.GetKey(4); err == nil {
		t.Fatal("should fail but got succeed")
	}
}

func Test_ChangeKey(t *testing.T) {
	oldKey, err := enclave.GetKey(2)
	if err != nil {
		t.Fatal(err)
	}

	if err = enclave.ChangeKey(); err != nil {
		t.Fatal(err)
	}

	newKey, err := enclave.GetKey(2)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(oldKey, newKey) {
		t.Fatal("keys should be different")
	}
}

func Test_DestroyKeys(t *testing.T) {
	if err := enclave.DestroyKey(); err != nil {
		t.Fatal(err)
	}
}
