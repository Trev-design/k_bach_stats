package keyorchestrator

import "auth_server/cmd/api/grpc/keyorchestrator/proto/proto"

type GetKeyResponse struct {
	Key     string
	Message error
}

type RetryType struct {
	Id      string
	Request *proto.GetKeyRequest
}

type GetKeySubscriptionType struct {
	Id           string
	Request      *proto.GetKeyRequest
	ResponsePipe chan GetKeyResponse
}
