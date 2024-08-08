package grpcclient

import (
	"listener/cmd/proto"

	"google.golang.org/grpc"
)

type GRPCConnection struct {
	loggerConnection          *grpc.ClientConn
	emailValidationConnection *grpc.ClientConn
}

type GRPCClientStructure struct {
	LoggerClient proto.LogMessageServiceClient

	EmailValidationClient proto.ValidationServiceClient
	//userModifyClient      *proto.UserModifyServiceClient
	connection *GRPCConnection
}

func (structure *GRPCClientStructure) CloseLoggerConnection() error {
	return structure.connection.loggerConnection.Close()
}

func (structure *GRPCClientStructure) CloseValidationEmailConnection() error {
	return structure.connection.emailValidationConnection.Close()
}
