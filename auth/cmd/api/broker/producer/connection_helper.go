package producer

func (conn *RMQConnection) Wait() {
	conn.waitgroup.Wait()
}

func (conn *RMQConnection) Close() error {
	for _, channel := range conn.channels {
		if err := channel.CloseChannel(); err != nil {
			return err
		}
	}

	return conn.connection.Close()
}

func (conn *RMQConnection) Conn() *RMQConnection {
	return conn
}
