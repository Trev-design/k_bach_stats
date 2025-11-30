package keyorchestrator

type GetKeyResponse struct {
	Key     string
	Message error
}

type requestType struct {
	Id    string
	Topic string
}

type RetryType struct {
	requestType
}

type GetKeySubscriptionType struct {
	requestType
	ResponsePipe chan *GetKeyResponse
}
