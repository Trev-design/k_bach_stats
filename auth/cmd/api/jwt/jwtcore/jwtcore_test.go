package jwtcore_test

import (
	"auth_server/cmd/api/jwt/jwtcore"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

var core *jwtcore.JWTService

func TestMain(m *testing.M) {
	newCore := jwtcore.NewJWTServiceBuilder().
		Identifier(uuid.New()).
		Interval(time.Second).
		Build()
	core = newCore

	core.ComputeBackgroundService()

	code := m.Run()

	core.CloseJWTService()

	os.Exit(code)
}

func TestInit(t *testing.T) {
	newCore := jwtcore.NewJWTServiceBuilder().
		Identifier(uuid.New()).
		Interval(time.Second).
		Build()

	if err := newCore.CloseJWTService(); err != nil {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	token, err := core.Sign(uuid.NewString(), uuid.NewString())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("token: %s", token)
}
