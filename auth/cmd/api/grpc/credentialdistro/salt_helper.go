package credentialdistro

import (
	"auth_server/cmd/api/grpc/credentialdistro/proto"
	"encoding/base64"
)

func (stream *grpcSaltStream) handleStream() {
	for {
		resp, err := stream.stream.Recv()
		if err != nil {
			// TODO: implement reconnect mechanism
			break
		}

		stream.handleIncome(resp)
	}
}

func (stream *grpcSaltStream) handleIncome(response *proto.Response) {
	switch resp := response.Value.(type) {
	case *proto.Response_Response:
		stream.handleResponse(resp)

	case *proto.Response_Confirm:
		stream.handleConfirm(resp)

	default:
		stream.handleNack("INVALID_RESPONSE_TYPE")
	}
}

func (stream *grpcSaltStream) handleResponse(response *proto.Response_Response) {
	resp := response.Response

	if stream.id != resp.Id {
		stream.handleNack("INVALID_ID_IN_RESPONSE")
		return
	}

	if stream.topic != resp.Topic {
		stream.handleNack("INVALID_TOPIC_IN_RESPONSE")
		return
	}

	saltBytes, err := base64.RawURLEncoding.DecodeString(resp.Payload)
	if err != nil {
		stream.handleNack("INVALID_SALT_BYTES")
		return
	}

	stream.messageChannel <- saltBytes

	stream.handleAck()
}

func (stream *grpcSaltStream) handleConfirm(response *proto.Response_Confirm) {
	resp := response.Confirm

	if stream.id != resp.Id {
		stream.handleNack("INVALID_ID_IN_CONFIRM")
		return
	}

	if stream.topic != resp.Topic {
		stream.handleNack("INVALID_TOPIC_IN_CONFIRM")
		return
	}
}

func (stream *grpcSaltStream) handleNack(message string) {
	if err := stream.stream.Send(&proto.SaltRequest{
		Value: &proto.SaltRequest_NackSalt{
			NackSalt: &proto.NackSalt{
				Id:    stream.id,
				Topic: stream.topic,
				Msg:   message,
			},
		},
	}); err != nil {
		// TODO: implement reconnect mechanism
		return
	}
}

func (stream *grpcSaltStream) handleAck() {
	if err := stream.stream.Send(&proto.SaltRequest{
		Value: &proto.SaltRequest_AckSalt{
			AckSalt: &proto.AckSalt{
				Id:    stream.id,
				Topic: stream.topic,
			},
		},
	}); err != nil {
		// TODO: implement reconnect mechanism
		return
	}
}
