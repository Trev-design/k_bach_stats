package userregistry

func (client *GRPCClient) SendMessage(message Message) error {
	errorChannel := make(chan error)
	client.messageIncomeChannel <- Request{
		Message:  message,
		Response: errorChannel,
	}

	return <-errorChannel
}
