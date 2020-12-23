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

func (this *Container) HttpRouter(router *mux.Router) *mux.Router {
	cnf := gql.Config{
		Resolvers: &Resolver{container: this},
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
	if this.HttpServer.GraphQL.Transports.Post {
		srv.AddTransport(transport.POST{})
	}
	
	if this.HttpServer.GraphQL.Transports.Websocket.KeepAlivePingInterval != 0 {
		srv.AddTransport(transport.Websocket{KeepAlivePingInterval: this.HttpServer.GraphQL.Transports.Websocket.KeepAlivePingInterval})
	}
	
	if this.HttpServer.GraphQL.Introspection {
		srv.Use(extension.Introspection{})
	}
	
	router.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		//  Verify JWT authorization if provided.
		if err := this.beforeServeHTTP(r); nil != err {
			this.respond403(w, err)
		} else {
			srv.ServeHTTP(w, r)
		}
	})
	
	if this.HttpServer.GraphQL.Playround.Enabled {
		hdl := playground.Handler(this.HttpServer.GraphQL.Playround.Title, "/query")
		router.Handle(this.HttpServer.GraphQL.Playround.Path, hdl)
	}
	
	return router
}

func (this *Container) beforeServeHTTP(r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	if "" != authHeader {
		bundle, err := this.bundles.Access()
		if nil != err {
			return errors.Wrap(err, util.ErrorCodeConfig.String())
		}
		
		claims, err := bundle.JwtService.Validate(authHeader)
		
		if err != nil {
			return err
		} else if nil != claims {
			ctx := context.WithValue(r.Context(), claim.ContextKey, claims)
			r = r.WithContext(ctx)
		}
	}
	
	return nil
}

func (this *Container) respond403(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusForbidden)
	errList := gqlerror.List{{Message: err.Error()}}
	body := graphql.Response{Errors: errList}
	content, _ := json.Marshal(body)
	
	_, err = w.Write(content)
	if nil != err {
		this.logger.Error("failed responding", zap.String("message", err.Error()))
	}
}
