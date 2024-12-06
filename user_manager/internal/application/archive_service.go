package application

import (
	"encoding/json"
	"log"
	"sync"
	"time"
	"user_manager/internal/core"
	"user_manager/types"
)

type ArchiveServiceAdapter struct {
	store        core.NotificationStore
	stream       core.Stream
	doneChannel  chan bool
	readyChannel chan bool
	errorChannel chan error
	waitgroup    *sync.WaitGroup
}

func NewArchiveStreamServiceAdapter(
	store core.NotificationStore,
	stream core.Stream,
) *ArchiveServiceAdapter {
	return &ArchiveServiceAdapter{
		store:        store,
		stream:       stream,
		errorChannel: make(chan error),
		readyChannel: make(chan bool),
		doneChannel:  make(chan bool),
		waitgroup:    &sync.WaitGroup{},
	}
}

func (asa *ArchiveServiceAdapter) StartArchiveLoop() {
	for {
		go asa.makeReady()
		select {
		case <-asa.readyChannel:
			asa.makeStreams()

		case err := <-asa.errorChannel:
			log.Println(err.Error())

		case <-asa.doneChannel:
			return
		}
	}
}

func (asa *ArchiveServiceAdapter) makeStreams() {
	go asa.makeInvitationStream()
	go asa.makeJoinRequestStream()
}

func (asa *ArchiveServiceAdapter) makeInvitationStream() {
	ids, payload, err := asa.store.GetInvitations()
	if err != nil {
		asa.errorChannel <- err
		return
	}

	sendData := make([][]byte, len(payload))

	for _, item := range payload {
		jsonData, err := json.Marshal(item)
		if err != nil {
			asa.errorChannel <- err
			continue
		}

		sendData = append(sendData, jsonData)
	}

	if err := asa.stream.SendStream(
		&types.StreamPayload{
			Kind:    "join_request",
			Payload: sendData,
		},
	); err != nil {
		asa.errorChannel <- err
	}

	if err := asa.store.RemoveSelectedInvitationData(ids); err != nil {
		asa.errorChannel <- err
	}
}

func (asa *ArchiveServiceAdapter) makeJoinRequestStream() {
	ids, payload, err := asa.store.GetJoinRequests()
	if err != nil {
		asa.errorChannel <- err
		return
	}

	sendData := make([][]byte, len(payload))

	for _, item := range payload {
		jsonData, err := json.Marshal(item)
		if err != nil {
			asa.errorChannel <- err
			continue
		}

		sendData = append(sendData, jsonData)
	}

	if err := asa.stream.SendStream(
		&types.StreamPayload{
			Kind:    "join_request",
			Payload: sendData,
		},
	); err != nil {
		asa.errorChannel <- err
	}

	if err := asa.store.RemoveSelectedJoinRequestData(ids); err != nil {
		asa.errorChannel <- err
	}
}

func (asa *ArchiveServiceAdapter) Close() error {
	asa.doneChannel <- true
	asa.waitgroup.Wait()
	return asa.stream.Close()
}

func (asa *ArchiveServiceAdapter) makeReady() {
	time.Sleep(5 * time.Minute)
	asa.readyChannel <- true
}
