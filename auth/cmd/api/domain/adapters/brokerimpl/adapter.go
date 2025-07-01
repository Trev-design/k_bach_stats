package brokerimpl

type Adapter interface {
	SendMessage(channelName string, message []byte) error
	ComputeBackgroundServices()
	CloseProducer() error
}
