package infra

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"bean/pkg/infra/gql"
	"bean/pkg/util"
)

func (this *Can) HttpRouter(router *mux.Router) *mux.Router {
	router.HandleFunc("/query", this.handleQueryRequest())

	if this.HttpServer.GraphQL.Playround.Enabled {
		hdl := playground.Handler(this.HttpServer.GraphQL.Playround.Title, "/query")
		router.Handle(this.HttpServer.GraphQL.Playround.Path, hdl)
	}

	return router
}

// Handle request to /query.
//  Verify JWT authorization if provided.
func (this *Can) handleQueryRequest() func(http.ResponseWriter, *http.Request) {
	cnf := gql.Config{Resolvers: this.graph}
	schema := gql.NewExecutableSchema(cnf)
	hdl := handler.NewDefaultServer(schema)

	return func(w http.ResponseWriter, r *http.Request) {
		err := this.beforeServeHTTP(r)
		if nil != err {
			w.WriteHeader(http.StatusForbidden)

			body := graphql.Response{
				Errors: gqlerror.List{
					{
						Message: err.Error(),
					},
				},
			}

			content, _ := json.Marshal(body)

			w.Write(content)
		} else {
			hdl.ServeHTTP(w, r)
		}
	}
}

func (this *Can) beforeServeHTTP(r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	if "" != authHeader {
		module, err := this.modules.Access()
		if nil != err {
			return errors.Wrap(err, util.ErrorCodeConfig.String())
		}

		claims, err := module.SessionResolver.JwtValidation(authHeader)
		if err != nil {
			return err
		} else if nil != claims {
			ctx := context.WithValue(r.Context(), "bean.claims", claims)
			r = r.WithContext(ctx)
		}
	}

	return nil
}
