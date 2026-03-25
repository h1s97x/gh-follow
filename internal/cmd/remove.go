package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	gh_client "github.com/h1s97x/gh-follow/internal/github"
	"github.com/h1s97x/gh-follow/internal/storage"
)

// Remove handles the follow remove command
func Remove(c *cli.Context) error {
	// Get usernames from arguments
	args := c.Args().Slice()
	if len(args) == 0 {
		return fmt.Errorf("please provide at least one username to unfollow")
	}

	// Get flags
	doSync := c.Bool("sync")
	force := c.Bool("force")
	silent := c.Bool("silent")

	// Initialize storage
	st := storage.NewStorage(storage.DefaultStoragePath())

	// Check if storage exists
	if !st.Exists() {
		return fmt.Errorf("no follow list found, please add some users first")
	}

	// Initialize GitHub client if sync is enabled
	var gc *gh_client.GitHubClient
	if doSync {
		token, err := gh_client.GetTokenFromGH()
		if err != nil {
			if !silent {
				fmt.Printf("Warning: Could not get GitHub token: %v\n", err)
				fmt.Println("Proceeding with local storage only...")
			}
			doSync = false
		} else {
			gc = gh_client.NewGitHubClient(token, c.String("hostname"))
		}
	}

	// Confirm if not forced
	if !force && !silent {
		fmt.Printf("Are you sure you want to unfollow %d user(s)? [y/N]: ", len(args))
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("Cancelled")
			return nil
		}
	}

	var succeeded []string
	var failed []string

	for _, username := range args {
		// Check if in local list
		list, err := st.Load()
		if err != nil {
			return fmt.Errorf("failed to load local list: %w", err)
		}

		if !list.Contains(username) {
			if !silent {
				fmt.Printf("User %s is not in your follow list\n", username)
			}
			continue
		}

		// Unfollow on GitHub if sync is enabled
		if doSync {
			ctx := context.Background()
			if err := gc.Unfollow(ctx, username); err != nil {
				if !silent {
					fmt.Printf("Failed to unfollow %s on GitHub: %v\n", username, err)
				}
				failed = append(failed, username)
				continue
			}
		}

		// Remove from local storage
		if err := st.Remove(username); err != nil {
			if !silent {
				fmt.Printf("Failed to remove %s from local list: %v\n", username, err)
			}
			failed = append(failed, username)
			continue
		}

		succeeded = append(succeeded, username)
		if !silent {
			action := "Removed"
			if doSync {
				action = "Unfollowed"
			}
			fmt.Printf("%s %s from your follow list\n", action, username)
		}
	}

	// Summary
	if !silent {
		fmt.Printf("\nSummary: %d succeeded, %d failed\n", len(succeeded), len(failed))
		if len(failed) > 0 {
			fmt.Printf("Failed: %v\n", failed)
		}
	}

	return nil
}
