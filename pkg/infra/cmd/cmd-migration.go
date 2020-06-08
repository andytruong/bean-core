package cmd

import (
	"github.com/urfave/cli/v2"

	"bean/pkg/infra"
)

func MigrationCommand(can *infra.Can) *cli.Command {
	return &cli.Command{
		Name: "migrate",
		Action: func(ctx *cli.Context) error {
			return can.Migrate(ctx.Context)
		},
	}
}
