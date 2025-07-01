package sessioncore_test

import (
	"auth_server/cmd/api/temporary/session/sessioncore"
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type sessionCreds struct {
	password string
	host     string
	port     string
}

var creds *sessionCreds
var session *sessioncore.Session

func TestMain(t *testing.M) {
	newCreds, cancel := redisTestContainer()
	creds = new(sessionCreds)
	creds = newCreds

	newSession, err := setup()
	if err != nil {
		log.Fatal(err)
	}
	session = newSession

	session.HandleBackground()

	code := t.Run()

	if err := session.CloseSession(); err != nil {
		log.Fatal(err)
	}

	cancel()

	os.Exit(code)
}

func TestSetData(t *testing.T) {
	if _, _, err := session.SetVerifyData(uuid.NewString()); err != nil {
		t.Fatal(err)
	}
}

func setup() (*sessioncore.Session, error) {
	return sessioncore.NewSessionBuilder().
		Host(creds.host).
		Port(creds.port).
		Password(creds.password).
		IntevalDuration(2 * time.Second).
		Build()
}

func redisTestContainer() (*sessionCreds, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		Cmd: []string{
			"redis-server",
			"--requirepass",
			"testpass",
		},
		WaitingFor: wait.ForLog("Ready to accept connections"),
	}

	instance, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req, Started: true,
		})
	if err != nil {
		log.Fatal(err)
	}

	host, err := instance.Host(ctx)
	if err != nil {
		log.Fatal()
	}

	port, err := instance.MappedPort(ctx, "6379/tcp")
	if err != nil {
		log.Fatal(err)
	}

	return &sessionCreds{
		port:     port.Port(),
		host:     host,
		password: "testpass",
	}, func() { instance.Terminate(ctx) }
}
