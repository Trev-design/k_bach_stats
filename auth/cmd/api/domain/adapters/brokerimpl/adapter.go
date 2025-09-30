package brokerimpl

type Adapter interface {
	// receives messages for computation
	SendMessage(channelName string, message []byte) error

	// registers and compute long living backround processes
	ComputeBackgroundServices()

	// cleans up the broker on certain scenarions like server shutdown or reconnect the broker
	CloseProducer() error
}
