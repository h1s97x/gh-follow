package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/h1s97x/gh-follow/internal/config"
	"github.com/h1s97x/gh-follow/internal/github"
	"github.com/h1s97x/gh-follow/internal/models"
	"github.com/h1s97x/gh-follow/internal/storage"
)

// AutoSync handles automatic synchronization
type AutoSync struct {
	config   *config.ConfigManager
	storage  *storage.Storage
	client   *github.GitHubClient
	enabled  bool
	interval time.Duration
}

// NewAutoSync creates a new AutoSync instance
func NewAutoSync() *AutoSync {
	cm := config.NewConfigManager(config.DefaultConfigPath())
	cfg, _ := cm.Load()

	return &AutoSync{
		config:   cm,
		storage:  storage.NewStorage(storage.DefaultStoragePath()),
		enabled:  cfg.Sync.AutoSync,
		interval: time.Duration(cfg.Sync.SyncInterval) * time.Second,
	}
}

// ShouldSync checks if auto sync should be performed
func (as *AutoSync) ShouldSync() bool {
	if !as.enabled {
		return false
	}

	cfg, err := as.config.Load()
	if err != nil {
		return false
	}

	// Check if enough time has passed since last sync
	if !cfg.Sync.LastSync.IsZero() {
		elapsed := time.Since(cfg.Sync.LastSync)
		if elapsed < as.interval {
			return false
		}
	}

	return true
}

// PerformSync performs automatic synchronization
func (as *AutoSync) PerformSync(ctx context.Context, useGist bool) error {
	token, err := github.GetTokenFromGH()
	if err != nil {
		return fmt.Errorf("failed to get GitHub token: %w", err)
	}

	as.client = github.NewGitHubClient(token, "github.com")

	// Load local list
	localList, err := as.storage.Load()
	if err != nil {
		return fmt.Errorf("failed to load local list: %w", err)
	}

	if useGist {
		gistID := as.config.GetGistID()
		if gistID == "" {
			return fmt.Errorf("Gist not configured")
		}

		gs := NewGistSync(as.client, gistID)
		if err := gs.Upload(ctx, localList); err != nil {
			return fmt.Errorf("failed to sync to Gist: %w", err)
		}
	} else {
		// Sync with GitHub
		_, err := as.client.SyncFollowing(ctx, localList)
		if err != nil {
			return err
		}
	}

	// Update last sync time
	_ = as.config.Update(func(c *models.Config) {
		c.Sync.LastSync = time.Now()
	})

	return nil
}

// SyncAfterOperation performs sync after add/remove operations if auto-sync is enabled
func SyncAfterOperation(ctx context.Context, useGist bool, silent bool) error {
	as := NewAutoSync()

	if !as.ShouldSync() {
		return nil
	}

	if !silent {
		fmt.Println("🔄 Auto-syncing...")
	}

	if err := as.PerformSync(ctx, useGist); err != nil {
		if !silent {
			fmt.Printf("⚠️  Auto-sync failed: %v\n", err)
		}
		return err
	}

	if !silent {
		fmt.Println("✅ Auto-sync complete")
	}

	return nil
}

// AutoSyncStatus shows the auto-sync status
func AutoSyncStatus(c *cli.Context) error {
	cm := config.NewConfigManager(config.DefaultConfigPath())
	cfg, err := cm.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("\n🔄 Auto-Sync Status")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if cfg.Sync.AutoSync {
		fmt.Println("Status: ✅ Enabled")
	} else {
		fmt.Println("Status: ❌ Disabled")
	}

	fmt.Printf("Interval: %d seconds\n", cfg.Sync.SyncInterval)

	if !cfg.Sync.LastSync.IsZero() {
		fmt.Printf("Last Sync: %s\n", cfg.Sync.LastSync.Format("2006-01-02 15:04:05"))
		nextSync := cfg.Sync.LastSync.Add(time.Duration(cfg.Sync.SyncInterval) * time.Second)
		fmt.Printf("Next Sync: %s\n", nextSync.Format("2006-01-02 15:04:05"))
	} else {
		fmt.Println("Last Sync: Never")
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("\nCommands:")
	fmt.Println("  gh follow config set sync.auto_sync true   # Enable auto-sync")
	fmt.Println("  gh follow config set sync.auto_sync false  # Disable auto-sync")
	fmt.Println("  gh follow config set sync.sync_interval 3600  # Set interval (seconds)")

	return nil
}

// TriggerSync manually triggers a sync operation
func TriggerSync(c *cli.Context) error {
	useGist := c.Bool("gist")
	silent := c.Bool("silent")

	ctx := context.Background()

	if !silent {
		fmt.Println("🔄 Triggering sync...")
	}

	if err := SyncAfterOperation(ctx, useGist, silent); err != nil {
		return err
	}

	return nil
}
