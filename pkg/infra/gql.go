package infra

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"

	"bean/pkg/infra/gql"
)

func NewQueryRoute(container *Container) (http.Handler, error) {
	resolver, err := gql.NewResolver(container)
	if nil != err {
		return nil, err
	}

	config := &gql.Config{Resolvers: resolver}
	schema := gql.NewExecutableSchema(*config)
	server := handler.NewDefaultServer(schema)

	return &QueryHandler{
		config:   config,
		resolver: resolver,
		server:   server,
	}, nil
}

type QueryHandler struct {
	config   *gql.Config
	resolver gql.ResolverRoot
	server   *handler.Server
}

func (this QueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: convert session token -> session, or just reject the request.

	this.server.ServeHTTP(w, r)
}
