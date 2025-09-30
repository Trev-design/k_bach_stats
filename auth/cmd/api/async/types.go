package async

type Fn func(chan Output, Input)

type Payload interface {
	JSONRepresentation() []byte
}

type Input interface {
	JSON() []byte
}

type Output interface {
	Error() error
	Data() Payload
}
