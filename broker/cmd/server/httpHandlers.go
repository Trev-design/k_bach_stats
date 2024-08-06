package server

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type handleFn func(ctx *fiber.Ctx) error

func (app *application) RabbitRequest(payloadHeader string) handleFn {
	return func(ctx *fiber.Ctx) error {
		messagePayload, err := setupPayload(ctx, payloadHeader)
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := app.computeValidateRequest(
			timeout,
			messagePayload,
			payloadHeader,
		); err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}
