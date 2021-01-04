package main

import (
	"os"

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
			cmdMigrate(container),
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
			server := c.HttpServer()

			c.Logger().Info("🚀 HTTP server is running", zap.String("address", c.Config.HttpServer.Address))

			return server.ListenAndServe()
		},
	}
}

func cmdMigrate(container *infra.Container) *cli.Command {
	return &cli.Command{
		Name: "migrate",
		Action: func(ctx *cli.Context) error {
			db, err := container.DBs.PrepareStmt(false).Master()
			if nil != err {
				return err
			}

			bundles := container.BundleList()

			return connect.Migrate(ctx.Context, bundles.Get(), db)
		},
	}
}
