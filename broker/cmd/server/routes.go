package server

func (app *application) routes() {
	app.httpServer.Post("/user/validate", app.RabbitRequest("user:validate"))
	app.httpServer.Post("/user/forgotten_password", app.RabbitRequest("forgotten:password"))
	app.httpServer.Post("/user/create", app.RabbitRequest("user:create"))
	app.httpServer.Post("/user/update", app.RabbitRequest("user:update"))
	app.httpServer.Post("/user/delete", app.RabbitRequest("user:delete"))
}
