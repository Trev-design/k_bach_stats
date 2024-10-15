package api

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func (app *App) ListenForErrors() {
	for {
		select {
		case err := <-app.ErrorChannel:
			log.Println(err.Error())

		case <-app.DoneChannel:
			return
		}
	}
}

func (app *App) ListenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutdown()
}

func (app *App) shutdown() {
	app.Consumer.Close()
	app.DoneChannel <- true
	app.Session.Close()
	app.DBIsntance.Close()
	close(app.ErrorChannel)
	close(app.DoneChannel)
	os.Exit(0)
}
