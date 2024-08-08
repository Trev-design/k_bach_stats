package grpcclient

import (
	"listener/cmd/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SetupClientStructure() (*GRPCClientStructure, error) {
	clients := &GRPCClientStructure{connection: &GRPCConnection{}}

	err := clients.setupLoggerClient()
	if err != nil {
		return nil, err
	}

	err = clients.setupvalidationEmailClient()
	if err != nil {
		return nil, err
	}

	return clients, nil
}

func (clients *GRPCClientStructure) setupLoggerClient() error {
	conn, err := grpc.NewClient("http://localhost:5298", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	client := proto.NewLogMessageServiceClient(conn)
	clients.LoggerClient = client
	clients.connection.loggerConnection = conn

	return nil
}

func (clients *GRPCClientStructure) setupvalidationEmailClient() error {
	conn, err := grpc.NewClient("localhost:5297", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	client := proto.NewValidationServiceClient(conn)
	clients.EmailValidationClient = client
	clients.connection.emailValidationConnection = conn

	return nil
}
