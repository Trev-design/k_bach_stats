package credentialdistro

import (
	"auth_server/cmd/api/grpc/credentialdistro/proto"
)

func (stream *grpcSaltStream) subscribe() {
	go stream.handleStream()

	stream.stream.Send(&proto.SaltRequest{
		Value: &proto.SaltRequest_SubToSalt{
			SubToSalt: &proto.SubscribeToSaltDistro{
				Topic: stream.topic,
				Id:    stream.id,
			},
		},
	})
}

func (stream *grpcNewCredsStreamHandler) subscribe(id string, streamChannel *grpcNewCredsStream) {
	stream.waitgroup.Wait()
	stream.waitgroup.Add(1)

	stream.streams[id] = streamChannel

	go stream.handleStream(id)

	if err := streamChannel.stream.Send(&proto.NewCredsRequest{
		Value: &proto.NewCredsRequest_SubToCredRotator{
			SubToCredRotator: &proto.SubscribeToConnectionCredRotator{
				Topic: streamChannel.topic,
				Id:    streamChannel.id,
			},
		},
	}); err != nil {
		return
	}

	stream.register(id)
}
