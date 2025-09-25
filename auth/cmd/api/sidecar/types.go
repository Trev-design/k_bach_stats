package sidecar

import "sync"

type Request interface {
	GetPayload() Payload
	Done()
}

type Payload interface {
	GetPayloadBytes() []byte
}

type BackgroundSignal interface {
	HandleBackgroundProcess(req Request)
}

type ResponsiveSignal interface {
	HandleResponsiveProcess(req Request, errChan chan error)
}

type request struct {
	message   Payload
	waitgroup *sync.WaitGroup
}

func (req *request) GetPayload() Payload {
	return req.message
}

func (req *request) Done() {
	req.waitgroup.Done()
}

type closeMessage []byte

func (msg closeMessage) GetPayloadBytes() []byte {
	return msg
}
