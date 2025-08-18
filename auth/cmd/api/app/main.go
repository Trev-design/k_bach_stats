package main

import (
	"auth_server/cmd/api/broker/channel"
	"auth_server/cmd/api/broker/producer"
	"auth_server/cmd/api/db/dbcore"
	"auth_server/cmd/api/impl"
	"auth_server/cmd/api/jwt/jwtcore"
	"auth_server/cmd/api/temporary/session/sessioncore"
	"auth_server/cmd/api/web/servercore"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
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
	implementation, err := impl.NewImplBuilder().
		DatabaseSetup(dbSetup()).
		SessionSetup(sessionSetup()).
		BrokerSetup(brokerSetup()).
		JWTSetup(
			jwtcore.NewJWTServiceBuilder().
				Identifier(uuid.New()).
				Interval(2 * time.Hour),
		).Build()
	if err != nil {
		log.Fatalf("line 52: %v", err)
	}

	server := servercore.NewServer(implementation)

	implementation.StartBackgroundServices()
	go listenForShutDown(implementation)

	log.Fatal(server.StartAndListen())
}

func dbSetup() *dbcore.DatabaseBuilder {

	return dbcore.NewDatabaseBuilder().
		User(os.Getenv("POSTGRES_USER")).
		Password(os.Getenv("POSTGRES_PASSWORD")).
		Port("5432").
		Host("postgres").
		DBName(os.Getenv("POSTGRES_DATABASE_NAME"))
	//WithTLS(
	//	tlsconf.NewTLSBuilder().
	//		CertPath(fmt.Sprintf("./certs/%s", os.Getenv("POSTGRES_DB_CERT"))).
	//		KeyPath(fmt.Sprintf("./certs/%s", os.Getenv("POSTGRES_DB_KEY"))).
	//		CACertPath(fmt.Sprintf("./certs/%s", os.Getenv("POSTGRES_DB_CA_CERT"))))
}

func sessionSetup() *sessioncore.SessionBuilder {
	return sessioncore.NewSessionBuilder().
		Host("redis").
		Port("6379").
		Password(os.Getenv("REDIS_PASSWORD")).
		IntevalDuration(2 * time.Hour)
	//WithTLS(
	//	tlsconf.NewTLSBuilder().
	//		CertPath(filepath.Join(wd, "certs", os.Getenv("REDIS_SESSION_CERT"))).
	//		KeyPath(filepath.Join(wd, "certs", os.Getenv("REDIS_SESSION_KEY"))))
}

func brokerSetup() *producer.RMQProducerBuilder {
	return producer.NewProducer().
		User(os.Getenv("RABBIT_USER")).
		Password(os.Getenv("RABBIT_PASSWORD")).
		Host("rabbitmq").
		Port("5672").
		VirtualHost(os.Getenv("RABBIT_VIRTUAL_HOST")).
		//WithTLS(
		//	tlsconf.NewTLSBuilder().
		//		CertPath(filepath.Join(wd, "certs", os.Getenv("RABBIT_CERT"))).
		//		KeyPath(filepath.Join(wd, "certs", os.Getenv("RABBIT_KEY")))).
		WithChannel(
			"send_verify_email",
			channel.NewPipeBuilder().
				Exchange("mails").
				Kind("direct").
				Queue("verify_mails").
				RoutingKey("send_verify_mail"))
}

func listenForShutDown(impl *impl.Impl) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := impl.CloseServices(); err != nil {
		log.Fatal(err)
	}
}
