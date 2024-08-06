package rabbitmq

import (
	"listener/cmd/messagetypes"
	"log"
	//"gorm.io/gorm"
)

func computeEmailMessage(message *messagetypes.Message) error {
	log.Printf("the email message is: %v\n", message)
	return nil
}

func computeModifyUserMessage(message *messagetypes.Message) error {
	log.Printf("the user message is: %v\n", message)
	return nil
}

func computeErrorLoggingMessage(message *messagetypes.Message) error {
	log.Printf("th error logging message is: %v\n", message)
	return nil
}
