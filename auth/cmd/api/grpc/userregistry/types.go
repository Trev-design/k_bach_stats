package userregistry

type Message struct {
	Name   string
	Email  string
	Entity string
}

type Request struct {
	Message  Message
	Response chan error
}
