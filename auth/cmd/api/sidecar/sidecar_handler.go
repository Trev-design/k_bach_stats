package sidecar

import (
	"errors"
	"sync"
	"sync/atomic"
)

func (processor *processor) MakeBackgroundProcess(name string, message Payload) error {
	topic, ok := processor.backgroundTopics[name]
	if !ok {
		return errors.New("no topic with this name")
	}

	topic.sendBackgroundMessage(message)

	return nil
}

func (processor *processor) MakeResponsiveProcess(name string, message Payload) error {
	topic, ok := processor.responsiveTopics[name]
	if !ok {
		return errors.New("no topic with this name")
	}

	return topic.sendResponsiveMessage(message)
}

func (processor *processor) RegisterBackgroundTopic(signal BackgroundSignal, name string) {
	processor.backgroundTopics[name] = &backgroundTopic{
		messageChannel: make(chan Payload, 5),
		handler:        signal,
		isFinished:     atomic.Bool{},
		waitgroup:      processor.waitgroup,
		once:           sync.Once{},
	}
}

func (processor *processor) RegisterResponsiveTopic(signal ResponsiveSignal, name string) {
	processor.responsiveTopics[name] = &responsiveTopic{
		messageChannel: make(chan responsiveRequest, 5),
		handle:         signal,
		isFinished:     atomic.Bool{},
		waitgroup:      processor.waitgroup,
		once:           sync.Once{},
	}
}

func (processor *processor) StartProcessing() {
	for _, topic := range processor.backgroundTopics {
		go topic.handleMessages()
	}

	for _, topic := range processor.responsiveTopics {
		go topic.handleMessages()
	}
}

func (processor *processor) StopSidecar() {
	for _, topic := range processor.backgroundTopics {
		topic.stopTopic()
	}

	for _, topic := range processor.responsiveTopics {
		topic.stopTopic()
	}

	processor.waitgroup.Wait()

	for key := range processor.backgroundTopics {
		delete(processor.backgroundTopics, key)
	}

	for key := range processor.responsiveTopics {
		delete(processor.responsiveTopics, key)
	}
}
