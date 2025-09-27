package servercore

import (
	"auth_server/cmd/api/domain/types"
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
		return setVerifySessionCreds(ctx)
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
		verifyCreds, ok := ctx.Locals("verify_creds").(*types.VerifyAccountDTO)
		if !ok {
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		session, err := srv.impl.Verify(verifyCreds)
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
		loginCreds, ok := ctx.Locals("login_creds").(*types.LoginAccountDTO)
		if !ok {
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		session, err := srv.impl.Login(loginCreds)
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
		session, ok := ctx.Locals("verify_session_creds").(*types.VerifySessionDTO)
		if !ok {
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		ctx.Cookie(&fiber.Cookie{
			Name:     "__HOST_VERIFY_",
			Value:    session.Cookie,
			Secure:   true,
			HTTPOnly: true,
			SameSite: "None",
		})

		return ctx.Status(http.StatusCreated).SendString(
			fmt.Sprintf("Hello %s. Please look in your email to get your verify code", session.Name))
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
		changePassCreds, ok := ctx.Locals("change_password_creds").(*types.ChangePasswordDTO)
		if !ok {
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		session, err := srv.impl.ChangePassword(changePassCreds)
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
		refreshCreds, ok := ctx.Locals("refresh_creds").(*types.RefreshSessionDTO)
		if !ok {
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		session, err := srv.impl.RefreshSession(refreshCreds)
		if err != nil {
			ctx.ClearCookie()
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
		refreshCreds, ok := ctx.Locals("refresh_creds").(*types.RefreshSessionDTO)
		if !ok {
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		if err := srv.impl.RemoveSession(refreshCreds); err != nil {
			return ctx.Status(http.StatusForbidden).SendString("invalid credentials")
		}

		ctx.ClearCookie()

		return ctx.SendStatus(http.StatusOK)
	}
}
