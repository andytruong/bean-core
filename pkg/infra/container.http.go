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
	"bean/pkg/infra/gql"
)

func (container *Container) HttpRouter(router *mux.Router) *mux.Router {
	cnf := gql.Config{
		Resolvers: &Resolver{container: container},
		Directives: gql.DirectiveRoot{
			// Comment: nil - just for comment
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
	if container.HttpServer.GraphQL.Transports.Post {
		srv.AddTransport(transport.POST{})
	}

	if container.HttpServer.GraphQL.Transports.Websocket.KeepAlivePingInterval != 0 {
		srv.AddTransport(transport.Websocket{KeepAlivePingInterval: container.HttpServer.GraphQL.Transports.Websocket.KeepAlivePingInterval})
	}

	if container.HttpServer.GraphQL.Introspection {
		srv.Use(extension.Introspection{})
	}

	router.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		//  Verify JWT authorization if provided.
		if ctx, err := container.beforeServeHTTP(r); nil != err {
			container.respond403(w, err)
		} else {
			if nil != ctx {
				srv.ServeHTTP(w, r.WithContext(ctx))
			} else {
				srv.ServeHTTP(w, r)
			}
		}
	})

	if container.HttpServer.GraphQL.Playround.Enabled {
		hdl := playground.Handler(container.HttpServer.GraphQL.Playround.Title, "/query")
		router.Handle(container.HttpServer.GraphQL.Playround.Path, hdl)
	}

	return router
}

func (container *Container) beforeServeHTTP(r *http.Request) (context.Context, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		bundle, err := container.bundles.Access()
		if nil != err {
			return nil, errors.Wrap(err, util.ErrorCodeConfig.String())
		}

		claims, err := bundle.JwtService.Validate(authHeader)

		if err != nil {
			return nil, err
		} else if nil != claims {
			ctx := context.WithValue(r.Context(), claim.ContextKey, claims)
			return ctx, nil
		}
	}

	return nil, nil
}

func (container *Container) respond403(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusForbidden)
	errList := gqlerror.List{{Message: err.Error()}}
	body := graphql.Response{Errors: errList}
	content, _ := json.Marshal(body)

	_, err = w.Write(content)
	if nil != err {
		container.logger.Error("failed responding", zap.String("message", err.Error()))
	}
}
