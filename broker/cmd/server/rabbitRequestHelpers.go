package server

import (
	messagetypes "broker-server/cmd/message_types"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (app *application) computeValidateRequest(timeout context.Context, messagePayload any, payloadHeader string) error {
	switch payloadHeader {

	case "user:validate":
		return app.rabbitClient.PublishMessage(
			timeout,
			"users",
			"user_validation",
			payloadHeader,
			messagePayload,
		)

	case "forgotten:password":
		return app.rabbitClient.PublishMessage(
			timeout,
			"users",
			"forgotten_password",
			payloadHeader,
			messagePayload,
		)

	case "user:create":
		return app.rabbitClient.PublishMessage(
			timeout,
			"users",
			"create_user",
			payloadHeader,
			messagePayload,
		)

	case "user:update":
		return app.rabbitClient.PublishMessage(
			timeout,
			"users",
			"update_user",
			payloadHeader,
			messagePayload,
		)

	case "user:delete":
		return app.rabbitClient.PublishMessage(
			timeout,
			"users",
			"delete_user",
			payloadHeader,
			messagePayload,
		)

	default:
		return fmt.Errorf("invalid header")
	}
}

func setupPayload(ctx *fiber.Ctx, payloadHeader string) (any, error) {
	switch payloadHeader {

	case "user:validate":
		return newPayload(ctx, new(messagetypes.ValidationNumberMessage))

	case "forgotten:password":
		return newPayload(ctx, new(messagetypes.ValidationNumberMessage))

	case "user:create":
		return newPayload(ctx, new(messagetypes.ChangeUserMessage))

	case "user:update":
		return newPayload(ctx, new(messagetypes.ChangeUserMessage))

	case "user:delete":
		return newPayload(ctx, new(messagetypes.DeleteUserMessage))

	default:
		return nil, fmt.Errorf("invalid payload header")
	}
}

func newPayload(ctx *fiber.Ctx, payload any) (any, error) {
	if err := ctx.BodyParser(payload); err != nil {
		return nil, err
	}

	return payload, nil
}
