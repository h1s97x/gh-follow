package flags

import "github.com/urfave/cli/v2"

// AddFlags returns the flags for the add command
func AddFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:  "sync",
			Usage: "Sync follow to GitHub (default: true)",
			Value: true,
		},
		&cli.BoolFlag{
			Name:  "gist",
			Usage: "Sync to Gist after adding",
			Value: false,
		},
		&cli.StringFlag{
			Name:  "gist-id",
			Usage: "Specify Gist ID for sync",
		},
		&cli.StringFlag{
			Name:  "notes",
			Usage: "Add notes for the user(s)",
			Value: "",
		},
		&cli.StringSliceFlag{
			Name:  "tags",
			Usage: "Add tags for the user(s)",
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
