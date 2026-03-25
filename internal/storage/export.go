package storage

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// Export handles the follow export command
func Export(c *cli.Context) error {
	// Get flags
	output := c.String("output")
	format := c.String("format")

	// Initialize storage
	storage := NewStorage(DefaultStoragePath())

	// Check if storage exists
	if !storage.Exists() {
		return fmt.Errorf("no follow list found, please add some users first")
	}

	// Export
	if err := storage.Export(output, format); err != nil {
		return fmt.Errorf("failed to export follow list: %w", err)
	}

	fmt.Printf("Exported follow list to %s\n", output)
	return nil
}
