package main

import (
	"auth_server/cmd/api/broker/channel"
	"auth_server/cmd/api/broker/producer"
	"auth_server/cmd/api/db/dbcore"
	"auth_server/cmd/api/impl"
	"auth_server/cmd/api/temporary/session/sessioncore"
	"auth_server/cmd/api/tlsconf"
	"auth_server/cmd/api/web/servercore"
	"fmt"
	"log"
	"os"
	"time"
)

//	@title						API Documentation
//	@version					1
//	@description				Documentation for an auth api which registers users and authorize them with some session credentials
//	@contact.name				Gerrit Flick
//	@host						localhost:4000
//	@BasePath					/api/v1
//	@securityDefinitions.apikey	VerifyCookie
//	@in							cookie
//	@name						__HOST_VERIFY_

//	@securityDefinitions.apikey	CSRFCookie
//	@in							cookie
//	@name						__HOST_CSRF_

//	@securityDefinitions.apikey	CSRFToken
//	@in							header
//	@name						X-CSRF-Token

// @securityDefinitions.apikey	RefreshCookie
// @in							cookie
// @name						__HOST_REFRESH_
func main() {
	if err := tlsconf.GenerateCertPool(os.Getenv("CA_CERT_PATH")); err != nil {
		log.Fatal(err)
	}

	implementation, err := impl.NewImplBuilder().
		DatabaseSetup(dbSetup()).
		SessionSetup(sessionSetup()).
		BrokerSetup(brokerSetup()).
		Build()
	if err != nil {
		log.Fatal(err)
	}

	server := servercore.NewServer(implementation)

	implementation.StartBackgroundServices()
	go listenForShutDown(implementation)

	log.Fatal(server.StartAndListen())
}

func dbSetup() *dbcore.DatabaseBuilder {
	return dbcore.NewDatabaseBuilder().
		User(os.Getenv("POSGRES_DB_USER")).
		Password(os.Getenv("POSTGRES_DB_PASSWORD")).
		Port(os.Getenv("POSTGRES_DB_PORT")).
		Host(os.Getenv("POSTGRES_DB_HOST")).
		DBName(os.Getenv("POSTGRES_DB_NAME")).
		WithTLS(
			tlsconf.NewTLSBuilder().
				CertPath(fmt.Sprintf(".%s", os.Getenv("POSTGRES_DB_CERT_PATH"))).
				KeyPath(fmt.Sprintf(".%s", os.Getenv("POSTGRES_DB_KEY_PATH"))))
}

func sessionSetup() *sessioncore.SessionBuilder {
	return sessioncore.NewSessionBuilder().
		Host(os.Getenv("REDIS_SESSION_HOST")).
		Port(os.Getenv("REDIS_SESSION_PORT")).
		Password(os.Getenv("REDIS_SESSION_PASSWORD")).
		IntevalDuration(2 * time.Hour).
		WithTLS(
			tlsconf.NewTLSBuilder().
				CertPath(fmt.Sprintf(".%s", os.Getenv("REDIS_SESSION_CERT_PATH"))).
				KeyPath(fmt.Sprintf(".%s", os.Getenv("REDIS_SESSION_KEY_PATH"))))
}

func brokerSetup() *producer.RMQProducerBuilder {
	return producer.NewProducer().
		User(os.Getenv("RABBIT_USER")).
		Password(os.Getenv("RABBIT_PASSWORD")).
		Host(os.Getenv("RABBIT_HOST")).
		Port(os.Getenv("RABBIT_PORT")).
		VirtualHost(os.Getenv("RABBIT_V_HOST")).
		WithTLS(
			tlsconf.NewTLSBuilder().
				CertPath(fmt.Sprintf(".%s", os.Getenv("RABBIT_CERT_PATH"))).
				KeyPath(fmt.Sprintf(".%s", os.Getenv("RABBIT_KEY_PATH")))).
		WithChannel(
			"send_verify_email",
			channel.NewPipeBuilder().
				Exchange("mails").
				Kind("direct").
				Queue("verify_mails").
				RoutingKey("send_verify_mail"))
}

func listenForShutDown(impl *impl.Impl) {

}
