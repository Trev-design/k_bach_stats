package sessionstore

import (
	"context"
	"errors"
)

// adds session data to the session store
func (client *RedisClient) SetSessionPayload(service, payload string) (string, error) {
	conn := client.conn.Get()
	conn.waitgroup.Add(1)
	defer conn.waitgroup.Done()

	id := client.makeNewID(service, conn)

	if _, err := conn.client.Set(
		context.Background(),
		id,
		payload,
		client.expiry,
	).Result(); err != nil {
		return "", err
	}

	return getUUIDFromRedisID(id)
}

// gets data from the session store
func (client *RedisClient) GetSessionPayload(service, guid string) (string, error) {
	conn := client.conn.Get()
	conn.waitgroup.Add(1)
	defer conn.waitgroup.Done()

	id := makeID(service, guid)

	return conn.client.Get(context.Background(), id).Result()
}

// deletes data from the session store
func (client *RedisClient) DeleteSessionPayload(service, guid string) error {
	conn := client.conn.Get()
	conn.waitgroup.Add(1)
	defer conn.waitgroup.Done()

	id := makeID(service, guid)

	numDeleted, err := conn.client.Del(context.Background(), id).Result()
	if err != nil {
		return err
	}

	if numDeleted == 0 {
		return errors.New("key not found")
	}

	return nil
}
