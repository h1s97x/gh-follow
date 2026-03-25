package flags

import "github.com/urfave/cli/v2"

// AutoSyncFlags returns the flags for auto-sync commands
func AutoSyncFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:  "gist",
			Usage: "Use Gist for sync",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "silent",
			Usage: "Do not output any content",
			Value: false,
		},
	}
}
