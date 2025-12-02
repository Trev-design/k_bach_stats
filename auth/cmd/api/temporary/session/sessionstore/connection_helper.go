package sessionstore

func (connection *RedisConnection) Wait() {
	connection.waitgroup.Wait()
}

func (connection *RedisConnection) Close() error {
	return connection.client.Close()
}

func (connection *RedisConnection) Conn() *RedisConnection {
	return connection
}
