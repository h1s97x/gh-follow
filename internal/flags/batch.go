package flags

import "github.com/urfave/cli/v2"

// BatchFlags returns the flags for the batch command
func BatchFlags() []cli.Flag {
	return []cli.Flag{}
}

// BatchFollowFlags returns the flags for batch follow command
func BatchFollowFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "file",
			Usage: "Read usernames from file (one per line)",
			Value: "",
		},
		&cli.IntFlag{
			Name:  "concurrency",
			Usage: "Number of concurrent operations",
			Value: 5,
		},
		&cli.IntFlag{
			Name:  "rate-limit",
			Usage: "Rate limit delay in milliseconds",
			Value: 100,
		},
		&cli.BoolFlag{
			Name:  "dry-run",
			Usage: "Preview without making changes",
			Value: false,
		},
	}
}

// BatchUnfollowFlags returns the flags for batch unfollow command
func BatchUnfollowFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "file",
			Usage: "Read usernames from file (one per line)",
			Value: "",
		},
		&cli.IntFlag{
			Name:  "concurrency",
			Usage: "Number of concurrent operations",
			Value: 5,
		},
		&cli.IntFlag{
			Name:  "rate-limit",
			Usage: "Rate limit delay in milliseconds",
			Value: 100,
		},
		&cli.BoolFlag{
			Name:  "dry-run",
			Usage: "Preview without making changes",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "force",
			Usage: "Skip confirmation",
			Value: false,
		},
	}
}

// BatchCheckFlags returns the flags for batch check command
func BatchCheckFlags() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:  "concurrency",
			Usage: "Number of concurrent operations",
			Value: 5,
		},
		&cli.IntFlag{
			Name:  "rate-limit",
			Usage: "Rate limit delay in milliseconds",
			Value: 100,
		},
	}
}

// BatchImportFlags returns the flags for batch import command
func BatchImportFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:  "follow",
			Usage: "Follow all imported users",
			Value: false,
		},
		&cli.IntFlag{
			Name:  "concurrency",
			Usage: "Number of concurrent operations",
			Value: 5,
		},
		&cli.IntFlag{
			Name:  "rate-limit",
			Usage: "Rate limit delay in milliseconds",
			Value: 100,
		},
		&cli.BoolFlag{
			Name:  "dry-run",
			Usage: "Preview without making changes",
			Value: false,
		},
	}
}
