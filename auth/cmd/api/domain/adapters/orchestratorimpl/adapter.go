package orchestratorimpl

import "auth_server/cmd/api/grpc/keyorchestrator"

type Adapter interface {
	OpenPipe(message *keyorchestrator.GetKeySubscriptionType)
}
