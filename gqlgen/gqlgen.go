package gqlgen

import (
	"github.com/varunturlapati/vtgqlgen/dataloaders"
	"net/http"

	"github.com/99designs/gqlgen/handler"

	ds "github.com/varunturlapati/vtgqlgen/datasource"
)

// NewHandler returns a new graphql endpoint handler.
func NewHandler(repo ds.Repository, dl dataloaders.Retriever) http.Handler {

	return handler.GraphQL(NewExecutableSchema(Config{
		Resolvers: &Resolver{
			Repository: repo,
			DataLoaders: dl,
		},
	}))
}

// NewPlaygroundHandler returns a new GraphQL Playground handler.
func NewPlaygroundHandler(endpoint string) http.Handler {
	return handler.Playground("GraphQL Playground", endpoint)
}
