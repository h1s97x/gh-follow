package flags

import "github.com/urfave/cli/v2"

// SyncFlags returns the flags for the sync command
func SyncFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "direction",
			Usage: "Sync direction: pull, push, both",
			Value: "both",
		},
		&cli.BoolFlag{
			Name:  "gist",
			Usage: "Use Gist as cloud storage",
			Value: false,
		},
		&cli.StringFlag{
			Name:  "gist-id",
			Usage: "Specify Gist ID",
		},
		&cli.StringFlag{
			Name:  "hostname",
			Usage: "GitHub hostname (for GitHub Enterprise)",
			Value: "github.com",
		},
		&cli.BoolFlag{
			Name:  "dry-run",
			Usage: "Simulate without making changes",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "force",
			Usage: "Skip confirmation",
			Value: false,
		},
	}
}
