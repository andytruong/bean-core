package cmd

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"bean/pkg/infra"
)

func HttpServerCommand(container *infra.Container) *cli.Command {
	return &cli.Command{
		Name: "http-server",
		Action: func(ctx *cli.Context) error {
			router := mux.NewRouter()
			r := infra.GraphqlHttpRouter{Container: container}

			server := http.Server{
				Addr:              container.HttpServer.Address,
				Handler:           r.Handler(router),
				ReadTimeout:       container.HttpServer.ReadTimeout,
				ReadHeaderTimeout: container.HttpServer.ReadTimeout,
				WriteTimeout:      container.HttpServer.WriteTimeout,
				IdleTimeout:       container.HttpServer.ReadTimeout,
			}

			container.Logger().Info("🚀 HTTP server is running", zap.String("address", container.HttpServer.Address))

			return server.ListenAndServe()
		},
	}
}
