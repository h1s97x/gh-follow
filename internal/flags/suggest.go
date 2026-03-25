package flags

import "github.com/urfave/cli/v2"

// SuggestFlags returns the flags for the suggest command
func SuggestFlags() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:  "limit",
			Usage: "Maximum number of suggestions to show",
			Value: 20,
		},
		&cli.StringFlag{
			Name:  "format",
			Usage: "Output format: table, json",
			Value: "table",
		},
		&cli.BoolFlag{
			Name:  "no-follow-of-follow",
			Usage: "Exclude suggestions from followed users",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "no-org-members",
			Usage: "Exclude organization members",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "no-star-contributors",
			Usage: "Exclude starred repo contributors",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "no-similar",
			Usage: "Exclude similar users",
			Value: false,
		},
	}
}

// SuggestTrendingFlags returns the flags for the trending command
func SuggestTrendingFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "language",
			Usage: "Filter by programming language",
			Value: "",
		},
		&cli.IntFlag{
			Name:  "limit",
			Usage: "Maximum number of users to show",
			Value: 10,
		},
	}
}

// SuggestMutualFlags returns the flags for the mutual command
func SuggestMutualFlags() []cli.Flag {
	return []cli.Flag{}
}

// SuggestInactiveFlags returns the flags for the inactive command
func SuggestInactiveFlags() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:  "days",
			Usage: "Days of inactivity threshold",
			Value: 365,
		},
	}
}
