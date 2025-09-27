package servercore

import (
	"auth_server/cmd/api/domain/types"
	"auth_server/cmd/api/web/webtypes"
	"crypto/subtle"
	"encoding/json"
	"net/http"
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

func (server *Server) registerUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		newUser := new(webtypes.NewAccountRepresentation)
		payload := ctx.Body()

		if err := json.Unmarshal(payload, newUser); err != nil {
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		sessionPayload, err := server.impl.Register(toNewAccountDTO(newUser))
		if err != nil {
			return ctx.Status(http.StatusForbidden).SendString("invalid login credentials")
		}

		ctx.Locals("verify_session", sessionPayload)

		return ctx.Next()
	}
}

func fetchVerifyCredentials(ctx *fiber.Ctx) error {
	newVerify := new(webtypes.VerifyRepresentation)
	payload := ctx.Body()

	if err := json.Unmarshal(payload, newVerify); err != nil {
		return ctx.Status(http.StatusForbidden).SendString("invalid verify credentials")
	}

	cookie := ctx.Cookies("__HOST_VERIFY_")
	if cookie == "" {
		return ctx.Status(http.StatusForbidden).SendString("invalid verify credentials")
	}

	userAgent := ctx.Get("User-Agent")
	ip := ctx.IP()

	ctx.Locals("verify_creds", &types.VerifyAccountDTO{
		Cookie:     cookie,
		VerifyCode: newVerify.Code,
		UserAgent:  userAgent,
		IPAddress:  ip,
	})

	return ctx.Next()
}

func fetchLoginCredentials(ctx *fiber.Ctx) error {
	login := new(webtypes.LoginRepresentation)

	payload := ctx.Body()
	if err := json.Unmarshal(payload, login); err != nil {
		return ctx.Status(http.StatusForbidden).SendString("invalid login credentials")
	}

	userAgent := ctx.Get("User-Agent")
	ip := ctx.IP()

	ctx.Locals(
		"login_creds",
		&types.LoginAccountDTO{
			Email:     login.Email,
			Password:  login.Password,
			UserAgent: userAgent,
			IPAddress: ip,
		})

	return ctx.Next()
}

func (server *Server) fetchNewPasswordCredentials() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		newPassword := new(webtypes.NewPasswordRepresentation)
		payload := ctx.Body()
		if err := json.Unmarshal(payload, newPassword); err != nil {
			return ctx.Status(http.StatusForbidden).SendString("invalid credentials")
		}

		session, err := server.impl.NewPassword(newPassword.Email)
		if err != nil {
			return ctx.Status(http.StatusForbidden).SendString("invalid credentials")
		}

		ctx.Locals("verify_session_creds", session)

		return ctx.Next()
	}
}

func fetchChangePasswordCredentials(ctx *fiber.Ctx) error {
	changePass := new(webtypes.ChangePasswordRepresentation)
	payload := ctx.Body()
	if err := json.Unmarshal(payload, changePass); err != nil {
		return ctx.Status(http.StatusForbidden).SendString("invalid credentials")
	}

	userAgent := ctx.Get("User-Agent")
	ip := ctx.IP()
	cookie := ctx.Cookies("__HOST_VERIFY_")
	if cookie == "" {
		return ctx.Status(http.StatusForbidden).SendString("invalid credentials")
	}

	ctx.Locals("change_password_creds", &types.ChangePasswordDTO{
		Email:        changePass.Email,
		Password:     changePass.Password,
		Confirmation: changePass.Confirmation,
		VerifyCode:   changePass.VerifyCode,
		UserAgent:    userAgent,
		IPAddress:    ip,
		Cookie:       cookie,
	})

	return ctx.Next()
}

func fetchRefreshSessionCreds(ctx *fiber.Ctx) error {
	refresh := ctx.Cookies("__HOST_REFRESH_")
	if refresh == "" {
		return ctx.Status(http.StatusForbidden).SendString("invalid credentials")
	}
	userAgent := ctx.Get("User-Agent")
	access := ctx.Get("Authorization")
	ip := ctx.IP()

	ctx.Locals("refresh_creds", &types.RefreshSessionDTO{
		JWT:       access,
		Cookie:    refresh,
		UserAgent: userAgent,
		IPAddress: ip,
	})

	return ctx.Next()
}
