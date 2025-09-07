package mocks

import (
	"auth_server/cmd/api/grpc/userregistry/proto"
	"io"
	"net"
	"time"

	"google.golang.org/grpc"
)

type RegistryServer struct {
	proto.UserRegistryServiceServer
	maxPrimary  int64
	maxOverflow int64
}

type requestHandler struct {
	stream          grpc.BidiStreamingServer[proto.RegistryRequest, proto.RegistryResponse]
	responseChannel chan *responseMessage
}

type messageHandler struct {
	index           uint64
	responseChannel chan *responseMessage
	request         *proto.RegistryRequest
}

type responseHandler struct {
	stream          grpc.BidiStreamingServer[proto.RegistryRequest, proto.RegistryResponse]
	responseChannel chan *responseMessage
}

type responseMessage struct {
	id       uint64
	response *proto.RegistryResponse
}

func NewRegistryServer(maxPrimary, maxOverflow int64) error {
	listener, err := net.Listen("tcp", ":5670")
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	proto.RegisterUserRegistryServiceServer(
		grpcServer,
		&RegistryServer{
			maxPrimary:  maxPrimary,
			maxOverflow: maxOverflow,
		})

	return grpcServer.Serve(listener)
}

func (server *RegistryServer) UserPrimaryStream(stream grpc.BidiStreamingServer[proto.RegistryRequest, proto.RegistryResponse]) error {
	responseChannel := make(chan *responseMessage, server.maxPrimary)

	handler := &requestHandler{
		stream:          stream,
		responseChannel: responseChannel,
	}

	return handler.handleRequests()
}

func (server *RegistryServer) UserOverflowStream(stream grpc.BidiStreamingServer[proto.RegistryRequest, proto.RegistryResponse]) error {
	responseChannel := make(chan *responseMessage, server.maxOverflow)

	handler := &requestHandler{
		stream:          stream,
		responseChannel: responseChannel,
	}

	return handler.handleRequests()
}

func (handler *requestHandler) handleRequests() error {
	responseHandler := &responseHandler{
		responseChannel: handler.responseChannel,
		stream:          handler.stream,
	}

	var index uint64 = 0

	go responseHandler.handleResponses()

	for {
		request, err := handler.stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		messageHandler := &messageHandler{
			index:           index,
			request:         request,
			responseChannel: handler.responseChannel,
		}

		go messageHandler.handleMessage()

		index++
	}
}

func (handler *messageHandler) handleMessage() {
	time.Sleep(200 * time.Millisecond)
	handler.responseChannel <- &responseMessage{
		id:       handler.index,
		response: &proto.RegistryResponse{Status: handler.request.Name},
	}
}

func (handler *responseHandler) handleResponses() {
	var index uint64 = 0
	pending := make(map[uint64]*proto.RegistryResponse)
	for response := range handler.responseChannel {
		pending[response.id] = response.response

		for {
			message, ok := pending[index]
			if !ok {
				break
			}

			handler.stream.Send(message)
			delete(pending, index)
			index++
		}
	}
}
