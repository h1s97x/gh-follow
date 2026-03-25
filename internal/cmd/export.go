package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/h1s97x/gh-follow/internal/storage"
	"github.com/h1s97x/gh-follow/internal/sync"
)

// Export handles the follow export command
func Export(c *cli.Context) error {
	// Get flags
	output := c.String("output")
	format := c.String("format")

	// Initialize storage
	st := storage.NewStorage(storage.DefaultStoragePath())

	// Check if storage exists
	if !st.Exists() {
		return fmt.Errorf("no follow list found, please add some users first")
	}

	// Export
	if err := st.Export(output, format); err != nil {
		return fmt.Errorf("failed to export follow list: %w", err)
	}

	fmt.Printf("Exported follow list to %s\n", output)
	return nil
}

// Import handles the follow import command (wrapper for storage.Import)
func Import(c *cli.Context) error {
	return storage.Import(c)
}

// AutoSyncStatus shows the auto-sync status
func AutoSyncStatus(c *cli.Context) error {
	return sync.AutoSyncStatus(c)
}

// TriggerSync manually triggers a sync operation
func TriggerSync(c *cli.Context) error {
	return sync.TriggerSync(c)
}
