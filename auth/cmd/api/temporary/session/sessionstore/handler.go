package sessionstore

import (
	"context"
	"errors"
)

func (client *RedisClient) SetSessionPayload(service, payload string) (string, error) {
	id := client.makeNewID(service)

	if _, err := client.client.Set(
		context.Background(),
		id,
		payload,
		client.expiry,
	).Result(); err != nil {
		return "", err
	}

	return getUUIDFromRedisID(id)
}

func (client *RedisClient) GetSessionPayload(service, guid string) (string, error) {
	id := makeID(service, guid)

	return client.client.Get(context.Background(), id).Result()
}

func (client *RedisClient) DeleteSessionPayload(service, guid string) error {
	id := makeID(service, guid)

	numDeleted, err := client.client.Del(context.Background(), id).Result()
	if err != nil {
		return err
	}

	if numDeleted == 0 {
		return errors.New("key not found")
	}

	return nil
}
