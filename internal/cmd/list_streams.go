package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

func init() {
}

var listStreams = &cli.Command{
	Name: "list",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	},
}
