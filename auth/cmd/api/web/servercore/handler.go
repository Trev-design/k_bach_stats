package servercore

import (
	"auth_server/cmd/api/domain/types"
	"auth_server/cmd/api/web/webtypes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// @Summary		CSRFToken for user
// @Description	Provides the user a csrf token
// @Tags			core
// @Success		200
// @Failure		401
// @Failure		403
// @Failure		500
// @Router			/csrf [get]
func (server *Server) GetCSRFToken() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token, err := newCSRFTokenPayload()
		if err != nil {
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		ctx.Cookie(&fiber.Cookie{
			Name:     "__HOST_CSRF_",
			Value:    token,
			Secure:   true,
			HTTPOnly: true,
			SameSite: "None",
		})

		ctx.Set("X-CSRF-Token", token)

		return ctx.SendStatus(http.StatusOK)
	}
}

// @Summary		new user
// @Description	Registers a new user
// @Tags			core
// @Param			request	body		webtypes.NewAccountRepresentation	true	"register param"
// @Success		200		{string}	string								"Response Message"
// @Failure		500
// @Failure		403	{string}	string	"failure message"
// @Security		CSRFCookie
// @Security		CSRFToken
// @Router			/register [post]
func (srv *Server) NewUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		newUser := new(webtypes.NewAccountRepresentation)
		payload := ctx.Body()

		if err := json.Unmarshal(payload, newUser); err != nil {
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		sessionPayload, err := srv.impl.Register(toNewAccountDTO(newUser))
		if err != nil {
			return ctx.Status(http.StatusForbidden).SendString("invalid login credentials")
		}

		ctx.Cookie(&fiber.Cookie{
			Name:     "__HOST_VERIFY_",
			Value:    sessionPayload.Cookie,
			Secure:   true,
			HTTPOnly: true,
			SameSite: "None",
		})

		return ctx.Status(http.StatusCreated).
			SendString(
				fmt.Sprintf(
					"Hello %s. Please look in your email to get your verify code",
					sessionPayload.Name))
	}
}

// @Summary		verify users account
// @Description	To unlock the account the user must verify the account with a code
// @Tags			core
// @Param			request		body	webtypes.VerifyRepresentation	true	"the verify code"
//
// @Param			User-Agent	header	string							true	"User-Agent string"
//
// @Success		200
// @Failure		403	{string}	string	"failure message"
// @Security		CSRFCookie
// @Security		CSRFToken
// @Security		VerifyCookie
// @Router			/verify [patch]
func (srv *Server) VerifyAccount() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
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

		session, err := srv.impl.Verify(&types.VerifyAccountDTO{
			Cookie:     cookie,
			VerifyCode: newVerify.Code,
			UserAgent:  userAgent,
			IPAddress:  ip,
		})
		if err != nil {
			return ctx.Status(http.StatusForbidden).SendString("invalid verify credentials")
		}

		setSessionCreds(ctx, session)
		ctx.ClearCookie("__HOST_VERIFY_")

		return ctx.SendStatus(http.StatusAccepted)
	}
}

// @Summary		Login user account
// @Description	route to log the user in required password and email
// @Tags			core
// @Param			request		body	webtypes.LoginRepresentation	true	"login credentials"
// @Param			User-Agent	header	string							true	"User-Agent string"
// @Success		201
// @Failure		403	{string}	string	"failure message"
// @Security		CSRFCookie
// @Security		CSRFToken
// @Router			/login [post]
func (srv *Server) LoginAccount() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		login := new(webtypes.LoginRepresentation)

		payload := ctx.Body()
		if err := json.Unmarshal(payload, login); err != nil {
			return ctx.Status(http.StatusForbidden).SendString("invalid login credentials")
		}

		userAgent := ctx.Get("User-Agent")
		ip := ctx.IP()

		session, err := srv.impl.Login(&types.LoginAccountDTO{
			Email:     login.Email,
			Password:  login.Password,
			UserAgent: userAgent,
			IPAddress: ip,
		})
		if err != nil {
			return ctx.Status(http.StatusForbidden).SendString("invalid login credentials")
		}

		setSessionCreds(ctx, session)

		return ctx.SendStatus(http.StatusAccepted)
	}
}

