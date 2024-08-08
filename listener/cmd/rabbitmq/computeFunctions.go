package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"listener/cmd/grpcclient"
	"listener/cmd/messagetypes"
	"listener/cmd/proto"
	"log"
	"time"

	"gorm.io/gorm"
)

func computeEmailMessage(
	message *messagetypes.Message,
	clients *grpcclient.GRPCClientStructure,
	db *gorm.DB,
) error {
	log.Printf("the email message is: %v\n", message)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jsonPayload, err := json.Marshal(message.Payload)
	if err != nil {
		return fmt.Errorf("something went wrong")
	}

	body := new(messagetypes.ValidationNumberMessage)
	if err := json.Unmarshal(jsonPayload, body); err != nil {
		return fmt.Errorf("something went wrong")
	}

	messageBody := &proto.EmailRequest{
		Name:       body.UserName,
		Email:      body.Email,
		Validation: body.ValidationNumber,
		Kind:       getEmailValidationKind(message.Type),
	}

	log.Printf("send %v\n", messageBody)

	_, err = clients.EmailValidationClient.SendValidationEmail(ctx, messageBody)
	if err != nil {
		return err
	}

	return nil
}

func getEmailValidationKind(kind string) string {
	switch kind {
	case "user:validate":
		return "validation"

	case "user:forgotten_password":
		return "forgotten_password"

	default:
		return "invalid"
	}
}

func computeModifyUserMessage(
	message *messagetypes.Message,
	client *grpcclient.GRPCClientStructure,
	db *gorm.DB) error {
	log.Printf("the user message is: %v\n", message)
	return nil
}
