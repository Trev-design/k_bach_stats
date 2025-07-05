package main

import (
	"auth_server/cmd/api/broker/channel"
	"auth_server/cmd/api/broker/producer"
	"auth_server/cmd/api/db/dbcore"
	"auth_server/cmd/api/impl"
	"auth_server/cmd/api/jwt/jwtcore"
	"auth_server/cmd/api/temporary/session/sessioncore"
	"auth_server/cmd/api/tlsconf"
	"auth_server/cmd/api/web/servercore"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
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
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(wd)

	if err := godotenv.Load("./.env"); err != nil {
		log.Fatal(err)
	}

	implementation, err := impl.NewImplBuilder().
		DatabaseSetup(dbSetup(wd)).
		SessionSetup(sessionSetup(wd)).
		BrokerSetup(brokerSetup(wd)).
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

func dbSetup(wd string) *dbcore.DatabaseBuilder {
	return dbcore.NewDatabaseBuilder().
		User(os.Getenv("POSTGRES_DB_USER")).
		Password(os.Getenv("POSTGRES_DB_PASSWORD")).
		Port(os.Getenv("POSTGRES_DB_PORT")).
		Host(os.Getenv("POSTGRES_DB_HOST")).
		DBName(os.Getenv("POSTGRES_DB_NAME")).
		WithTLS(
			tlsconf.NewTLSBuilder().
				CertPath(fmt.Sprintf("./certs/%s", os.Getenv("POSTGRES_DB_CERT"))).
				KeyPath(fmt.Sprintf("./certs/%s", os.Getenv("POSTGRES_DB_KEY"))).
				CACertPath(fmt.Sprintf("./certs/%s", os.Getenv("POSTGRES_DB_CA_CERT"))))
}

func sessionSetup(wd string) *sessioncore.SessionBuilder {
	return sessioncore.NewSessionBuilder().
		Host(os.Getenv("REDIS_SESSION_HOST")).
		Port(os.Getenv("REDIS_SESSION_PORT")).
		Password(os.Getenv("REDIS_SESSION_PASSWORD")).
		IntevalDuration(2 * time.Hour).
		WithTLS(
			tlsconf.NewTLSBuilder().
				CertPath(filepath.Join(wd, "certs", os.Getenv("REDIS_SESSION_CERT"))).
				KeyPath(filepath.Join(wd, "certs", os.Getenv("REDIS_SESSION_KEY"))))
}

func brokerSetup(wd string) *producer.RMQProducerBuilder {
	return producer.NewProducer().
		User(os.Getenv("RABBIT_USER")).
		Password(os.Getenv("RABBIT_PASSWORD")).
		Host(os.Getenv("RABBIT_HOST")).
		Port(os.Getenv("RABBIT_PORT")).
		VirtualHost(os.Getenv("RABBIT_V_HOST")).
		WithTLS(
			tlsconf.NewTLSBuilder().
				CertPath(filepath.Join(wd, "certs", os.Getenv("RABBIT_CERT"))).
				KeyPath(filepath.Join(wd, "certs", os.Getenv("RABBIT_KEY")))).
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
