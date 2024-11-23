package application

import (
	"user_manager/internal/core"
)

type ListenerServiceAdapter struct {
	Listener core.Listener
}

func NewListenerServiceAdapter(listener core.Listener) *ListenerServiceAdapter {
	return &ListenerServiceAdapter{
		Listener: listener,
	}
}

func (lsa *ListenerServiceAdapter) Start() error {
	return lsa.Listener.Consume()
}

func (lsa *ListenerServiceAdapter) ShutDown() error {
	if err := lsa.Listener.Disconnect(); err != nil {
		return err
	}

	lsa.Listener.Wait()

	if err := lsa.Listener.CloseListener(); err != nil {
		return err
	}

	return nil
}
