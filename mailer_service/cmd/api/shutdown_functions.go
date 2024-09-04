package api

import (
	"os"
	"os/signal"
	"syscall"
)

func (application *app) listenForShutDown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	application.shutdown()
}

func (application *app) shutdown() {
	application.rmqClient.CloseChannel()
	application.rmqClient.CloseConnection()

	application.mailHost.Wait.Wait()
	application.mailHost.DoneChannel <- true

	close(application.mailerChannel)
	close(application.mailHost.ErrorChannel)
	close(application.mailHost.DoneChannel)

	application.wait.Done()
}
