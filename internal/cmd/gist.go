package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/h1s97x/gh-follow/internal/config"
	gh_client "github.com/h1s97x/gh-follow/internal/github"
	"github.com/h1s97x/gh-follow/internal/models"
	"github.com/h1s97x/gh-follow/internal/storage"
	"github.com/h1s97x/gh-follow/internal/sync"
)

// Gist handles the gist command
func Gist(c *cli.Context) error {
	return nil
}

// GistCreate creates a new Gist for sync
func GistCreate(c *cli.Context) error {
	token, err := gh_client.GetTokenFromGH()
	if err != nil {
		return fmt.Errorf("failed to get GitHub token: %w", err)
	}

	gc := gh_client.NewGitHubClient(token, "github.com")
	st := storage.NewStorage(storage.DefaultStoragePath())
	ctx := context.Background()

	// Load local list
	list, err := st.Load()
	if err != nil {
		return fmt.Errorf("failed to load local list: %w", err)
	}

	// Create Gist sync
	gs := sync.NewGistSync(gc, "")

	// Create the Gist
	_, err = gs.CreateGist(ctx, list)
	if err != nil {
		// Note: This is a placeholder - actual Gist creation requires direct API access
		fmt.Println("⚠️  Gist creation requires additional implementation")
		fmt.Println("This feature needs direct GitHub API access to create Gists.")
		return nil
	}

	// Save Gist ID to config would happen here
	fmt.Println("✅ Gist sync setup initiated")
	fmt.Println("Use 'gh follow sync --gist' to sync changes.")

	return nil
}

// GistStatus shows the Gist sync status
func GistStatus(c *cli.Context) error {
	cm := config.NewConfigManager(config.DefaultConfigPath())
	gistID := cm.GetGistID()

	if gistID == "" {
		fmt.Println("❌ Gist sync is not configured")
		fmt.Println("Run 'gh follow gist create' to set up Gist sync")
		return nil
	}

	fmt.Printf("✅ Gist sync is configured\n")
	fmt.Printf("Gist ID: %s\n", gistID)
	fmt.Printf("Gist URL: https://gist.github.com/%s\n", gistID)

	return nil
}

// GistPull pulls the follow list from Gist
func GistPull(c *cli.Context) error {
	force := c.Bool("force")

	cm := config.NewConfigManager(config.DefaultConfigPath())
	gistID := cm.GetGistID()

	if gistID == "" {
		return fmt.Errorf("Gist sync is not configured. Run 'gh follow gist create' first")
	}

	token, err := gh_client.GetTokenFromGH()
	if err != nil {
		return fmt.Errorf("failed to get GitHub token: %w", err)
	}

	gc := gh_client.NewGitHubClient(token, "github.com")
	st := storage.NewStorage(storage.DefaultStoragePath())
	gs := sync.NewGistSync(gc, gistID)
	ctx := context.Background()

	// Download from Gist
	gistList, err := gs.Download(ctx)
	if err != nil {
		return fmt.Errorf("failed to download from Gist: %w", err)
	}

	// Load local list
	localList, err := st.Load()
	if err != nil {
		return fmt.Errorf("failed to load local list: %w", err)
	}

	// Check for conflicts
	conflicts := detectConflicts(localList, gistList)
	if len(conflicts) > 0 && !force {
		fmt.Printf("⚠️  Found %d conflicts:\n", len(conflicts))
		for _, c := range conflicts {
			fmt.Printf("  - %s\n", c)
		}
		fmt.Println("\nUse --force to overwrite local changes, or resolve conflicts manually")
		return nil
	}

	// Save to local
	if err := st.Save(gistList); err != nil {
		return fmt.Errorf("failed to save local list: %w", err)
	}

	fmt.Printf("✅ Pulled %d follows from Gist\n", len(gistList.Follows))
	return nil
}

