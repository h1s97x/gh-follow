package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	gh_client "github.com/h1s97x/gh-follow/internal/github"
	"github.com/h1s97x/gh-follow/internal/storage"
)

// Batch handles the batch command
func Batch(c *cli.Context) error {
	return fmt.Errorf("please specify a subcommand: follow, unfollow, or check")
}

// BatchFollow handles batch follow command
func BatchFollow(c *cli.Context) error {
	if c.Args().Len() == 0 {
		return fmt.Errorf("please provide usernames or use --file flag")
	}

	token, err := gh_client.GetTokenFromGH()
	if err != nil {
		return fmt.Errorf("failed to get GitHub token: %w", err)
	}

	gc := gh_client.NewGitHubClient(token, "github.com")
	st := storage.NewStorage(storage.DefaultStoragePath())

	// Get usernames
	usernames := c.Args().Slice()
	if file := c.String("file"); file != "" {
		usersFromFile, err := readUsersFromFile(file)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		usernames = append(usernames, usersFromFile...)
	}

	if len(usernames) == 0 {
		return fmt.Errorf("no usernames provided")
	}

	// Remove duplicates
	usernames = uniqueStrings(usernames)

	// Check dry run
	dryRun := c.Bool("dry-run")
	if dryRun {
		fmt.Printf("\n⚠️  DRY RUN - No changes will be made\n")
	}

	fmt.Printf("\n👥 Batch following %d users...\n", len(usernames))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	opts := &gh_client.BatchOptions{
		Concurrency: c.Int("concurrency"),
		RateLimit:   time.Duration(c.Int("rate-limit")) * time.Millisecond,
		DryRun:      dryRun,
		Progress:    true,
	}

	bp := gh_client.NewBatchProcessor(gc, opts)
	results := bp.BatchFollow(context.Background(), usernames, opts)

	// Display results
	gh_client.DisplayBatchResults(results, "Follow")

	// Save to storage
	if !dryRun {
		list, err := st.Load()
		if err != nil {
			return fmt.Errorf("failed to load follow list: %w", err)
		}
		for _, r := range results {
			if r.Success {
				list.Add(r.Username, "", nil)
			}
		}
		if err := st.Save(list); err != nil {
			return fmt.Errorf("failed to save follow list: %w", err)
		}
	}

	return nil
}

// BatchUnfollow handles batch unfollow command
func BatchUnfollow(c *cli.Context) error {
	if c.Args().Len() == 0 {
		return fmt.Errorf("please provide usernames or use --file flag")
	}

	token, err := gh_client.GetTokenFromGH()
	if err != nil {
		return fmt.Errorf("failed to get GitHub token: %w", err)
	}

	gc := gh_client.NewGitHubClient(token, "github.com")
	st := storage.NewStorage(storage.DefaultStoragePath())

	// Get usernames
	usernames := c.Args().Slice()
	if file := c.String("file"); file != "" {
		usersFromFile, err := readUsersFromFile(file)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		usernames = append(usernames, usersFromFile...)
	}

	if len(usernames) == 0 {
		return fmt.Errorf("no usernames provided")
	}

	// Remove duplicates
	usernames = uniqueStrings(usernames)

	// Check dry run
	dryRun := c.Bool("dry-run")
	if dryRun {
		fmt.Printf("\n⚠️  DRY RUN - No changes will be made\n")
	}

	// Confirm
	if !c.Bool("force") && !dryRun {
		fmt.Printf("Unfollow %d users? [y/N]: ", len(usernames))
		var confirm string
		if _, err := fmt.Scanln(&confirm); err != nil {
			fmt.Println("Cancelled")
			return nil
		}
		if confirm != "y" && confirm != "Y" {
			fmt.Println("Cancelled")
			return nil
		}
	}

	fmt.Printf("\n👥 Batch unfollowing %d users...\n", len(usernames))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	opts := &gh_client.BatchOptions{
		Concurrency: c.Int("concurrency"),
		RateLimit:   time.Duration(c.Int("rate-limit")) * time.Millisecond,
		DryRun:      dryRun,
		Progress:    true,
	}

	bp := gh_client.NewBatchProcessor(gc, opts)
	results := bp.BatchUnfollow(context.Background(), usernames, opts)

	// Display results
	gh_client.DisplayBatchResults(results, "Unfollow")

	// Update storage
	if !dryRun {
		list, err := st.Load()
		if err != nil {
			return fmt.Errorf("failed to load follow list: %w", err)
		}
		for _, r := range results {
			if r.Success {
				list.Remove(r.Username)
			}
		}
		if err := st.Save(list); err != nil {
			return fmt.Errorf("failed to save follow list: %w", err)
		}
	}

	return nil
}

