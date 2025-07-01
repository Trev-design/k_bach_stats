package main

import (
	"os"
	"os/signal"
	"syscall"
)

func (srv *app) listenForShutDown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	srv.shutdown()
}

func (srv *app) shutdown() {
	srv.rmqSrv.CloseRabbit()
	srv.mailSrv.CloseMailhost()
}
