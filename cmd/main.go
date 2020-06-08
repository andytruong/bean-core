package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"bean/pkg/infra"
	"bean/pkg/infra/cmd"
	"bean/pkg/util"
)

func main() {
	path := os.Getenv("CONFIG")
	if "" == path {
		err := errors.Wrap(util.ErrorConfig, "missing env CONFIG")
		panic(err)
	}

	can, err := infra.NewCan(path)
	if nil != err {
		panic("failed creating can: " + err.Error())
	}

	app := cli.App{
		Name: "bean",
		Commands: []*cli.Command{
			cmd.HttpServerCommand(can),
			cmd.MigrationCommand(can),
			cmd.KeyGenCommand(can),
		},
	}

	if err := app.Run(os.Args); nil != err {
		panic(err)
	}
}
