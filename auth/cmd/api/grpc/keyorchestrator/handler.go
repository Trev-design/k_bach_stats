package keyorchestrator

import "log"

func (orchestrator *GRPCClient) OpenPipe(message *GetKeySubscriptionType) {
	stream := newGRPCStream(message.ResponsePipe)

	go stream.handle(orchestrator.client)

	orchestrator.mutex.Lock()
	orchestrator.streams[message.Id] = stream
	orchestrator.mutex.Unlock()

	stream.sendSubscribe(message)
}

func (orchestrator *GRPCClient) Retry(message *RetryType) {
	orchestrator.mutex.RLock()
	defer orchestrator.mutex.RUnlock()

	stream, ok := orchestrator.streams[message.Id]
	if !ok {
		log.Println("something went wrong")
	}

	stream.sendRetry(message)
}
