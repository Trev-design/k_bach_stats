package core

type Listener interface {
	Consume() error
	CloseListener() error
	Disconnect() error
	Wait()
}