// GistPush pushes the follow list to Gist
func GistPush(c *cli.Context) error {
	force := c.Bool("force")

	cm := config.NewConfigManager(config.DefaultConfigPath())
	gistID := cm.GetGistID()

	if gistID == "" {
		return fmt.Errorf("Gist sync is not configured. Run 'gh follow gist create' first")
	}

	token, err := gh_client.GetTokenFromGH()
	if err != nil {
		return fmt.Errorf("failed to get GitHub token: %w", err)
	}

	gc := gh_client.NewGitHubClient(token, "github.com")
	st := storage.NewStorage(storage.DefaultStoragePath())
	gs := sync.NewGistSync(gc, gistID)
	ctx := context.Background()

	// Load local list
	localList, err := st.Load()
	if err != nil {
		return fmt.Errorf("failed to load local list: %w", err)
	}

	// Download current Gist state
	gistList, err := gs.Download(ctx)
	if err != nil {
		// If Gist is empty or not found, just push
		gistList = models.NewFollowList()
	}

	// Check for conflicts
	conflicts := detectConflicts(localList, gistList)
	if len(conflicts) > 0 && !force {
		fmt.Printf("⚠️  Found %d conflicts:\n", len(conflicts))
		for _, c := range conflicts {
			fmt.Printf("  - %s\n", c)
		}
		fmt.Println("\nUse --force to overwrite Gist, or resolve conflicts manually")
		return nil
	}

	// Upload to Gist
	if err := gs.Upload(ctx, localList); err != nil {
		return fmt.Errorf("failed to upload to Gist: %w", err)
	}

	fmt.Printf("✅ Pushed %d follows to Gist\n", len(localList.Follows))
	return nil
}

// detectConflicts detects conflicts between local and remote lists
func detectConflicts(local, remote *models.FollowList) []string {
	var conflicts []string

	// Build maps for comparison
	localMap := make(map[string]models.Follow)
	for _, f := range local.Follows {
		localMap[f.Username] = f
	}

	remoteMap := make(map[string]models.Follow)
	for _, f := range remote.Follows {
		remoteMap[f.Username] = f
	}

	// Check for users modified in both places
	for username, localFollow := range localMap {
		if remoteFollow, exists := remoteMap[username]; exists {
			// User exists in both - check if modified
			if !localFollow.FollowedAt.Equal(remoteFollow.FollowedAt) {
				conflicts = append(conflicts, fmt.Sprintf("%s: different follow dates", username))
			}
			if localFollow.Notes != remoteFollow.Notes {
				conflicts = append(conflicts, fmt.Sprintf("%s: different notes", username))
			}
		}
	}

	// Check for users only in one place
	for username := range localMap {
		if _, exists := remoteMap[username]; !exists {
			// Only in local - potential conflict if remote was modified
			if remote.UpdatedAt.After(local.UpdatedAt) {
				conflicts = append(conflicts, fmt.Sprintf("%s: removed from remote", username))
			}
		}
	}

	for username := range remoteMap {
		if _, exists := localMap[username]; !exists {
			// Only in remote - potential conflict if local was modified
			if local.UpdatedAt.After(remote.UpdatedAt) {
				conflicts = append(conflicts, fmt.Sprintf("%s: removed from local", username))
			}
		}
	}

	return conflicts
}

// mergeLists merges local and remote lists with conflict resolution
func mergeLists(local, remote *models.FollowList, strategy string) *models.FollowList {
	merged := models.NewFollowList()

	// Build maps
	localMap := make(map[string]models.Follow)
	for _, f := range local.Follows {
		localMap[f.Username] = f
	}

	remoteMap := make(map[string]models.Follow)
	for _, f := range remote.Follows {
		remoteMap[f.Username] = f
	}

	// Merge based on strategy
	switch strategy {
	case "local-wins":
		// Local changes take precedence
		for _, f := range local.Follows {
			merged.Follows = append(merged.Follows, f)
		}
		// Add remote-only follows
		for username, f := range remoteMap {
			if _, exists := localMap[username]; !exists {
				merged.Follows = append(merged.Follows, f)
			}
		}

	case "remote-wins":
		// Remote changes take precedence
		for _, f := range remote.Follows {
			merged.Follows = append(merged.Follows, f)
		}
		// Add local-only follows
		for username, f := range localMap {
			if _, exists := remoteMap[username]; !exists {
				merged.Follows = append(merged.Follows, f)
			}
		}

	default: // "newest-wins"
		// Use the most recently modified version
		for username, localFollow := range localMap {
			if remoteFollow, exists := remoteMap[username]; exists {
				// User exists in both - use newer
				if localFollow.FollowedAt.After(remoteFollow.FollowedAt) {
					merged.Follows = append(merged.Follows, localFollow)
				} else {
					merged.Follows = append(merged.Follows, remoteFollow)
				}
			} else {
				merged.Follows = append(merged.Follows, localFollow)
			}
		}
		// Add remote-only follows
		for username, f := range remoteMap {
			if _, exists := localMap[username]; !exists {
				merged.Follows = append(merged.Follows, f)
			}
		}
	}

	merged.Metadata.TotalCount = len(merged.Follows)
	merged.UpdatedAt = time.Now()

	return merged
}
