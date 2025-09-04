package userregistry

import (
	"auth_server/cmd/api/grpc/userregistry/proto"
	"context"
	"errors"
	"io"

	"github.com/gofiber/fiber/v2/log"
	"google.golang.org/grpc"
)

type registryResponseHandler struct {
	stream      grpc.BidiStreamingClient[proto.RegistryRequest, proto.RegistryResponse]
	errorPipe   chan chan error
	doneChannel chan struct{}
}

type registryRequestHandler struct {
	stream         grpc.BidiStreamingClient[proto.RegistryRequest, proto.RegistryResponse]
	errorPipe      chan chan error
	messageChannel chan Request
}

func (client *GRPCClient) handlePrimaryStream(
	channelSize int64,
	doneChannel chan struct{},
	messageChannel chan Request,
) {
	stream, err := client.client.UserPrimaryStream(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	client.handleStream(stream, channelSize, doneChannel, messageChannel)
}

func (client *GRPCClient) handleOverflowStream(
	channelSize int64,
	doneChannel chan struct{},
	messageChannel chan Request,
) {
	stream, err := client.client.UserOverflowStream(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	client.handleStream(stream, channelSize, doneChannel, messageChannel)
}

func (client *GRPCClient) handleStream(
	stream grpc.BidiStreamingClient[proto.RegistryRequest, proto.RegistryResponse],
	channelSize int64,
	doneChannel chan struct{},
	messageChannel chan Request,
) {

	errorPipe := make(chan chan error, channelSize)

	responseHandler := &registryResponseHandler{
		stream:      stream,
		errorPipe:   errorPipe,
		doneChannel: doneChannel,
	}

	go client.receiveResponses(responseHandler)

	requestHandler := &registryRequestHandler{
		stream:         stream,
		errorPipe:      errorPipe,
		messageChannel: messageChannel,
	}

	client.handleRequests(requestHandler)
}

func (client *GRPCClient) receiveResponses(streamHandler *registryResponseHandler) {
	for errChan := range streamHandler.errorPipe {
		message, err := streamHandler.stream.Recv()
		if err == io.EOF {
			log.Error(err)
		}

		if err != nil {
			errChan <- err
		}

		if message.Status != "OK" && message.Status != "ACCEPTED" {
			errChan <- errors.New(message.Status)
		}

		client.currentNumPrimaryRequests.Add(-1)
		client.currentNumOverflowRequests.Add(-1)
		streamHandler.doneChannel <- struct{}{}
		client.serviceWaitgroup.Done()
	}
}

func (client *GRPCClient) handleRequests(streamHandler *registryRequestHandler) {
	for message := range streamHandler.messageChannel {
		client.currentNumPrimaryRequests.Add(1)
		client.currentNumOverflowRequests.Add(1)
		client.serviceWaitgroup.Add(1)

		if err := streamHandler.stream.Send(&proto.RegistryRequest{
			Name:   message.Message.Name,
			Email:  message.Message.Email,
			Entity: message.Message.Entity,
		}); err != nil {
			message.Response <- errors.New("something went wrong")
		}

		streamHandler.errorPipe <- message.Response
	}
}