// BatchCheck handles batch check command
func BatchCheck(c *cli.Context) error {
	if c.Args().Len() == 0 {
		return fmt.Errorf("please provide usernames")
	}

	token, err := gh_client.GetTokenFromGH()
	if err != nil {
		return fmt.Errorf("failed to get GitHub token: %w", err)
	}

	gc := gh_client.NewGitHubClient(token, "github.com")

	usernames := c.Args().Slice()
	usernames = uniqueStrings(usernames)

	fmt.Printf("\n🔍 Checking if %d users follow you...\n", len(usernames))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	opts := &gh_client.BatchOptions{
		Concurrency: c.Int("concurrency"),
		RateLimit:   time.Duration(c.Int("rate-limit")) * time.Millisecond,
	}

	bp := gh_client.NewBatchProcessor(gc, opts)
	results := bp.BatchCheckFollowers(context.Background(), usernames)

	following := make([]string, 0)
	notFollowing := make([]string, 0)

	for user, isFollowing := range results {
		if isFollowing {
			following = append(following, user)
		} else {
			notFollowing = append(notFollowing, user)
		}
	}

	fmt.Printf("\n✅ Following you (%d):\n", len(following))
	for _, u := range following {
		fmt.Printf("   %s\n", u)
	}

	fmt.Printf("\n❌ Not following you (%d):\n", len(notFollowing))
	for _, u := range notFollowing {
		fmt.Printf("   %s\n", u)
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	return nil
}

// BatchImport handles batch import from file
func BatchImport(c *cli.Context) error {
	if c.Args().Len() == 0 {
		return fmt.Errorf("please provide a file path")
	}

	filePath := c.Args().Get(0)
	usernames, err := readUsersFromFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if len(usernames) == 0 {
		return fmt.Errorf("no usernames found in file")
	}

	fmt.Printf("Found %d usernames in %s\n", len(usernames), filePath)

	if c.Bool("follow") {
		// Follow all imported users
		return batchFollowUsers(usernames, c)
	}

	// Just list them
	for _, u := range usernames {
		fmt.Println(u)
	}

	return nil
}

// batchFollowUsers follows a list of users
func batchFollowUsers(usernames []string, c *cli.Context) error {
	token, err := gh_client.GetTokenFromGH()
	if err != nil {
		return fmt.Errorf("failed to get GitHub token: %w", err)
	}

	gc := gh_client.NewGitHubClient(token, "github.com")
	st := storage.NewStorage(storage.DefaultStoragePath())

	dryRun := c.Bool("dry-run")
	if dryRun {
		fmt.Printf("\n⚠️  DRY RUN - No changes will be made\n")
	}

	fmt.Printf("\n👥 Following %d users...\n", len(usernames))

	opts := &gh_client.BatchOptions{
		Concurrency: c.Int("concurrency"),
		RateLimit:   time.Duration(c.Int("rate-limit")) * time.Millisecond,
		DryRun:      dryRun,
		Progress:    true,
	}

	bp := gh_client.NewBatchProcessor(gc, opts)
	results := bp.BatchFollow(context.Background(), usernames, opts)

	gh_client.DisplayBatchResults(results, "Follow")

	if !dryRun {
		list, err := st.Load()
		if err != nil {
			return fmt.Errorf("failed to load follow list: %w", err)
		}
		for _, r := range results {
			if r.Success {
				list.Add(r.Username, "", nil)
			}
		}
		if err := st.Save(list); err != nil {
			return fmt.Errorf("failed to save follow list: %w", err)
		}
	}

	return nil
}

// readUsersFromFile reads usernames from a file
func readUsersFromFile(filePath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse usernames (one per line)
	lines := strings.Split(string(data), "\n")
	usernames := make([]string, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			usernames = append(usernames, line)
		}
	}

	return usernames, nil
}

// uniqueStrings removes duplicates from a string slice
func uniqueStrings(s []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)

	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}

	return result
}
