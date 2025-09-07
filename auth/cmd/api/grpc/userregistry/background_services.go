package userregistry

func (client *GRPCClient) handleMessage() {
	for message := range client.messageIncomeChannel {
		if client.isOverflowFreeToUse() {
			client.messageOverflowChannel <- message
		} else {
			client.messagePrimaryChannel <- message
		}
	}
}

func (client *GRPCClient) computePrimaryStream() {
	client.handlePrimaryStream()
}

func (client *GRPCClient) computeOverflowStream() {
	client.handleOverflowStream()
}

func (client *GRPCClient) isOverflowFull() bool {
	return client.currentNumOverflowRequests.Load() >= client.maxNumOverflowRequests
}

func (client *GRPCClient) isPrimaryFull() bool {
	return client.currentNumPrimaryRequests.Load() >= client.maxNumPrimaryRequests
}

func (client *GRPCClient) isOverflowFreeToUse() bool {
	return !client.isOverflowFull() && client.isPrimaryFull()
}
