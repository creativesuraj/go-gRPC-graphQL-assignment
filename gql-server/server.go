package main

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aeon/gql-server/api"
	"github.com/aeon/gql-server/api/resolvers"
	"github.com/aeon/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	gqlEnv := utils.LoadEnvConfig("GQL")
	url := gqlEnv.GetString("SERVER_URL")
	srv := handler.NewDefaultServer(api.NewExecutableSchema(api.Config{Resolvers: &resolvers.Resolver{}}))

	http.Handle("/play", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	if url != "" {
		logrus.Infof("Listening to %s GraphQL server", url)
		logrus.Fatal(http.ListenAndServe(url, nil))
	} else {
		logrus.Fatal("Graphql server port is undefined")
	}

}
