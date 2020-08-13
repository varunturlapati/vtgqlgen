package gqlgen

import (
	"net/http"

	"github.com/99designs/gqlgen/handler"

	ds "github.com/varunturlapati/vtgqlgen/datasource"
)

// NewHandler returns a new graphql endpoint handler.
func NewHandler(repo ds.Repository) http.Handler {

	return handler.GraphQL(NewExecutableSchema(Config{
		Resolvers: &Resolver{
			Repository: repo,
		},
	}))
}

// NewPlaygroundHandler returns a new GraphQL Playground handler.
func NewPlaygroundHandler(endpoint string) http.Handler {
	return handler.Playground("GraphQL Playground", endpoint)
}
