package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"bean/components/connect"
	"bean/components/util"
	"bean/pkg/infra"
)

func main() {
	path := os.Getenv("CONFIG")
	if path == "" {
		err := errors.Wrap(util.ErrorConfig, "missing env CONFIG")
		panic(err)
	}

	container, err := infra.NewContainer(path)
	if nil != err {
		panic("failed creating container: " + err.Error())
	}

	app := cli.App{
		Name: "bean",
		Commands: []*cli.Command{
			cmdHttpServer(container),
			{
				Name: "migrate",
				Action: func(ctx *cli.Context) error {
					db, err := container.DBs.Master()
					if nil != err {
						return err
					}

					bundles := container.BundleList()

					return connect.Migrate(ctx.Context, bundles.Get(), db)
				},
			},
		},
	}

	if err := app.Run(os.Args); nil != err {
		panic(err)
	}
}

func cmdHttpServer(c *infra.Container) *cli.Command {
	return &cli.Command{
		Name: "http-server",
		Action: func(ctx *cli.Context) error {
			router := mux.NewRouter()
			r := infra.GraphqlHttpRouter{Container: c}

			server := http.Server{
				Addr:              c.Config.HttpServer.Address,
				Handler:           r.Handler(router),
				ReadTimeout:       c.Config.HttpServer.ReadTimeout,
				ReadHeaderTimeout: c.Config.HttpServer.ReadTimeout,
				WriteTimeout:      c.Config.HttpServer.WriteTimeout,
				IdleTimeout:       c.Config.HttpServer.ReadTimeout,
			}

			c.Logger().Info("🚀 HTTP server is running", zap.String("address", c.Config.HttpServer.Address))

			return server.ListenAndServe()
		},
	}
}
