package flags

import "github.com/urfave/cli/v2"

// GistFlags returns the flags for the gist command
func GistFlags() []cli.Flag {
	return []cli.Flag{}
}

// GistPullFlags returns the flags for gist pull command
func GistPullFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:  "force",
			Usage: "Force overwrite local changes",
			Value: false,
		},
	}
}

// GistPushFlags returns the flags for gist push command
func GistPushFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:  "force",
			Usage: "Force overwrite remote Gist",
			Value: false,
		},
	}
}
