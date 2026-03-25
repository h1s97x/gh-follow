package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/h1s97x/gh-follow/internal/cmd"
	"github.com/h1s97x/gh-follow/internal/flags"
)

func main() {
	app := &cli.App{
		Name:    "gh-follow",
		Usage:   "Manage your GitHub follow list from the terminal",
		Version: Version,
		Authors: []*cli.Author{
			{
				Name:  "h1s97x",
				Email: "h1s97x@users.noreply.github.com",
			},
		},
		Commands: []*cli.Command{
			{
				Name:     "add",
				Aliases:  []string{"follow"},
				Usage:    "Follow one or more GitHub users",
				Category: "Management",
				Flags:    flags.AddFlags(),
				Action:   cmd.Add,
			},
			{
				Name:     "remove",
				Aliases:  []string{"unfollow", "rm", "delete"},
				Usage:    "Unfollow one or more GitHub users",
				Category: "Management",
				Flags:    flags.RemoveFlags(),
				Action:   cmd.Remove,
			},
			{
				Name:     "list",
				Aliases:  []string{"ls"},
				Usage:    "List all users you follow",
				Category: "Display",
				Flags:    flags.ListFlags(),
				Action:   cmd.List,
			},
			{
				Name:     "sync",
				Usage:    "Sync follow list with GitHub or Gist",
				Category: "Sync",
				Flags:    flags.SyncFlags(),
				Action:   cmd.Sync,
			},
			{
				Name:     "export",
				Aliases:  []string{"exp"},
				Usage:    "Export follow list to a file",
				Category: "Data",
				Flags:    flags.ExportFlags(),
				Action:   cmd.Export,
			},
			{
				Name:     "import",
				Aliases:  []string{"imp"},
				Usage:    "Import follow list from a file",
				Category: "Data",
				Flags:    flags.ImportFlags(),
				Action:   cmd.Import,
			},
			{
				Name:     "stats",
				Usage:    "Show follow list statistics",
				Category: "Display",
				Flags:    flags.StatsFlags(),
				Action:   cmd.StatsCmd,
			},
			{
				Name:     "config",
				Usage:    "Manage configuration",
				Category: "Settings",
				Action:   cmd.ConfigCmd,
				Subcommands: []*cli.Command{
					{
						Name:   "show",
						Usage:  "Show current configuration",
						Action: cmd.ConfigShow,
					},
					{
						Name:   "get",
						Usage:  "Get a configuration value",
						Action: cmd.ConfigGet,
					},
					{
						Name:   "set",
						Usage:  "Set a configuration value",
						Action: cmd.ConfigSet,
					},
					{
						Name:   "reset",
						Usage:  "Reset configuration to defaults",
						Flags:  flags.ConfigResetFlags(),
						Action: cmd.ConfigReset,
					},
				},
			},
			{
				Name:     "gist",
				Usage:    "Manage Gist sync for follow list",
				Category: "Sync",
				Action:   cmd.Gist,
				Subcommands: []*cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new Gist for sync",
						Action: cmd.GistCreate,
					},
					{
						Name:   "status",
						Usage:  "Show Gist sync status",
						Action: cmd.GistStatus,
					},
					{
						Name:   "pull",
						Usage:  "Pull follow list from Gist",
						Flags:  flags.GistPullFlags(),
						Action: cmd.GistPull,
					},
					{
						Name:   "push",
						Usage:  "Push follow list to Gist",
						Flags:  flags.GistPushFlags(),
						Action: cmd.GistPush,
					},
				},
			},
			{
				Name:     "autosync",
				Usage:    "Manage auto-sync settings",
				Category: "Sync",
				Subcommands: []*cli.Command{
					{
						Name:   "status",
						Usage:  "Show auto-sync status",
						Action: cmd.AutoSyncStatus,
					},
					{
						Name:   "trigger",
						Usage:  "Manually trigger a sync",
						Flags:  flags.AutoSyncFlags(),
						Action: cmd.TriggerSync,
					},
				},
			},
			{
				Name:     "cache",
				Usage:    "Manage user information cache",
				Category: "Performance",
				Action:   cmd.CacheCmd,
				Subcommands: []*cli.Command{
					{
						Name:   "status",
						Usage:  "Show cache status",
						Action: cmd.CacheStatus,
					},
					{
						Name:   "list",
						Usage:  "List cached users",
						Flags:  flags.CacheListFlags(),
						Action: cmd.CacheList,
					},
					{
						Name:   "clear",
						Usage:  "Clear the cache",
						Flags:  flags.CacheClearFlags(),
						Action: cmd.CacheClear,
					},
					{
						Name:   "cleanup",
						Usage:  "Remove expired cache entries",
						Action: cmd.CacheCleanup,
					},
					{
						Name:   "refresh",
						Usage:  "Refresh cache for all followed users",
						Action: cmd.CacheRefresh,
					},
					{
						Name:   "show",
						Usage:  "Show cached user details",
						Action: cmd.CacheShow,
					},
				},
			},
			{
				Name:     "suggest",
				Aliases:  []string{"recommend"},
				Usage:    "Get follow suggestions",
				Category: "Discovery",
				Action:   cmd.Suggest,
				Flags:    flags.SuggestFlags(),
				Subcommands: []*cli.Command{
					{
						Name:   "trending",
						Usage:  "Show trending users",
						Flags:  flags.SuggestTrendingFlags(),
						Action: cmd.SuggestTrending,
					},
					{
						Name:   "mutual",
						Usage:  "Check mutual followers",
						Flags:  flags.SuggestMutualFlags(),
						Action: cmd.SuggestMutual,
					},
					{
						Name:   "inactive",
						Usage:  "Find inactive followed users",
						Flags:  flags.SuggestInactiveFlags(),
						Action: cmd.SuggestInactive,
					},
				},
			},
			{
				Name:     "batch",
				Usage:    "Perform batch operations",
				Category: "Management",
				Action:   cmd.Batch,
				Subcommands: []*cli.Command{
					{
						Name:   "follow",
						Usage:  "Follow multiple users at once",
						Flags:  flags.BatchFollowFlags(),
						Action: cmd.BatchFollow,
					},
					{
						Name:   "unfollow",
						Usage:  "Unfollow multiple users at once",
						Flags:  flags.BatchUnfollowFlags(),
						Action: cmd.BatchUnfollow,
					},
					{
						Name:   "check",
						Usage:  "Check if multiple users follow you",
						Flags:  flags.BatchCheckFlags(),
						Action: cmd.BatchCheck,
					},
					{
						Name:   "import",
						Usage:  "Import usernames from file",
						Flags:  flags.BatchImportFlags(),
						Action: cmd.BatchImport,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
