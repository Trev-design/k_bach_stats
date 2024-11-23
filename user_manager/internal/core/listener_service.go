package core

type ListenerService interface {
	Start() error
	ShutDown() error
}
