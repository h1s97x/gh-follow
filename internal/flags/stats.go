package flags

import "github.com/urfave/cli/v2"

// StatsFlags returns the flags for the stats command
func StatsFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "format",
			Usage: "Output format: table, json",
			Value: "table",
		},
	}
}
