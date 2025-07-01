package sessionstore

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func setupClient(password, host, port string, tlsConfig *tls.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Network:   "tcp",
		Addr:      fmt.Sprintf("%s:%s", host, port),
		Password:  password,
		TLSConfig: tlsConfig,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func (client *RedisClient) makeNewID(kind string) string {
	for {
		guid := uuid.New().String()
		redisID := fmt.Sprintf("%s:%s", kind, guid)
		status, err := client.client.Exists(
			context.Background(),
			redisID,
		).Result()

		if err != nil {
			continue
		}

		if status > 0 {
			continue
		}

		return redisID
	}
}

func makeID(kind, guid string) string {
	return fmt.Sprintf("%s:%s", kind, guid)
}

func getUUIDFromRedisID(redisID string) (string, error) {
	sequences := strings.Split(redisID, ":")

	if len(sequences) != 2 {
		return "", errors.New("invalid id")
	}

	if _, err := uuid.Parse(sequences[1]); err != nil {
		return "", errors.New("invalid id")
	}

	return sequences[1], nil
}
