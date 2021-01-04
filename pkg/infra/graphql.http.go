package infra

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.uber.org/zap"

	"bean/components/claim"
	"bean/components/connect"
	"bean/components/util"
	"bean/pkg/infra/gql"
)

type graphqlHttpHandler struct {
	container *Container
}

func (h *graphqlHttpHandler) config() gql.Config {
	cnf := gql.Config{
		Resolvers: &Resolver{container: h.container},
		Directives: gql.DirectiveRoot{
			Constraint: func(ctx context.Context, obj interface{}, next graphql.Resolver, maxLength *int, minLength *int) (
				res interface{}, err error,
			) {
				// TODO: implement me

				return next(ctx)
			},
			RequireAuth: func(ctx context.Context, obj interface{}, next graphql.Resolver) (
				res interface{}, err error,
			) {
				// TODO: implement me

				return next(ctx)
			},
		},
	}

	return cnf
}

func (h *graphqlHttpHandler) Get(router *mux.Router) *mux.Router {
	router.HandleFunc("/query", h.callback())

	if h.container.Config.HttpServer.GraphQL.Playround.Enabled {
		router.Handle(
			h.container.Config.HttpServer.GraphQL.Playround.Path,
			playground.Handler(
				h.container.Config.HttpServer.GraphQL.Playround.Title,
				"/query",
			),
		)
	}

	return router
}

func (h *graphqlHttpHandler) defaultCallback() func(http.ResponseWriter, *http.Request) {
	config := h.config()
	schema := gql.NewExecutableSchema(config)
	server := handler.New(schema)

	if h.container.Config.HttpServer.GraphQL.Transports.Post {
		server.AddTransport(transport.POST{})
	}

	if h.container.Config.HttpServer.GraphQL.Transports.Websocket.KeepAlivePingInterval != 0 {
		trans := transport.Websocket{KeepAlivePingInterval: h.container.Config.HttpServer.GraphQL.Transports.Websocket.KeepAlivePingInterval}
		server.AddTransport(trans)
	}

	if h.container.Config.HttpServer.GraphQL.Introspection {
		server.Use(extension.Introspection{})
	}

	return server.ServeHTTP
}

func (h *graphqlHttpHandler) callback() func(http.ResponseWriter, *http.Request) {
	defaultCallback := h.defaultCallback()

	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		//  Verify JWT authorization if provided.
		claimContext, err := h.beforeServe(req)
		if nil != err {
			h.respondError(w, err, "failed responding", http.StatusForbidden)
		} else if nil != claimContext {
			ctx = claimContext
		}

		// Inject DB connection to context
		ctx = connect.WithContextValue(ctx, h.container.DBs)

		defaultCallback(w, req.WithContext(ctx))
	}
}

func (h *graphqlHttpHandler) beforeServe(req *http.Request) (context.Context, error) {
	authHeader := req.Header.Get("Authorization")

	// no authentication header, no need to validate
	if authHeader == "" {
		return nil, nil
	}

	// invoke access bundle to validate the auth header
	ctx := req.Context()
	accessBundle, err := h.container.bundles.Access()
	if nil != err {
		return nil, errors.Wrap(err, util.ErrorCodeConfig.String())
	}

	claims, err := accessBundle.JwtService.Validate(authHeader)
	if err != nil {
		return nil, err
	}

	if nil != claims {
		ctx = claim.PayloadToContext(ctx, claims)
	}

	return ctx, nil
}

func (h *graphqlHttpHandler) respondError(w http.ResponseWriter, err error, msg string, status int) {
	w.WriteHeader(status)
	errList := gqlerror.List{{Message: err.Error()}}
	body := graphql.Response{Errors: errList}
	content, _ := json.Marshal(body)

	_, err = w.Write(content)
	if nil != err {
		h.container.logger.Error(msg, zap.String("message", err.Error()))
	}
}
