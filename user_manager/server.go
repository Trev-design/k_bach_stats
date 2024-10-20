package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"user_manager/api"
	"user_manager/database"
	"user_manager/graph"
	"user_manager/middleware"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/google/uuid"
)

const defaultPort = "4004"

func getUserCredentials(storeHandler database.StoreHandler) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if entity := request.Context().Value(api.ContextKey("entity")); entity != nil {
			entityString, ok := entity.(string)
			if !ok {
				log.Println("invalid type")
				writer.WriteHeader(http.StatusForbidden)
				return
			}

			log.Println("have a valid context")

			userID, err := storeHandler.InitialCredentials(entityString)
			if err != nil {
				log.Println("invalid token")
				writer.WriteHeader(http.StatusForbidden)
				return
			}

			log.Println("have valid user entity")

			userGuid, err := uuid.FromBytes([]byte(userID))
			if err != nil {
				log.Println("invalid id")
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Println("have valid user id")

			writer.Header().Add("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			err = json.NewEncoder(writer).Encode(struct {
				ID string `json:"id"`
			}{
				ID: userGuid.String(),
			})

			if err != nil {
				log.Println("something went wrong")
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Println("finished successfully")

			return
		}

		log.Println("something went wrong with the context")
		writer.WriteHeader(http.StatusInternalServerError)
	})
}

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

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Database: app.DBIsntance}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", middleware.EnableCORS(middleware.Auth(app.Session, srv)))
	http.Handle("/initial", middleware.EnableCORS(middleware.InitialAuth(app.Session, getUserCredentials(app.DBIsntance))))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
