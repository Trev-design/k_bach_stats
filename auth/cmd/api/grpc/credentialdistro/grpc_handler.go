package credentialdistro

import (
	"auth_server/cmd/api/utils/connection"
	"context"
)

func (client *GRPCClient) MakeSaltStream(topic, id string, messageChannel chan []byte) {

}

func (client *GRPCClient) MakeNewCredsStream(topic, id string, messageChannel chan *connection.Credentials) error {
	stream, err := client.client.NewCredsStream(context.Background())
	if err != nil {
		return err
	}

	streamChannel := &grpcNewCredsStream{
		topic:           topic,
		id:              id,
		stream:          stream,
		messageChannels: messageChannel,
	}

	client.newCredsStream.subscribe(id, streamChannel)

	return nil
}
