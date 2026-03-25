package flags

import "github.com/urfave/cli/v2"

// CacheFlags returns the flags for the cache command
func CacheFlags() []cli.Flag {
	return []cli.Flag{}
}

// CacheListFlags returns the flags for cache list command
func CacheListFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "format",
			Usage: "Output format: table, json",
			Value: "table",
		},
	}
}

// CacheClearFlags returns the flags for cache clear command
func CacheClearFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:  "force",
			Usage: "Skip confirmation",
			Value: false,
		},
	}
}

// CacheRefreshFlags returns the flags for cache refresh command
func CacheRefreshFlags() []cli.Flag {
	return []cli.Flag{}
}

// CacheShowFlags returns the flags for cache show command
func CacheShowFlags() []cli.Flag {
	return []cli.Flag{}
}
