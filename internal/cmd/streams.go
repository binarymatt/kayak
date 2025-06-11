package cmd

import "github.com/urfave/cli/v3"

var streams = &cli.Command{
	Name: "streams",
	Commands: []*cli.Command{
		streamStatsCommand,
		listStreams,
	},
}
