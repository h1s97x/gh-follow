package flags

import "github.com/urfave/cli/v2"

// ConfigFlags returns the flags for the config command
func ConfigFlags() []cli.Flag {
	return []cli.Flag{}
}

// ConfigSetFlags returns the flags for the config set subcommand
func ConfigSetFlags() []cli.Flag {
	return []cli.Flag{}
}

// ConfigGetFlags returns the flags for the config get subcommand
func ConfigGetFlags() []cli.Flag {
	return []cli.Flag{}
}

// ConfigResetFlags returns the flags for the config reset subcommand
func ConfigResetFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:  "force",
			Usage: "Skip confirmation",
			Value: false,
		},
	}
}
