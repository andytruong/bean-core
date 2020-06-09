package cmd

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"bean/pkg/infra"
)

func HttpServerCommand(can *infra.Can) *cli.Command {
	router := mux.NewRouter()

	return &cli.Command{
		Name: "http-server",
		Action: func(ctx *cli.Context) error {
			server := http.Server{
				Addr:              can.HttpServer.Address,
				Handler:           can.HttpRouter(router),
				ReadTimeout:       can.HttpServer.ReadTimeout,
				ReadHeaderTimeout: can.HttpServer.ReadTimeout,
				WriteTimeout:      can.HttpServer.WriteTimeout,
				IdleTimeout:       can.HttpServer.ReadTimeout,
			}

			can.Logger().Info("🚀 HTTP server is running", zap.String("address", can.HttpServer.Address))

			return server.ListenAndServe()
		},
	}
}
