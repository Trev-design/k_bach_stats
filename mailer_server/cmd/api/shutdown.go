package api

import (
	"os"
	"os/signal"
	"syscall"
)

func (server *app) listenForShutDown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	server.shutDown()
	os.Exit(0)
}

func (server *app) shutDown() {
	server.wait.Wait()
}
