package infra

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (this *Container) ListenAndServe() error {
	handler, err := NewQueryRoute(this)
	if nil != err {
		return err
	}

	router := mux.NewRouter()
	router.Handle("/query", handler)
	if this.HttpServer.GraphQL.Playround.Enabled {
		handler := playground.Handler(this.HttpServer.GraphQL.Playround.Title, "/query")
		router.Handle(this.HttpServer.GraphQL.Playround.Path, handler)
	}

	server := http.Server{
		Addr:              this.HttpServer.Address,
		Handler:           router,
		ReadTimeout:       this.HttpServer.ReadTimeout,
		ReadHeaderTimeout: this.HttpServer.ReadTimeout,
		WriteTimeout:      this.HttpServer.WriteTimeout,
		IdleTimeout:       this.HttpServer.ReadTimeout,
	}

	this.Logger.Info("🚀 HTTP server is running", zap.String("address", this.HttpServer.Address))

	return server.ListenAndServe()
}
