package flags

import "github.com/urfave/cli/v2"

// RemoveFlags returns the flags for the remove command
func RemoveFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:  "sync",
			Usage: "Sync unfollow to GitHub (default: true)",
			Value: true,
		},
		&cli.BoolFlag{
			Name:  "gist",
			Usage: "Sync to Gist after removing",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "force",
			Usage: "Skip confirmation",
			Value: false,
		},
		&cli.StringFlag{
			Name:  "hostname",
			Usage: "GitHub hostname (for GitHub Enterprise)",
			Value: "github.com",
		},
		&cli.BoolFlag{
			Name:  "silent",
			Usage: "Do not output any content",
			Value: false,
		},
	}
}
