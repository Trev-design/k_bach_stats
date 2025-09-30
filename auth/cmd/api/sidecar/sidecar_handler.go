package sidecar

import (
	"errors"
	"sync"
	"sync/atomic"
)

// sends your payload to a background process it will return an error if the topic doesn't exist
func (processor *processor) MakeBackgroundProcess(name string, message Payload) error {
	topic, ok := processor.backgroundTopics[name]
	if !ok {
		return errors.New("no topic with this name")
	}

	topic.sendBackgroundMessage(message)

	return nil
}

// sends your payload to process which will be responsive.
// it will return an error if the topic doesn't exist or if the request doesn't succeed
func (processor *processor) MakeResponsiveProcess(name string, message Payload) error {
	topic, ok := processor.responsiveTopics[name]
	if !ok {
		return errors.New("no topic with this name")
	}

	return topic.sendResponsiveMessage(message)
}

// here we can register some background topics
func (processor *processor) RegisterBackgroundTopic(signal BackgroundSignal, name string) {
	processor.backgroundTopics[name] = &backgroundTopic{
		messageChannel: make(chan Payload, 5),
		handler:        signal,
		isFinished:     atomic.Bool{},
		waitgroup:      processor.waitgroup,
		once:           sync.Once{},
	}
}

// here we can register some responsive topics
func (processor *processor) RegisterResponsiveTopic(signal ResponsiveSignal, name string) {
	processor.responsiveTopics[name] = &responsiveTopic{
		messageChannel: make(chan responsiveRequest, 5),
		handle:         signal,
		isFinished:     atomic.Bool{},
		waitgroup:      processor.waitgroup,
		once:           sync.Once{},
	}
}

// here we starting the topic processes
func (processor *processor) StartProcessing() {
	for _, topic := range processor.backgroundTopics {
		go topic.handleMessages()
	}

	for _, topic := range processor.responsiveTopics {
		go topic.handleMessages()
	}
}

// here we're stopping the sidecar on server shutdown
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
