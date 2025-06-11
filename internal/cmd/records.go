package cmd

import "github.com/urfave/cli/v3"

var records = &cli.Command{
	Name: "records",
	Commands: []*cli.Command{
		getRecords,
	},
}
