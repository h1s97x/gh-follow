package config

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/h1s97x/gh-follow/internal/models"
)

// RunConfig runs the config command
func RunConfig(c *cli.Context) error {
	cm := NewConfigManager(DefaultConfigPath())

	// If no arguments, show current config
	if c.NArg() == 0 {
		return showConfig(cm)
	}

	action := c.Args().First()

	switch action {
	case "get":
		return getConfig(c, cm)
	case "set":
		return setConfig(c, cm)
	case "list":
		return showConfig(cm)
	case "reset":
		return resetConfig(cm)
	default:
		return fmt.Errorf("unknown config action: %s", action)
	}
}

// showConfig displays the current configuration
func showConfig(cm *ConfigManager) error {
	config, err := cm.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("\n⚙️  GH-Follow Configuration")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	fmt.Println("\n[Storage]")
	fmt.Printf("  local_path:    %s\n", config.Storage.LocalPath)
	fmt.Printf("  use_gist:      %v\n", config.Storage.UseGist)
	fmt.Printf("  gist_id:       %s\n", formatEmpty(config.Storage.GistID))

	fmt.Println("\n[Sync]")
	fmt.Printf("  auto_sync:     %v\n", config.Sync.AutoSync)
	fmt.Printf("  sync_interval: %d seconds\n", config.Sync.SyncInterval)

	if !config.Sync.LastSync.IsZero() {
		fmt.Printf("  last_sync:     %s\n", config.Sync.LastSync.Format("2006-01-02 15:04:05"))
	} else {
		fmt.Println("  last_sync:     Never")
	}

	fmt.Println("\n[Display]")
	fmt.Printf("  default_format: %s\n", config.Display.DefaultFormat)
	fmt.Printf("  default_sort:   %s\n", config.Display.DefaultSort)
	fmt.Printf("  default_order:  %s\n", config.Display.DefaultOrder)

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("\nCommands:")
	fmt.Println("  gh follow config get <key>              # Get a config value")
	fmt.Println("  gh follow config set <key> <value>      # Set a config value")
	fmt.Println("  gh follow config list                   # List all config values")
	fmt.Println("  gh follow config reset                  # Reset to defaults")

	return nil
}

// getConfig gets a configuration value
func getConfig(c *cli.Context, cm *ConfigManager) error {
	if c.NArg() < 2 {
		return fmt.Errorf("usage: gh follow config get <key>")
	}

	key := c.Args().Get(1)
	value, err := cm.Get(key)
	if err != nil {
		return err
	}

	fmt.Printf("%s: %v\n", key, value)
	return nil
}

// setConfig sets a configuration value
func setConfig(c *cli.Context, cm *ConfigManager) error {
	if c.NArg() < 3 {
		return fmt.Errorf("usage: gh follow config set <key> <value>")
	}

	key := c.Args().Get(1)
	value := c.Args().Get(2)

	// Parse the value based on the key
	var parsedValue interface{}
	switch key {
	case "sync.auto_sync", "storage.use_gist":
		parsedValue = value == "true"
	case "sync.sync_interval":
		var intVal int
		_, err := fmt.Sscanf(value, "%d", &intVal)
		if err != nil {
			return fmt.Errorf("invalid integer value: %s", value)
		}
		parsedValue = intVal
	default:
		parsedValue = value
	}

	if err := cm.Set(key, parsedValue); err != nil {
		return err
	}

	fmt.Printf("✅ Set %s = %v\n", key, parsedValue)
	return nil
}

// resetConfig resets configuration to defaults
func resetConfig(cm *ConfigManager) error {
	defaultConfig := models.NewDefaultConfig()

	if err := cm.Save(defaultConfig); err != nil {
		return err
	}

	fmt.Println("✅ Configuration reset to defaults")
	return nil
}

// formatEmpty formats an empty string
func formatEmpty(s string) string {
	if s == "" {
		return "(not set)"
	}
	return s
}
