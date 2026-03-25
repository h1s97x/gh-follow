package flags

import "github.com/urfave/cli/v2"

// ImportFlags returns the flags for the import command
func ImportFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "input",
			Usage: "Input file path (required)",
		},
		&cli.BoolFlag{
			Name:  "merge",
			Usage: "Merge with existing list",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "sync",
			Usage: "Sync to GitHub after import",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "dry-run",
			Usage: "Simulate without making changes",
			Value: false,
		},
	}
}
