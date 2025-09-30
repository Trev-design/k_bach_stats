package sidecar

import "sync"

type Request interface {
	// here we get our payload from the reques
	GetPayload() Payload

	// because we have no destructors we must make the request done explicitly with this function
	Done()
}

type Payload interface {
	// gives you the payload bytes if you use some json based infrastructures
	GetPayloadBytes() []byte
}

type BackgroundSignal interface {
	// bridges the functionality of a background service to the sidecar
	// and you can send some messages in the background
	HandleBackgroundProcess(req Request)
}

type ResponsiveSignal interface {
	// bridges the functionality of a responsive service to the sidecar
	// and you can send some responsive messages
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
