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

type GraphqlHttpRouter struct {
	Container *Container
}

func (r *GraphqlHttpRouter) config() gql.Config {
	cnf := gql.Config{
		Resolvers: &Resolver{container: r.Container},
		Directives: gql.DirectiveRoot{
			Constraint: func(ctx context.Context, obj interface{}, next graphql.Resolver, maxLength *int, minLength *int) (res interface{}, err error) {
				// TODO: implement me

				return next(ctx)
			},
			RequireAuth: func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
				// TODO: implement me

				return next(ctx)
			},
		},
	}

	return cnf
}

func (r *GraphqlHttpRouter) Handler(router *mux.Router) *mux.Router {
	config := r.config()
	schema := gql.NewExecutableSchema(config)
	server := handler.New(schema)

	if r.Container.Config.HttpServer.GraphQL.Transports.Post {
		server.AddTransport(transport.POST{})
	}

	if r.Container.Config.HttpServer.GraphQL.Transports.Websocket.KeepAlivePingInterval != 0 {
		server.AddTransport(transport.Websocket{KeepAlivePingInterval: r.Container.Config.HttpServer.GraphQL.Transports.Websocket.KeepAlivePingInterval})
	}

	if r.Container.Config.HttpServer.GraphQL.Introspection {
		server.Use(extension.Introspection{})
	}

	router.HandleFunc("/query", r.handleFunc(server))

	if r.Container.Config.HttpServer.GraphQL.Playround.Enabled {
		hdl := playground.Handler(r.Container.Config.HttpServer.GraphQL.Playround.Title, "/query")
		router.Handle(r.Container.Config.HttpServer.GraphQL.Playround.Path, hdl)
	}

	return router
}

func (r *GraphqlHttpRouter) handleFunc(srv *handler.Server) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		//  Verify JWT authorization if provided.
		claimContext, err := r.beforeServe(req)
		if nil != err {
			r.respondError(w, err, "failed responding", http.StatusForbidden)
		} else if nil != claimContext {
			ctx = claimContext
		}

		// Inject DB connection to context
		ctx = connect.WithContextValue(ctx, r.Container.dbs)

		srv.ServeHTTP(w, req.WithContext(ctx))
	}
}

func (r *GraphqlHttpRouter) beforeServe(req *http.Request) (context.Context, error) {
	ctx := req.Context()
	authHeader := req.Header.Get("Authorization")
	if authHeader != "" {
		bundle, err := r.Container.bundles.Access()
		if nil != err {
			return nil, errors.Wrap(err, util.ErrorCodeConfig.String())
		}

		claims, err := bundle.JwtService.Validate(authHeader)
		if err != nil {
			return nil, err
		}

		if nil != claims {
			ctx = claim.PayloadToContext(ctx, claims)
		}
	}

	return ctx, nil
}

func (r *GraphqlHttpRouter) respondError(w http.ResponseWriter, err error, msg string, status int) {
	w.WriteHeader(status)
	errList := gqlerror.List{{Message: err.Error()}}
	body := graphql.Response{Errors: errList}
	content, _ := json.Marshal(body)

	_, err = w.Write(content)
	if nil != err {
		r.Container.logger.Error(msg, zap.String("message", err.Error()))
	}
}
