package servercore

import (
	"auth_server/cmd/api/domain/types"
	"crypto/rand"
	"encoding/base64"
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
		Secure:   true,
		HTTPOnly: true,
		SameSite: "None",
		Expires:  time.Now().Add(2 * time.Hour),
	})

	ctx.Set("Authorization", session.Access)
}
