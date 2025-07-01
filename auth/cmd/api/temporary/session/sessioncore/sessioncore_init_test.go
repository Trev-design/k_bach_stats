package sessioncore_test

import (
	"auth_server/cmd/api/temporary/session/sessioncore"
	"testing"
	"time"
)

func TestInitFailedFalseHost(t *testing.T) {
	_, err := sessioncore.NewSessionBuilder().
		Host("hoho").
		Port(creds.port).
		Password(creds.password).
		IntevalDuration(2 * time.Second).
		Build()

	if err == nil {
		t.Fatal("should fail because of false host but got succeed")
	}

	t.Logf("the err: %s", err.Error())
}

func TestInitFailedFalsePort(t *testing.T) {
	_, err := sessioncore.NewSessionBuilder().
		Host(creds.host).
		Port("haha").
		Password(creds.password).
		IntevalDuration(2 * time.Second).
		Build()

	if err == nil {
		t.Fatal("should fail because of false port but got succeed")
	}

	t.Logf("the err: %s", err.Error())
}

func TestInitFailedFalsePassword(t *testing.T) {
	_, err := sessioncore.NewSessionBuilder().
		Host(creds.host).
		Port(creds.port).
		Password("affe").
		IntevalDuration(2 * time.Second).
		Build()

	if err == nil {
		t.Fatal("should fail because of false port but got succeed")
	}

	t.Logf("the err: %s", err.Error())
}

func TestInitAndCloseSuccess(t *testing.T) {
	newSession, err := sessioncore.NewSessionBuilder().
		Host(creds.host).
		Port(creds.port).
		Password(creds.password).
		IntevalDuration(2 * time.Second).
		Build()

	if err != nil {
		t.Fatal("should fail because of false port but got succeed")
	}

	if err := newSession.CloseSession(); err != nil {
		t.Fatal(err)
	}
}
