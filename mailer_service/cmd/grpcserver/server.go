package grpcserver

import (
	"mailerservice/cmd/email"
	"mailerservice/cmd/proto"
	"net"

	"google.golang.org/grpc"
)

type ValidationMailerServer struct {
	proto.ValidationServiceServer
	Host *email.MailHost
}

func StartAndListen() error {
	listener, err := net.Listen("tcp", ":5297")
	if err != nil {
		return err
	}

	host := email.NewValidationMailer()
	server := grpc.NewServer()
	validationServer := &ValidationMailerServer{Host: host}

	proto.RegisterValidationServiceServer(server, validationServer)

	go host.ListenForEmails()

	return server.Serve(listener)
}
