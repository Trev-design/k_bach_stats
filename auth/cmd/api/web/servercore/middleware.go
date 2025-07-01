package servercore

import (
	"crypto/subtle"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CSRFAuth(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("__HOST_CSRF_")
	token := ctx.Get("X-CSRF-Token")
	if token == "" || cookie == "" || subtle.ConstantTimeCompare([]byte(cookie), []byte(token)) != 1 {
		return ctx.Status(403).SendString("Forbidden")
	}

	newPayload, err := newCSRFTokenPayload()
	if err != nil {
		return ctx.SendStatus(500)
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "__HOST_CSRF_",
		Value:    newPayload,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "None",
		Expires:  time.Now().Add(2 * time.Hour),
	})

	ctx.Set("X-CSRF-Token", newPayload)

	return ctx.Next()
}
