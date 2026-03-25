package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	gh_client "github.com/h1s97x/gh-follow/internal/github"
	"github.com/h1s97x/gh-follow/internal/storage"
)

// Add handles the follow add command
func Add(c *cli.Context) error {
	// Get usernames from arguments
	args := c.Args().Slice()
	if len(args) == 0 {
		return fmt.Errorf("please provide at least one username to follow")
	}

	// Get flags
	doSync := c.Bool("sync")
	notes := c.String("notes")
	tags := c.StringSlice("tags")
	silent := c.Bool("silent")

	// Initialize storage
	st := storage.NewStorage(storage.DefaultStoragePath())

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

	var succeeded []string
	var failed []string

	for _, username := range args {
		// Check if already in local list
		list, err := st.Load()
		if err != nil {
			return fmt.Errorf("failed to load local list: %w", err)
		}

		if list.Contains(username) {
			if !silent {
				fmt.Printf("User %s is already in your follow list\n", username)
			}
			continue
		}

		// Follow on GitHub if sync is enabled
		if doSync {
			ctx := context.Background()
			if err := gc.Follow(ctx, username); err != nil {
				if !silent {
					fmt.Printf("Failed to follow %s on GitHub: %v\n", username, err)
				}
				failed = append(failed, username)
				continue
			}
		}

		// Add to local storage
		if err := st.Add(username, notes, tags); err != nil {
			if !silent {
				fmt.Printf("Failed to add %s to local list: %v\n", username, err)
			}
			failed = append(failed, username)
			continue
		}

		succeeded = append(succeeded, username)
		if !silent {
			action := "Added"
			if doSync {
				action = "Followed"
			}
			fmt.Printf("%s %s to your follow list\n", action, username)
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
