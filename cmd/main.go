package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"bean/components/util"
	"bean/pkg/infra"
	"bean/pkg/infra/cmd"
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
			cmd.HttpServerCommand(container),
			cmd.MigrationCommand(container),
			cmd.KeyGenCommand(container),
		},
	}

	if err := app.Run(os.Args); nil != err {
		panic(err)
	}
}
