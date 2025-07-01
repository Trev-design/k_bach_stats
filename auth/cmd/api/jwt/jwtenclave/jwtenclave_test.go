package jwtenclave_test

import (
	"auth_server/cmd/api/jwt/jwtenclave"
	"bytes"
	"os"
	"testing"
)

var enclave *jwtenclave.Enclave

func TestMain(m *testing.M) {
	newEnclave := jwtenclave.NewEnclave()
	enclave = newEnclave

	code := m.Run()

	enclave.DestroySeeds()

	os.Exit(code)
}

func TestInit(t *testing.T) {
	newEnclave := jwtenclave.NewEnclave()
	t.Log(newEnclave)

	if err := newEnclave.DestroySeeds(); err != nil {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	if _, err := enclave.GetSeed(); err != nil {
		t.Fatal(err)
	}
}

func TestGetFailedDestroyedBuffer(t *testing.T) {
	newEnclave := jwtenclave.NewEnclave()

	if err := newEnclave.DestroySeeds(); err != nil {
		t.Fatal(err)
	}

	_, err := newEnclave.GetSeed()
	if err == nil {
		t.Fatal("should fail because of destroyed buffer but got succeed")
	}

	t.Logf("got err: %s", err.Error())
}

func TestDestroy(t *testing.T) {
	newEnclave := jwtenclave.NewEnclave()
	if err := newEnclave.DestroySeeds(); err != nil {
		t.Fatal(err)
	}
}

func TestDestroyFailedDoubleDestroy(t *testing.T) {
	newEnclave := jwtenclave.NewEnclave()

	if err := newEnclave.DestroySeeds(); err != nil {
		t.Fatal(err)
	}

	err := newEnclave.DestroySeeds()
	if err == nil {
		t.Fatal("should fail because of already destroyed buffer but got succeed")
	}

	t.Logf("got err: %s", err.Error())
}

func TestRotation(t *testing.T) {
	oldSeed, err := enclave.GetSeed()
	if err != nil {
		t.Fatal(err)
	}

	if err = enclave.NewSeeds(); err != nil {
		t.Fatal(err)
	}

	newSeed, err := enclave.GetSeed()
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(oldSeed, newSeed) {
		t.Fatal("should be another key")
	}
}

func TestNewFailedDestroyedKeys(t *testing.T) {
	newEnclave := jwtenclave.NewEnclave()

	if err := newEnclave.DestroySeeds(); err != nil {
		t.Fatal(err)
	}

	err := newEnclave.NewSeeds()
	if err == nil {
		t.Fatal("should fail because of destroyed buffer but got succeed")
	}

	t.Logf("got err: %s", err.Error())
}
