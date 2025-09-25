package mocks

import (
	"auth_server/cmd/api/sidecar"
	"log"
	"time"
)

type LoggingProcessor struct{}
type LoggingPayload []byte

func (processor *LoggingProcessor) HandleBackgroundProcess(req sidecar.Request) {
	defer req.Done()
	payload := req.GetPayload()
	time.Sleep(50 * time.Millisecond)
	log.Println(string(payload.GetPayloadBytes()))
}

func (payload LoggingPayload) GetPayloadBytes() []byte {
	return payload
}
