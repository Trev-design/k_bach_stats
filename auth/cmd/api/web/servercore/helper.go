package servercore

import (
	"auth_server/cmd/api/domain/types"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func newCSRFTokenPayload() (string, error) {
	randomBytes := make([]byte, 48)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	tokenPayload := base64.RawURLEncoding.EncodeToString(randomBytes)

	return tokenPayload, nil
}

func setSessionCreds(ctx *fiber.Ctx, session *types.NewAccountSessionDTO) {
	ctx.Cookie(&fiber.Cookie{
		Name:     "__HOST_REFRESH_",
		Value:    session.Refresh,
		Secure:   false,
		HTTPOnly: true,
		SameSite: "Strict",
		Expires:  time.Now().Add(2 * time.Hour),
	})

	ctx.Set("Authorization", session.Access)
}

func setVerifySessionCreds(ctx *fiber.Ctx) error {
	session, ok := ctx.Locals("verify_session").(*types.VerifySessionDTO)
	if !ok {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "__HOST_VERIFY_",
		Value:    session.Cookie,
		Secure:   false,
		HTTPOnly: true,
		SameSite: "Strict",
		Expires:  time.Now().Add(time.Hour),
	})

	return ctx.Status(http.StatusCreated).
		SendString(fmt.Sprintf(
			"Hello %s. Please look in your email to get your verify code",
			session.Name))
}
