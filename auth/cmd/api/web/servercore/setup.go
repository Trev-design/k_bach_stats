package servercore

import (
	"auth_server/cmd/api/domain/adapters/domainimpl"
	_ "auth_server/docs"
	"os"

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
	return srv.app.ListenTLS(":4000", os.Getenv("AUTH_API_CERT_PATH"), os.Getenv("AUTH_API_KEY_PATH"))
}

func (srv *Server) registerRoutes() {
	group := srv.app.Group("/api/v1")
	group.Get("/swagger", swagger.HandlerDefault)
	group.Get("/csrf", srv.GetCSRFToken())
	group.Post("/register", CSRFAuth, srv.NewUser())
	group.Patch("/verify", CSRFAuth, srv.VerifyAccount())
	group.Post("/login", CSRFAuth, srv.LoginAccount())
	group.Post("/new_password", CSRFAuth, srv.NewPassword())
	group.Patch("/change_password", CSRFAuth, srv.ChangePassword())
}
