package cmd

import (
	"github.com/urfave/cli/v2"

	"bean/pkg/infra"
)

func MigrationCommand(container *infra.Container) *cli.Command {
	return &cli.Command{
		Name: "migrate",
		Action: func(ctx *cli.Context) error {
			return container.Install(ctx.Context)
		},
	}
}