// @Summary		new password for user
// @Description	If the user of an account has forgotten the password for some reason or just want to change the password this is the route. The user can request a session to change the password
// @Tags			special situations
// @Param			request		body		webtypes.NewPasswordRepresentation	true	"new password session request type"
// @Param			User-Agent	header		string								true	"User-Agent string"
// @Success		200			{string}	string								"Response Message"
// @Failure		403			{string}	string								"failure message"
// @Security		CSRFCookie
// @Security		CSRFToken
// @Router			/new_password [post]
func (srv *Server) NewPassword() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		newPassword := new(webtypes.NewPasswordRepresentation)
		payload := ctx.Body()
		if err := json.Unmarshal(payload, newPassword); err != nil {
			return ctx.Status(http.StatusForbidden).SendString("invalid credentials")
		}

		session, err := srv.impl.NewPassword(newPassword.Email)
		if err != nil {
			return ctx.Status(http.StatusForbidden).SendString("invalid credentials")
		}

		ctx.Cookie(&fiber.Cookie{
			Name:     "__HOST_VERIFY_",
			Value:    session.Cookie,
			Secure:   true,
			HTTPOnly: true,
			SameSite: "None",
		})

		return ctx.Status(http.StatusCreated).
			SendString(
				fmt.Sprintf(
					"Hello %s. Please look in your email to get your verify code",
					session.Name))
	}
}

// @Summary		change password for user
// @Description	If the client has a valid change password session the user can change the password on this route
// @Tags			special situations
// @Param			request		body	webtypes.ChangePasswordRepresentation	true	"required change password credentials"
// @Param			User-Agent	header	string									true	"User-Agent string"
// @Success		201
// @Failure		403	{string}	string	"failure response"
// @Security		CSRFToken
// @Security		CSRFCookie
// @Security		VerifyCookie
// @Router			/change_password [patch]
func (srv *Server) ChangePassword() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
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

		session, err := srv.impl.ChangePassword(&types.ChangePasswordDTO{
			Email:        changePass.Email,
			Password:     changePass.Password,
			Confirmation: changePass.Confirmation,
			VerifyCode:   changePass.VerifyCode,
			UserAgent:    userAgent,
			IPAddress:    ip,
			Cookie:       cookie,
		})
		if err != nil {
			return ctx.Status(http.StatusForbidden).SendString("invalid credentials")
		}

		setSessionCreds(ctx, session)
		ctx.ClearCookie("__HOST_VERIFY_")
		return ctx.SendStatus(http.StatusAccepted)
	}
}

// @Summary		refresh session for user access
// @Description	if the client has a valid refresh session the session can be refreshed on this route
// @Tags			protected
// @Param			User-Agent	header	string	true	"User-Agent string"
// @Success		201
// @Failure		403	{string}	string	"failure response"
// @Security		CSRFToken
// @Security		CSRFCookie
// @Security		RefreshCookie
// @Router			/refresh [post]
func (srv *Server) RefreshSession() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		refresh := ctx.Cookies("__HOST_REFRESH_")
		if refresh == "" {
			return ctx.Status(http.StatusForbidden).SendString("invalid credentials")
		}
		userAgent := ctx.Get("User-Agent")
		ip := ctx.IP()

		session, err := srv.impl.RefreshSession(&types.RefreshSessionDTO{
			Cookie:    refresh,
			UserAgent: userAgent,
			IPAddress: ip,
		})
		if err != nil {
			return ctx.Status(http.StatusForbidden).SendString("invalid credentials")
		}

		setSessionCreds(ctx, session)

		return ctx.SendStatus(http.StatusAccepted)
	}
}

// @Summary		remove session for user
// @Description	some explicit logout functionality to expire the session
// @Tags			protected
// @Param			User-Agent	header	string	true	"User-Agent string"
// @Success		200
// @Failure		403	{string}	string	"failure response"
// @Security		CSRFToken
// @Security		CSRFCookie
// @Security		RefreshCookie
// @Router			/remove [post]
func (srv *Server) RemoveSession() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		refresh := ctx.Cookies("__HOST_REFRESH_")
		if refresh == "" {
			return ctx.Status(http.StatusForbidden).SendString("invalid credentials")
		}
		userAgent := ctx.Get("User-Agent")
		ip := ctx.IP()

		if err := srv.impl.RemoveSession(&types.RemoveSessionDTO{
			Cookie:    refresh,
			UserAgent: userAgent,
			IPAddress: ip,
		}); err != nil {
			return ctx.Status(http.StatusForbidden).SendString("invalid credentials")
		}

		ctx.ClearCookie()

		return ctx.SendStatus(http.StatusOK)
	}
}
