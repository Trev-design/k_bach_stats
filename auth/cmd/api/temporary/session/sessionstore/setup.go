package sessionstore

import (
	"crypto/tls"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	expiry time.Duration
}

func NewRedisClient(expiry time.Duration, password, host, port string, tlsConfig *tls.Config) (*RedisClient, error) {
	client, err := setupClient(password, host, port, tlsConfig)
	if err != nil {
		return nil, err
	}

	return &RedisClient{
		client: client,
		expiry: expiry,
	}, nil
}

func (client *RedisClient) CloseRedisStore() error {
	return client.client.Close()
}
