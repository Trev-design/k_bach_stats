package credentialdistro

import (
	"auth_server/cmd/api/grpc/credentialdistro/proto"
	"auth_server/cmd/api/utils/connection"
	"encoding/base64"
	"encoding/json"
)

func (handler *grpcNewCredsStreamHandler) handleStream(id string) {
	stream := handler.streams[id]

	for {
		response, err := stream.stream.Recv()
		if err != nil {

		}

		handler.handleIncome(response)
	}
}

func (handler *grpcNewCredsStreamHandler) handleIncome(response *proto.Response) {
	switch resp := response.Value.(type) {
	case *proto.Response_Response:
		handler.handleResponse(resp)

	case *proto.Response_Confirm:
		handler.handleConfirm(resp)

	case *proto.Response_InitialResponse:
		handler.handleInitialResponse(resp)
	}
}

func (handler *grpcNewCredsStreamHandler) handleResponse(response *proto.Response_Response) {
	resp := response.Response
	stream, ok := handler.streams[resp.Id]
	if !ok {
		// TODO: implement reconnect mechanism
		return
	}

	if resp.Topic != stream.topic {
		handler.handleNack(stream, "FALSE_TOPIC_IN_RESPONSE")
		return
	}

	payload, err := base64.RawURLEncoding.DecodeString(resp.Payload)
	if err != nil {
		handler.handleNack(stream, "INVALID_PAYLOAD_IN_RESPONSE")
		return
	}

	decrypted, err := handler.keyManager.Decrypt(payload)
	if err != nil {
		handler.handleNack(stream, "INVALID_PAYLOAD_IN_RESPONSE")
		return
	}

	credentials := new(connection.Credentials)

	if err := json.Unmarshal(decrypted, credentials); err != nil {
		handler.handleNack(stream, "INVALID_PAYLOAD_IN_RESPONSE")
		return
	}

	stream.messageChannels <- credentials

	handler.handleAck(stream)
}

func (handler *grpcNewCredsStreamHandler) handleConfirm(response *proto.Response_Confirm) {
	resp := response.Confirm
	stream, ok := handler.streams[resp.Id]
	if !ok {
		// TODO: implement reconnect mechanism
		return
	}

	if stream.topic != resp.Topic {
		handler.handleNack(stream, "INVALID_TOPIC_IN_CONFIRM")
		return
	}

	if handler.isScheduled() {
		handler.handleConfirmAck(stream)
		return
	}
}

func (handler *grpcNewCredsStreamHandler) handleInitialResponse(response *proto.Response_InitialResponse) {
	resp := response.InitialResponse
	stream, ok := handler.streams[resp.Id]
	if !ok {
		// TODO: implement reconnect mechanism
		return
	}

	if stream.topic != resp.Topic {
		handler.handleNack(stream, "FALSE_TOPIC_IN_INITIAL")
		return
	}

	payload, err := base64.RawURLEncoding.DecodeString(resp.Payload)
	if err != nil {
		handler.handleNack(stream, "INVALID_PAYLOAD_IN_INIT")
		return
	}

	decripted, err := handler.keyManager.Decrypt(payload)
	if err != nil {
		handler.handleNack(stream, "INVALID_PAYLOAD_IN_INIT")
		return
	}

	credentials := new(connection.Credentials)

	if err := json.Unmarshal(decripted, credentials); err != nil {
		handler.handleNack(stream, "INVALID_PAYLOAD_IN_INIT")
		return
	}

	stream.messageChannels <- credentials

	handler.handleConfirmAck(stream)
}

func (handler *grpcNewCredsStreamHandler) handleAck(stream *grpcNewCredsStream) {
	if err := stream.stream.Send(&proto.NewCredsRequest{
		Value: &proto.NewCredsRequest_AckNewCreds{
			AckNewCreds: &proto.AckNewCreds{
				Topic: stream.topic,
				Id:    stream.id,
			},
		},
	}); err != nil {
		// TODO: implement reconnect mechanism
		return
	}

	handler.confirm(stream.id)
}

func (handler *grpcNewCredsStreamHandler) handleNack(stream *grpcNewCredsStream, message string) {
	if err := stream.stream.Send(&proto.NewCredsRequest{
		Value: &proto.NewCredsRequest_NackNewCreds{
			NackNewCreds: &proto.NackNewCreds{
				Topic: stream.topic,
				Id:    stream.id,
				Msg:   message,
			},
		},
	}); err != nil {
		// TODO: implement reconnect mechanism
		return
	}
}

func (handler *grpcNewCredsStreamHandler) handleConfirmAck(stream *grpcNewCredsStream) {
	key, err := handler.keyManager.SwapAndGet()
	if err != nil {
		// TODO: implement reconnect mechanism
		return
	}

	if err := stream.stream.Send(&proto.NewCredsRequest{
		Value: &proto.NewCredsRequest_AckConfirm{
			AckConfirm: &proto.AckConfirm{
				Id:        stream.id,
				Topic:     stream.topic,
				PublicKey: key,
			},
		},
	}); err != nil {
		// TODO: implement reconnect mechanism
		return
	}
}
