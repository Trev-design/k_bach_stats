package userregistry

func (client *GRPCClient) handleMessageLoad() {
	for {
		select {
		case <-client.promptRequestChannel:
			client.handleSendPolicy()
		case <-client.primaryRequestDoneChannel:
			client.primaryReadyChannel <- struct{}{}
		case <-client.overflowReadyChannel:
			client.handleSendPolicy()
		}
	}
}

func (client *GRPCClient) handleMessage() {
	for message := range client.messageIncomeChannel {
		client.promptRequestChannel <- struct{}{}

		select {
		case <-client.primaryReadyChannel:
			client.messagePrimaryChannel <- message

		case <-client.overflowReadyChannel:
			client.messageOverflowChannel <- message
		}
	}
}

func (client *GRPCClient) computePrimaryStream() {
	client.handlePrimaryStream(
		client.maxNumPrimaryRequests,
		client.primaryRequestDoneChannel,
		client.messagePrimaryChannel,
	)
}

func (client *GRPCClient) computeOverflowStream() {
	client.handleOverflowStream(
		client.maxNumOverflowRequests,
		client.overflowRequestDoneChannel,
		client.messageOverflowChannel,
	)
}

func (client *GRPCClient) handleSendPolicy() {
	if client.currentNumOverflowRequests.Load() < client.maxNumOverflowRequests && client.currentNumPrimaryRequests.Load() >= client.maxNumPrimaryRequests {
		client.overflowReadyChannel <- struct{}{}
	} else if client.currentNumPrimaryRequests.Load() < client.maxNumPrimaryRequests {
		client.primaryReadyChannel <- struct{}{}
	}
}
