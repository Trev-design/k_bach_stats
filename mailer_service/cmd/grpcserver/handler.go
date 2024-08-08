package grpcserver

import (
	"context"
	"fmt"
	"mailerservice/cmd/email"
	"mailerservice/cmd/proto"
)

func (server *ValidationMailerServer) SendValidationEmail(
	ctx context.Context,
	message *proto.EmailRequest,
) (*proto.EmailResponse, error) {

	kind, err := getSubject(message.Kind)
	if err != nil {
		return nil, err
	}

	newMailMSG := email.Message{
		To:      message.Email,
		Subject: kind,
		Payload: &email.ValidationMessage{
			Kind:             message.Kind,
			Email:            message.Email,
			ValidationNumber: message.Validation,
			Name:             message.Name,
		},
	}

	server.Host.Wait.Add(1)
	server.Host.MailerChannel <- newMailMSG

	return &proto.EmailResponse{Result: "OK"}, nil
}

func getSubject(kind string) (string, error) {
	switch kind {
	case "validation":
		return "Validate your KBach account!", nil

	case "forgotten_password":
		return "Validate your identity to reset your password!", nil

	default:
		return "", fmt.Errorf("invalid message type")
	}
}
