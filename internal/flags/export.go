package flags

import "github.com/urfave/cli/v2"

// ExportFlags returns the flags for the export command
func ExportFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "output",
			Usage: "Output file path",
			Value: "follows.json",
		},
		&cli.StringFlag{
			Name:  "format",
			Usage: "Output format: json, csv",
			Value: "json",
		},
	}
}
