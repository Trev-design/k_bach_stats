package keyorchestrator

import (
	"auth_server/cmd/api/grpc/keyorchestrator/proto/proto"
	"context"
	"errors"
	"io"
	"log"
	"sync"

	"google.golang.org/grpc"
)

func newGRPCStream(responseChannel chan *GetKeyResponse) *grpcStream {
	return &grpcStream{
		serviceWaitgroup:    &sync.WaitGroup{},
		requestChannel:      make(chan *requestType),
		doneChannel:         make(chan struct{}),
		getResponseChannels: responseChannel,
	}
}

func (clientStream *grpcStream) sendSubscribe(sub *GetKeySubscriptionType) {
	clientStream.sendRequest(&sub.requestType)
}

func (clientStream *grpcStream) sendRetry(retry *RetryType) {
	clientStream.sendRequest(&retry.requestType)
}

func (clientStream *grpcStream) sendRequest(request *requestType) {
	clientStream.requestChannel <- request
}

func (clientStream *grpcStream) handle(client proto.KeyOrchestratorServiceClient) {
	stream, err := client.KeyPipe(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	clientStream.handleStream(stream)
}

func (clientStream *grpcStream) handleStream(stream grpc.BidiStreamingClient[proto.GetKeyRequest, proto.KeyResponse]) {
	go clientStream.handleResponse(stream)

	for {
		select {
		case request := <-clientStream.requestChannel:
			clientStream.handleRequest(request, stream)

		case <-clientStream.doneChannel:
			return
		}
	}
}

func (clientStream *grpcStream) handleResponse(stream grpc.BidiStreamingClient[proto.GetKeyRequest, proto.KeyResponse]) {
	for {
		message, err := stream.Recv()
		if err == io.EOF {
			log.Println(err)
			break
		}

		if err != nil {
			log.Println(err)
		}

		var responseError error = nil

		if message.Message != "ACCEPTED" {
			responseError = errors.New(message.Message)
		}

		clientStream.getResponseChannels <- &GetKeyResponse{
			Key:     message.Key,
			Message: responseError,
		}
	}
}

func (clientStream *grpcStream) handleRequest(
	request *requestType,
	stream grpc.BidiStreamingClient[proto.GetKeyRequest, proto.KeyResponse],
) {
	if err := stream.Send(&proto.GetKeyRequest{
		Id:    request.Id,
		Topic: request.Topic,
	}); err != nil {
		log.Println(err.Error())
	}
}
