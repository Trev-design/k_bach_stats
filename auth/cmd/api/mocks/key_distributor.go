package mocks

import (
	"auth_server/cmd/api/sidecar"
	"errors"
	"log"
	"time"
)

type KeyDistributor struct{}

func (distributor *KeyDistributor) HandleResponsiveProcess(req sidecar.Request, errChan chan error) {
	defer req.Done()
	payload := req.GetPayload()
	time.Sleep(50 * time.Millisecond)
	log.Println(string(payload.GetPayloadBytes()))
	if string(payload.GetPayloadBytes()) == "some invalid payload" {
		errChan <- errors.New("you've send some invalid payload")
	} else {
		errChan <- nil
	}
}
