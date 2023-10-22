package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aeon/gql-server/api"
	"github.com/aeon/gql-server/api/resolvers"
)

const defaultPort = "5001"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(api.NewExecutableSchema(api.Config{Resolvers: &resolvers.Resolver{}}))

	http.Handle("/play", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("Listening to http://localhost:%s GraphQL server", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
