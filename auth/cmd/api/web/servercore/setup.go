package servercore

import (
	"auth_server/cmd/api/domain/adapters/domainimpl"
	_ "auth_server/docs"

	"github.com/gofiber/fiber/v2"
	swagger "github.com/gofiber/swagger"
)

type Server struct {
	app  *fiber.App
	impl domainimpl.Adapter
}

func NewServer(impl domainimpl.Adapter) *Server {
	app := fiber.New()
	server := new(Server)
	server.app = app
	server.impl = impl

	server.registerRoutes()

	return server
}

func (srv *Server) StartAndListen() error {
	return srv.app.Listen(":4000")
}

func (srv *Server) registerRoutes() {
	group := srv.app.Group("/api/v1")
	group.Get("/swagger", swagger.HandlerDefault)
	group.Get("/csrf", srv.GetCSRFToken())
	group.Post("/register", srv.registerUser(), srv.NewUser())
	group.Patch("/verify", fetchVerifyCredentials, srv.VerifyAccount())
	group.Post("/login", fetchLoginCredentials, srv.LoginAccount())
	group.Post("/new_password", srv.fetchNewPasswordCredentials(), srv.NewPassword())
	group.Patch("/change_password", fetchChangePasswordCredentials, srv.ChangePassword())
	group.Post("/refresh", fetchRefreshSessionCreds, srv.RefreshSession())
	group.Delete("/logout", fetchRefreshSessionCreds, srv.RemoveSession())
}
