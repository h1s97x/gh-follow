package flags

import "github.com/urfave/cli/v2"

// ListFlags returns the flags for the list command
func ListFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "format",
			Usage: "Output format: table, json, csv, simple",
			Value: "table",
		},
		&cli.StringFlag{
			Name:  "sort",
			Usage: "Sort field: name, date, notes",
			Value: "date",
		},
		&cli.StringFlag{
			Name:  "order",
			Usage: "Sort order: asc, desc",
			Value: "desc",
		},
		&cli.IntFlag{
			Name:  "limit",
			Usage: "Limit output count",
			Value: 0,
		},
		&cli.StringFlag{
			Name:  "filter",
			Usage: "Filter usernames by pattern",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "tag",
			Usage: "Filter by tag",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "date-from",
			Usage: "Filter follows from date (YYYY-MM-DD)",
			Value: "",
		},
		&cli.StringFlag{
			Name:  "date-to",
			Usage: "Filter follows to date (YYYY-MM-DD)",
			Value: "",
		},
	}
}
