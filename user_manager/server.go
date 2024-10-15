package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"user_manager/api"
	"user_manager/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbOptions := flag.String("database", "keep", "set the reset flag for the sql database")
	sessionOptions := flag.String("session", "keep", "set th reset flag for the session")
	flag.Parse()

	serverOptions := &api.ServerOptions{Database: *dbOptions, Session: *sessionOptions}

	app, err := serverOptions.Setup()
	if err != nil {
		log.Fatalf("could not start app %v", err)
	}

	go app.ListenForErrors()
	go app.ListenForShutdown()

	go app.Consumer.ComputeMessages("start_session", app.Session, app.ErrorChannel)
	go app.Consumer.ComputeMessages("stop_session", app.Session, app.ErrorChannel)
	go app.Consumer.ComputeMessages("add_user", app.DBIsntance, app.ErrorChannel)
	go app.Consumer.ComputeMessages("remove_user", app.DBIsntance, app.ErrorChannel)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
