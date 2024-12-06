package main

import (
	"log"
	"net/http"
	"os"
	"user_manager/app"
	"user_manager/graph"
	"user_manager/internal/application"
	"user_manager/internal/plugins/database/sqldb"
	"user_manager/internal/plugins/listener/rabbit"
	redissession "user_manager/internal/plugins/session/redis_session.go"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	// TODO show environment variable management
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// create database
	dba, err := sqldb.NewDatabaseAdapter("reset")
	if err != nil {
		panic(err)
	}

	log.Println("finshed database")

	// create session store
	sessiona, err := redissession.NewSessionAdapter()
	if err != nil {
		panic(err)
	}

	// api infrastructure
	apia := application.NewAPIServiceAdapter(dba)

	// rabbitMQ listener
	// TODO safe password storage
	la, err :=
		rabbit.NewListenerAdapterBuilder().
			Connection("IAmTheUser", "ThisIsMyPassword", "localhost", "kbach", 5672).
			Channel("start_user_session", "session", "send_session_credentials", "add_session_consumer").
			Channel("stop_user_session", "session", "remove_session", "remove_session_consumer").
			Channel("add_account", "account", "add_account_request", "add_user_consumer").
			Channel("delete_accoubt", "account", "delete_account_request", "remove_user_consumer").
			SessionStore(sessiona).
			UserStore(dba).
			Build()
	if err != nil {
		panic(err)
	}

	// listener service
	// TODO builder pattern
	lsa := application.NewListenerServiceAdapter(la)

	if err := lsa.Start(); err != nil {
		panic(err)
	}

	// app config
	app, err :=
		app.NewApplicationBuilder().
			ApiService(apia).
			ListenerService(lsa).
			Build()
	if err != nil {
		panic(err)
	}

	go app.ListenForShutdown()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{ApiService: apia}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
