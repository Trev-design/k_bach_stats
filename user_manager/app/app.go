package app

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"user_manager/internal/core"
)

type Application struct {
	listenerService core.ListenerService
	apiService      core.ApiService
}

type applicationBulider struct {
	application Application
}

func NewApplicationBuilder() *applicationBulider {
	return &applicationBulider{
		application: Application{},
	}
}

func (builder *applicationBulider) ListenerService(ls core.ListenerService) *applicationBulider {
	builder.application.listenerService = ls
	return builder
}

func (builder *applicationBulider) ApiService(as core.ApiService) *applicationBulider {
	builder.application.apiService = as
	return builder
}

func (builder *applicationBulider) Build() (*Application, error) {
	if builder.application.apiService == nil || builder.application.listenerService == nil {
		return nil, errors.New("necessary services missed")
	}

	return &builder.application, nil
}

func (app *Application) StartListenerService() error {
	return app.listenerService.Start()
}

func (app *Application) ListenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutdown()
}

func (app *Application) shutdown() {
	if err := app.apiService.Close(); err != nil {
		panic("something went wrong")
	}

	if err := app.listenerService.ShutDown(); err != nil {
		panic("something went wrong")
	}

	os.Exit(0)
}
