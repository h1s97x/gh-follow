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

// Sync handles the follow sync command
func Sync(c *cli.Context) error {
	// Get flags
	direction := c.String("direction")
	useGist := c.Bool("gist")
	dryRun := c.Bool("dry-run")
	force := c.Bool("force")

	// Initialize storage
	st := storage.NewStorage(storage.DefaultStoragePath())

	// Get GitHub token
	token, err := gh_client.GetTokenFromGH()
	if err != nil {
		return fmt.Errorf("failed to get GitHub token: %w", err)
	}

	gc := gh_client.NewGitHubClient(token, c.String("hostname"))
	ctx := context.Background()

	// Load local list
	localList, err := st.Load()
	if err != nil {
		return fmt.Errorf("failed to load local list: %w", err)
	}

	if dryRun {
		fmt.Println("🔍 Dry run mode - no changes will be made")
	}

	// Handle Gist sync if enabled
	if useGist {
		return syncWithGist(ctx, gc, st, localList, direction, dryRun, force, c.String("gist-id"))
	}

	// Regular GitHub sync
	switch direction {
	case "pull":
		return syncPull(ctx, gc, st, localList, dryRun)
	case "push":
		return syncPush(ctx, gc, st, localList, dryRun, force)
	case "both":
		fallthrough
	default:
		// Pull first, then push
		if err := syncPull(ctx, gc, st, localList, dryRun); err != nil {
			return err
		}
		return syncPush(ctx, gc, st, localList, dryRun, force)
	}
}

// syncWithGist handles syncing with Gist
func syncWithGist(ctx context.Context, gc *gh_client.GitHubClient, st *storage.Storage, localList *models.FollowList, direction string, dryRun bool, force bool, gistID string) error {
	// Get or create Gist ID
	cm := config.NewConfigManager(config.DefaultConfigPath())
	if gistID == "" {
		gistID = cm.GetGistID()
	}

	if gistID == "" {
		fmt.Println("⚠️  No Gist configured. Creating one...")
		gs := sync.NewGistSync(gc, "")
		_, err := gs.CreateGist(ctx, localList)
		if err != nil {
			fmt.Println("⚠️  Gist creation requires additional implementation")
			return nil
		}
		fmt.Println("✅ Created Gist sync placeholder")
		return nil
	}

	gs := sync.NewGistSync(gc, gistID)

	switch direction {
	case "pull":
		return syncPullFromGist(ctx, gs, st, localList, dryRun, force)
	case "push":
		return syncPushToGist(ctx, gs, st, localList, dryRun, force)
	case "both":
		fallthrough
	default:
		// Pull then push
		if err := syncPullFromGist(ctx, gs, st, localList, dryRun, force); err != nil {
			return err
		}
		return syncPushToGist(ctx, gs, st, localList, dryRun, force)
	}
}

// syncPullFromGist pulls the follow list from Gist
func syncPullFromGist(ctx context.Context, gs *sync.GistSync, st *storage.Storage, localList *models.FollowList, dryRun bool, force bool) error {
	fmt.Println("📥 Pulling from Gist...")

	// Download from Gist
	gistList, err := gs.Download(ctx)
	if err != nil {
		return fmt.Errorf("failed to download from Gist: %w", err)
	}

	// Detect conflicts
	conflicts := detectConflicts(localList, gistList)
	if len(conflicts) > 0 && !force && !dryRun {
		fmt.Printf("⚠️  Found %d conflicts:\n", len(conflicts))
		for _, c := range conflicts {
			fmt.Printf("  - %s\n", c)
		}
		fmt.Println("\nUse --force to overwrite local changes")
		return nil
	}

	if dryRun {
		fmt.Printf("Would pull %d follows from Gist\n", len(gistList.Follows))
		if len(conflicts) > 0 {
			fmt.Printf("Would resolve %d conflicts\n", len(conflicts))
		}
		return nil
	}

	// Merge and save
	merged := mergeLists(localList, gistList, "newest-wins")
	if err := st.Save(merged); err != nil {
		return fmt.Errorf("failed to save local list: %w", err)
	}

	fmt.Printf("✅ Pulled %d follows from Gist\n", len(gistList.Follows))
	return nil
}

