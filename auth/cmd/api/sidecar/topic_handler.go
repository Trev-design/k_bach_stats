package sidecar

import "io"

func (topic *backgroundTopic) sendBackgroundMessage(message Payload) {
	if !topic.isFinished.Load() {
		topic.waitgroup.Add(1)
		topic.messageChannel <- message
	} else {
		topic.once.Do(func() { close(topic.messageChannel) })
	}
}

func (topic *responsiveTopic) sendResponsiveMessage(payload Payload) error {
	if !topic.isFinished.Load() {
		topic.waitgroup.Add(1)
		errorChannel := make(chan error)
		topic.messageChannel <- responsiveRequest{message: payload, errorChannel: errorChannel}
		return <-errorChannel
	} else {
		topic.once.Do(func() { close(topic.messageChannel) })
		return io.EOF
	}
}

func (topic *backgroundTopic) handleMessages() {
	for message := range topic.messageChannel {
		go topic.handler.HandleBackgroundProcess(
			&request{
				message:   message,
				waitgroup: topic.waitgroup,
			})
	}
}

func (topic *responsiveTopic) handleMessages() {
	for message := range topic.messageChannel {
		go topic.handle.HandleResponsiveProcess(
			&request{
				message:   message.message,
				waitgroup: topic.waitgroup,
			}, message.errorChannel)
	}
}

func (topic *backgroundTopic) stopTopic() {
	topic.isFinished.Store(true)
	topic.sendBackgroundMessage(closeMessage([]byte("closing")))
}

func (topic *responsiveTopic) stopTopic() {
	topic.isFinished.Store(true)
	topic.sendResponsiveMessage(closeMessage([]byte("closing")))
}
