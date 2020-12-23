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
	"bean/components/util"
	"bean/components/util/connect"
	"bean/pkg/infra/gql"
)

type GraphqlHttpRouter struct {
	Container *Container
}

func (r *GraphqlHttpRouter) Handler(router *mux.Router) *mux.Router {
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

	schema := gql.NewExecutableSchema(cnf)
	srv := handler.New(schema)
	if r.Container.HttpServer.GraphQL.Transports.Post {
		srv.AddTransport(transport.POST{})
	}

	if r.Container.HttpServer.GraphQL.Transports.Websocket.KeepAlivePingInterval != 0 {
		srv.AddTransport(transport.Websocket{KeepAlivePingInterval: r.Container.HttpServer.GraphQL.Transports.Websocket.KeepAlivePingInterval})
	}

	if r.Container.HttpServer.GraphQL.Introspection {
		srv.Use(extension.Introspection{})
	}

	router.HandleFunc("/query", r.handleFunc(srv))

	if r.Container.HttpServer.GraphQL.Playround.Enabled {
		hdl := playground.Handler(r.Container.HttpServer.GraphQL.Playround.Title, "/query")
		router.Handle(r.Container.HttpServer.GraphQL.Playround.Path, hdl)
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
		con, err := r.Container.dbs.master()
		if nil != err {
			r.respondError(w, err, "failed to make DB connection", http.StatusInternalServerError)
		} else {
			context.WithValue(ctx, connect.DatabaseContextKey, con.WithContext(ctx))
		}

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
			ctx = context.WithValue(ctx, claim.ClaimsContextKey, claims)
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