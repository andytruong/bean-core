package cmd

import (
	"github.com/urfave/cli/v2"

	"bean/pkg/infra"
)

func HttpServerCommand(container *infra.Container) *cli.Command {
	return &cli.Command{
		Name: "http-server",
		Action: func(ctx *cli.Context) error {
			return container.ListenAndServe()
		},
	}
}
