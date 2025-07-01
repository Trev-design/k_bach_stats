package jwtmanager_test

import (
	"auth_server/cmd/api/jwt/jwtmanager"
	"bytes"
	"os"
	"testing"
	"time"
)

var manager *jwtmanager.SeedManager

func TestMain(m *testing.M) {
	newManager := jwtmanager.NewSeedManager(time.Second)
	manager = newManager
	go manager.ComputeInterval()

	code := m.Run()

	manager.CloseSeedManager()

	os.Exit(code)
}

func TestInit(t *testing.T) {
	newManager := jwtmanager.NewSeedManager(time.Second)

	if err := newManager.CloseSeedManager(); err != nil {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	timestamp := time.Now()
	_, err := manager.GetSeed(timestamp)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetOld(t *testing.T) {
	time.Sleep(700 * time.Millisecond)
	timestamp := time.Now()
	seed, err := manager.GetSeed(timestamp)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(700 * time.Millisecond)
	other, err := manager.GetSeed(timestamp)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(seed, other) {
		t.Fatal("should be the same")
	}
}

func TestGetFailedExpired(t *testing.T) {
	timestamp := time.Now()
	time.Sleep(1200 * time.Millisecond)
	_, err := manager.GetSeed(timestamp)
	if err == nil {
		t.Fatal("should fail but got succeed")
	}
	t.Logf("got err: %s", err.Error())
}
