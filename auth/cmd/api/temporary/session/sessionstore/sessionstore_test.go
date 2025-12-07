package sessionstore_test

import (
	"auth_server/cmd/api/temporary/session/sessionstore"
	"context"
	"errors"
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

func Test_SessionStore(t *testing.T) {
	redisInstance(t)
}

func redisInstance(t *testing.T) {
	t.Run("redis_instance", func(t *testing.T) {
		sessionCreds, cancel := setupRedisContainer(t)
		defer cancel()

		client := newRedisClientTest(t, sessionCreds)
		id := setPayloadTest(t, client)
		getPayloadTest(t, client, id)
		deletePayloadTest(t, client, id)
		getPayloadFailedTest(t, client)
		deletePayloadFailedTest(t, client)
		closeClientTest(t, client)
	})
}

func setupRedisContainer(t *testing.T) (*sessionCreds, func()) {
	ctx := context.Background()

	var redisInstance testcontainers.Container
	var redisHost string
	var redisPort string

	t.Run("setup_redis_container", func(t *testing.T) {
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
			t.Fatal(err)
		}

		host, err := instance.Host(ctx)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("my very cool host is %s", host)

		port, err := instance.MappedPort(ctx, "6379/tcp")
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("my very cool port is %s", port.Port())

		redisInstance = instance
		redisHost = host
		redisPort = port.Port()

		time.Sleep(2 * time.Second)
	})

	return &sessionCreds{
		password: "testpass",
		host:     redisHost,
		port:     redisPort,
	}, func() { redisInstance.Terminate(ctx) }
}

func newRedisClientTest(t *testing.T, creds *sessionCreds) *sessionstore.RedisClient {
	var redisSession *sessionstore.RedisClient

	t.Run("new_redis_client", func(t *testing.T) {
		session, err := sessionstore.NewRedisClientBuilder().
			Host(creds.host).
			Port(creds.port).
			Password(creds.password).
			WithDuration(2 * time.Second).
			Build()
		if err != nil {
			t.Fatal(err)
		}

		redisSession = session
	})

	return redisSession
}

func setPayloadTest(t *testing.T, client *sessionstore.RedisClient) string {
	var sessionID string

	t.Run("redis_set_test", func(t *testing.T) {
		id, err := client.SetSessionPayload("verify", "my_super_mega_verify_payload")
		if err != nil {
			t.Fatal(err)
		}

		sessionID = id
		t.Log(sessionID)
	})
	t.Log(sessionID)

	return sessionID
}

func getPayloadTest(t *testing.T, client *sessionstore.RedisClient, id string) {
	t.Run("redis_get_test", func(t *testing.T) {
		payload, err := client.GetSessionPayload("verify", id)
		if err != nil {
			t.Fatal(err)
		}

		if payload != "my_super_mega_verify_payload" {
			t.Fatal(errors.New("invalid payload"))
		}
	})
}

func deletePayloadTest(t *testing.T, client *sessionstore.RedisClient, id string) {
	t.Run("redis_delete_test", func(t *testing.T) {
		if err := client.DeleteSessionPayload("verify", id); err != nil {
			t.Fatal(err)
		}
	})
}

func getPayloadFailedTest(t *testing.T, client *sessionstore.RedisClient) {
	t.Run("redis_set_failed_test", func(t *testing.T) {
		if _, err := client.SetSessionPayload("verify", "payload"); err != nil {
			t.Fatal(err)
		}

		if _, err := client.GetSessionPayload("verify", uuid.NewString()); err == nil {
			t.Fatal("should fail but succeed")
		}
	})
}

func deletePayloadFailedTest(t *testing.T, client *sessionstore.RedisClient) {
	t.Run("redis_delete_failed_test", func(t *testing.T) {
		if _, err := client.SetSessionPayload("verify", "payload"); err != nil {
			t.Fatal(err)
		}

		if err := client.DeleteSessionPayload("verify", uuid.NewString()); err == nil {
			t.Fatal(errors.New("should fail but succeed"))
		}
	})
}

func closeClientTest(t *testing.T, client *sessionstore.RedisClient) {
	t.Run("close_redis_test", func(t *testing.T) {
		if err := client.CloseRedisStore(); err != nil {
			t.Fatal(err)
		}
	})
}
