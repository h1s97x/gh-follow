package storage

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/h1s97x/gh-follow/internal/github"
	"github.com/h1s97x/gh-follow/internal/models"
)

// Import handles the follow import command
func Import(c *cli.Context) error {
	// Get flags
	input := c.String("input")
	merge := c.Bool("merge")
	sync := c.Bool("sync")
	dryRun := c.Bool("dry-run")

	if input == "" {
		return fmt.Errorf("please provide an input file with --input")
	}

	// Initialize storage
	storage := NewStorage(DefaultStoragePath())

	// Load input file to preview
	list, err := storage.Load()
	if err != nil {
		list = models.NewFollowList()
	}

	// Count before
	countBefore := len(list.Follows)

	if dryRun {
		fmt.Printf("Dry run: would import from %s\n", input)
		fmt.Printf("Current list has %d users\n", countBefore)
		return nil
	}

	// Import
	if err := storage.Import(input, merge); err != nil {
		return fmt.Errorf("failed to import follow list: %w", err)
	}

	// Load after import
	newList, err := storage.Load()
	if err != nil {
		return fmt.Errorf("failed to load new list: %w", err)
	}

	countAfter := len(newList.Follows)
	added := countAfter - countBefore

	fmt.Printf("Imported follow list from %s\n", input)
	fmt.Printf("Added %d new users, total: %d\n", added, countAfter)

	// Sync to GitHub if requested
	if sync {
		token, err := github.GetTokenFromGH()
		if err != nil {
			fmt.Printf("Warning: Could not get GitHub token: %v\n", err)
			return nil
		}

		gc := github.NewGitHubClient(token, "github.com")
		ctx := context.Background()

		// Follow all users in the list on GitHub
		for _, f := range newList.Follows {
			if err := gc.Follow(ctx, f.Username); err != nil {
				fmt.Printf("Warning: Failed to follow %s on GitHub: %v\n", f.Username, err)
			} else {
				fmt.Printf("Followed %s on GitHub\n", f.Username)
			}
		}
	}

	return nil
}
