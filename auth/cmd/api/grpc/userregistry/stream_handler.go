package userregistry

import (
	"auth_server/cmd/api/grpc/userregistry/proto"
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"

	"github.com/gofiber/fiber/v2/log"
	"google.golang.org/grpc"
)

type registryResponseHandler struct {
	stream             grpc.BidiStreamingClient[proto.RegistryRequest, proto.RegistryResponse]
	errorPipe          chan chan error
	currentNumRequests *atomic.Int64
	waitgroup          *sync.WaitGroup
}

type registryRequestHandler struct {
	stream             grpc.BidiStreamingClient[proto.RegistryRequest, proto.RegistryResponse]
	errorPipe          chan chan error
	messageChannel     chan Request
	currentNumRequests *atomic.Int64
	waitgroup          *sync.WaitGroup
}

type registryStreamHandler struct {
	stream             grpc.BidiStreamingClient[proto.RegistryRequest, proto.RegistryResponse]
	channelSize        int64
	messageChannel     chan Request
	currentNumRequests *atomic.Int64
	waitgroup          *sync.WaitGroup
}

func (client *GRPCClient) handlePrimaryStream() {
	stream, err := client.client.UserPrimaryStream(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	streamHandler := &registryStreamHandler{
		stream:             stream,
		channelSize:        client.maxNumPrimaryRequests,
		messageChannel:     client.messagePrimaryChannel,
		currentNumRequests: &client.currentNumPrimaryRequests,
		waitgroup:          client.serviceWaitgroup,
	}

	streamHandler.handleStream()
}

func (client *GRPCClient) handleOverflowStream() {
	stream, err := client.client.UserOverflowStream(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	streamHandler := &registryStreamHandler{
		stream:             stream,
		channelSize:        client.maxNumOverflowRequests,
		messageChannel:     client.messageOverflowChannel,
		currentNumRequests: &client.currentNumOverflowRequests,
		waitgroup:          client.serviceWaitgroup,
	}

	streamHandler.handleStream()
}

func (streamHandler *registryStreamHandler) handleStream() {

	errorPipe := make(chan chan error, streamHandler.channelSize)

	responseHandler := &registryResponseHandler{
		stream:             streamHandler.stream,
		errorPipe:          errorPipe,
		currentNumRequests: streamHandler.currentNumRequests,
		waitgroup:          streamHandler.waitgroup,
	}

	go responseHandler.receiveResponses()

	requestHandler := &registryRequestHandler{
		stream:             streamHandler.stream,
		errorPipe:          errorPipe,
		messageChannel:     streamHandler.messageChannel,
		currentNumRequests: streamHandler.currentNumRequests,
		waitgroup:          streamHandler.waitgroup,
	}

	requestHandler.handleRequests()
}

func (streamHandler *registryResponseHandler) receiveResponses() {
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

		streamHandler.currentNumRequests.Add(-1)
		streamHandler.waitgroup.Done()
	}
}

func (streamHandler *registryRequestHandler) handleRequests() {
	for message := range streamHandler.messageChannel {
		streamHandler.currentNumRequests.Add(1)
		streamHandler.waitgroup.Add(1)

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
