package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/h1s97x/gh-follow/internal/models"
	"github.com/h1s97x/gh-follow/internal/storage"
)

// StatsCmd handles the follow stats command
func StatsCmd(c *cli.Context) error {
	// Get flags
	format := c.String("format")

	// Initialize storage
	st := storage.NewStorage(storage.DefaultStoragePath())

	// Load follow list
	list, err := st.Load()
	if err != nil {
		return fmt.Errorf("failed to load follow list: %w", err)
	}

	if len(list.Follows) == 0 {
		fmt.Println("Your follow list is empty")
		return nil
	}

	// Get stats
	stats := list.GetStats()

	// Output based on format
	switch format {
	case "json":
		return outputStatsJSON(stats)
	case "table":
		fallthrough
	default:
		return outputStatsTable(stats)
	}
}

// outputStatsTable outputs stats in table format
func outputStatsTable(stats *models.Stats) error {
	fmt.Println("\n📊 Follow List Statistics")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("Total follows: %d\n", stats.TotalFollows)
	fmt.Printf("Last updated: %s\n", stats.LastUpdated.Format(time.RFC3339))

	if stats.OldestFollow != nil {
		fmt.Printf("Oldest follow: %s (%s)\n",
			stats.OldestFollow.Username,
			stats.OldestFollow.FollowedAt.Format("2006-01-02"))
	}

	fmt.Println("\n📅 Recent Follows:")
	fmt.Println(strings.Repeat("-", 40))
	for _, f := range stats.RecentFollows {
		fmt.Printf("  %s - %s\n", f.Username, f.FollowedAt.Format("2006-01-02"))
	}

	if len(stats.Tags) > 0 {
		fmt.Println("\n🏷️ Tags Used:")
		fmt.Println(strings.Repeat("-", 40))
		fmt.Printf("  %s\n", strings.Join(stats.Tags, ", "))
	}

	return nil
}

// outputStatsJSON outputs stats in JSON format
func outputStatsJSON(stats *models.Stats) error {
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
