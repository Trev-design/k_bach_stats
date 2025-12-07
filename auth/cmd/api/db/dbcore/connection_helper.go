package dbcore

func (conn *Connection) Wait() {
	conn.waitgroup.Wait()
}

func (conn *Connection) Close() error {
	conn.waitgroup.Wait()
	db, err := conn.conn.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

func (conn *Connection) Conn() *Connection {
	return conn
}