// syncPushToGist pushes the follow list to Gist
func syncPushToGist(ctx context.Context, gs *sync.GistSync, st *storage.Storage, localList *models.FollowList, dryRun bool, force bool) error {
	fmt.Println("📤 Pushing to Gist...")

	// Download current Gist state for comparison
	gistList, err := gs.Download(ctx)
	if err != nil {
		// Gist might be empty, continue with push
		gistList = models.NewFollowList()
	}

	// Detect conflicts
	conflicts := detectConflicts(localList, gistList)
	if len(conflicts) > 0 && !force && !dryRun {
		fmt.Printf("⚠️  Found %d conflicts:\n", len(conflicts))
		for _, c := range conflicts {
			fmt.Printf("  - %s\n", c)
		}
		fmt.Println("\nUse --force to overwrite Gist")
		return nil
	}

	if dryRun {
		fmt.Printf("Would push %d follows to Gist\n", len(localList.Follows))
		return nil
	}

	// Update timestamp
	localList.UpdatedAt = time.Now()
	localList.Metadata.LastSync = time.Now()
	localList.Metadata.SyncStatus = "success"

	// Upload to Gist
	if err := gs.Upload(ctx, localList); err != nil {
		return fmt.Errorf("failed to upload to Gist: %w", err)
	}

	fmt.Printf("✅ Pushed %d follows to Gist\n", len(localList.Follows))
	return nil
}

// syncPull pulls the following list from GitHub to local storage
func syncPull(ctx context.Context, gc *gh_client.GitHubClient, st *storage.Storage, localList *models.FollowList, dryRun bool) error {
	fmt.Println("📥 Pulling following list from GitHub...")

	// Get following list from GitHub
	githubFollowing, err := gc.GetFollowing(ctx)
	if err != nil {
		return fmt.Errorf("failed to get following list from GitHub: %w", err)
	}

	// Create maps for comparison
	githubUsers := make(map[string]bool)
	for _, user := range githubFollowing {
		githubUsers[user.GetLogin()] = true
	}

	localUsers := make(map[string]bool)
	for _, f := range localList.Follows {
		localUsers[f.Username] = true
	}

	// Find users to add (on GitHub but not in local)
	var toAdd []string
	for username := range githubUsers {
		if !localUsers[username] {
			toAdd = append(toAdd, username)
		}
	}

	// Find users to remove (in local but not on GitHub)
	var toRemove []string
	for username := range localUsers {
		if !githubUsers[username] {
			toRemove = append(toRemove, username)
		}
	}

	fmt.Printf("Found %d users to add, %d users to remove\n", len(toAdd), len(toRemove))

	if dryRun {
		if len(toAdd) > 0 {
			fmt.Printf("Would add: %v\n", toAdd)
		}
		if len(toRemove) > 0 {
			fmt.Printf("Would remove: %v\n", toRemove)
		}
		return nil
	}

	// Add new users
	for _, username := range toAdd {
		if err := st.Add(username, "", nil); err != nil {
			fmt.Printf("⚠️  Failed to add %s: %v\n", username, err)
		} else {
			fmt.Printf("✅ Added %s from GitHub\n", username)
		}
	}

	// Remove users not on GitHub
	for _, username := range toRemove {
		if err := st.Remove(username); err != nil {
			fmt.Printf("⚠️  Failed to remove %s: %v\n", username, err)
		} else {
			fmt.Printf("🗑️  Removed %s (not following on GitHub)\n", username)
		}
	}

	fmt.Printf("\n✅ Sync complete: %d added, %d removed\n", len(toAdd), len(toRemove))
	return nil
}

// syncPush pushes the local follow list to GitHub
func syncPush(ctx context.Context, gc *gh_client.GitHubClient, st *storage.Storage, localList *models.FollowList, dryRun bool, force bool) error {
	fmt.Println("📤 Pushing follow list to GitHub...")

	// Get current GitHub following
	githubFollowing, err := gc.GetFollowing(ctx)
	if err != nil {
		return fmt.Errorf("failed to get following list from GitHub: %w", err)
	}

	// Create map for quick lookup
	githubUsers := make(map[string]bool)
	for _, user := range githubFollowing {
		githubUsers[user.GetLogin()] = true
	}

	// Find users to follow
	var toFollow []string
	for _, f := range localList.Follows {
		if !githubUsers[f.Username] {
			toFollow = append(toFollow, f.Username)
		}
	}

	fmt.Printf("Found %d users to follow on GitHub\n", len(toFollow))

	if dryRun {
		if len(toFollow) > 0 {
			fmt.Printf("Would follow: %v\n", toFollow)
		}
		return nil
	}

	// Follow users
	for _, username := range toFollow {
		if err := gc.Follow(ctx, username); err != nil {
			fmt.Printf("⚠️  Failed to follow %s: %v\n", username, err)
		} else {
			fmt.Printf("✅ Followed %s on GitHub\n", username)
		}
	}

	fmt.Printf("\n✅ Push complete: %d users followed\n", len(toFollow))
	return nil
}
