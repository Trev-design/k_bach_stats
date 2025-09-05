package mocks

import (
	"auth_server/cmd/api/grpc/userregistry/proto"
	"auth_server/cmd/api/utils"
	"context"
	"io"
	"net"
	"time"

	"golang.org/x/sync/semaphore"
	"google.golang.org/grpc"
)

type RegistryServer struct {
	proto.UserRegistryServiceServer
	maxPrimary  int64
	maxOverflow int64
}

type responseHandler struct {
	stream    grpc.BidiStreamingServer[proto.RegistryRequest, proto.RegistryResponse]
	semaphore *semaphore.Weighted
	queue     *utils.MessageQueue[chan *proto.RegistryResponse]
}

type messageHandler struct {
	responseChannel chan *proto.RegistryResponse
	semaphore       *semaphore.Weighted
	request         *proto.RegistryRequest
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
	semaphore := semaphore.NewWeighted(server.maxPrimary)
	queue := utils.NewMessageQueue[chan *proto.RegistryResponse]()

	handler := &responseHandler{
		stream:    stream,
		semaphore: semaphore,
		queue:     queue,
	}

	return handler.handleResponses()
}

func (server *RegistryServer) UserOverflowStream(stream grpc.BidiStreamingServer[proto.RegistryRequest, proto.RegistryResponse]) error {
	semaphore := semaphore.NewWeighted(server.maxOverflow)
	queue := utils.NewMessageQueue[chan *proto.RegistryResponse]()

	handler := &responseHandler{
		stream:    stream,
		semaphore: semaphore,
		queue:     queue,
	}

	return handler.handleResponses()
}

func (handler *responseHandler) handleResponses() error {
	for {
		request, err := handler.stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		responseChannel := make(chan *proto.RegistryResponse)
		handler.queue.Enqueue(responseChannel)
		handler.semaphore.Acquire(context.Background(), 1)

		messageHandler := &messageHandler{
			request:         request,
			responseChannel: responseChannel,
			semaphore:       handler.semaphore,
		}

		go messageHandler.handleMessage()
	}
}

func (handler *messageHandler) handleMessage() {
	defer handler.semaphore.Release(1)

	time.Sleep(200 * time.Millisecond)

	handler.responseChannel <- &proto.RegistryResponse{
		Status: handler.request.Name,
	}
}
