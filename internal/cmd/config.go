package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/h1s97x/gh-follow/internal/config"
	"github.com/h1s97x/gh-follow/internal/models"
)

// ConfigCmd handles the config command
func ConfigCmd(c *cli.Context) error {
	// Get subcommand
	if c.Args().Len() == 0 {
		return showConfig(c)
	}

	return nil
}

// ConfigShow shows the current configuration
func ConfigShow(c *cli.Context) error {
	cm := config.NewConfigManager(config.DefaultConfigPath())
	cfg, err := cm.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}

// ConfigSet sets a configuration value
func ConfigSet(c *cli.Context) error {
	if c.Args().Len() < 2 {
		return fmt.Errorf("usage: gh follow config set <key> <value>")
	}

	key := c.Args().Get(0)
	value := c.Args().Get(1)

	cm := config.NewConfigManager(config.DefaultConfigPath())

	err := cm.Update(func(cfg *models.Config) {
		switch key {
		case "storage.local_path":
			cfg.Storage.LocalPath = value
		case "storage.use_gist":
			cfg.Storage.UseGist = value == "true"
		case "storage.gist_id":
			cfg.Storage.GistID = value
		case "sync.auto_sync":
			cfg.Sync.AutoSync = value == "true"
		case "sync.sync_interval":
			var interval int
			if _, err := fmt.Sscanf(value, "%d", &interval); err != nil {
				fmt.Printf("Warning: invalid interval value: %v\n", err)
				return
			}
			cfg.Sync.SyncInterval = interval
		case "display.default_format":
			cfg.Display.DefaultFormat = value
		case "display.default_sort":
			cfg.Display.DefaultSort = value
		case "display.default_order":
			cfg.Display.DefaultOrder = value
		default:
			fmt.Printf("Unknown config key: %s\n", key)
			fmt.Println("Available keys:")
			fmt.Println("  storage.local_path")
			fmt.Println("  storage.use_gist")
			fmt.Println("  storage.gist_id")
			fmt.Println("  sync.auto_sync")
			fmt.Println("  sync.sync_interval")
			fmt.Println("  display.default_format")
			fmt.Println("  display.default_sort")
			fmt.Println("  display.default_order")
		}
	})

	if err != nil {
		return fmt.Errorf("failed to set config: %w", err)
	}

	fmt.Printf("Set %s = %s\n", key, value)
	return nil
}

// ConfigGet gets a configuration value
func ConfigGet(c *cli.Context) error {
	if c.Args().Len() < 1 {
		return fmt.Errorf("usage: gh follow config get <key>")
	}

	key := c.Args().Get(0)
	cm := config.NewConfigManager(config.DefaultConfigPath())
	cfg, err := cm.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	var value interface{}
	switch key {
	case "storage.local_path":
		value = cfg.Storage.LocalPath
	case "storage.use_gist":
		value = cfg.Storage.UseGist
	case "storage.gist_id":
		value = cfg.Storage.GistID
	case "sync.auto_sync":
		value = cfg.Sync.AutoSync
	case "sync.sync_interval":
		value = cfg.Sync.SyncInterval
	case "display.default_format":
		value = cfg.Display.DefaultFormat
	case "display.default_sort":
		value = cfg.Display.DefaultSort
	case "display.default_order":
		value = cfg.Display.DefaultOrder
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	fmt.Printf("%s: %v\n", key, value)
	return nil
}

// ConfigReset resets the configuration to defaults
func ConfigReset(c *cli.Context) error {
	force := c.Bool("force")

	if !force {
		fmt.Print("Are you sure you want to reset configuration to defaults? [y/N]: ")
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

	cm := config.NewConfigManager(config.DefaultConfigPath())
	defaultConfig := models.NewDefaultConfig()
	if err := cm.Save(defaultConfig); err != nil {
		return fmt.Errorf("failed to reset config: %w", err)
	}

	fmt.Println("Configuration reset to defaults")
	return nil
}

// showConfig shows the current configuration (default action)
func showConfig(c *cli.Context) error {
	cm := config.NewConfigManager(config.DefaultConfigPath())
	cfg, err := cm.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("\n📋 GH-Follow Configuration")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("\nStorage:")
	fmt.Printf("  Local Path:    %s\n", cfg.Storage.LocalPath)
	fmt.Printf("  Use Gist:      %v\n", cfg.Storage.UseGist)
	if cfg.Storage.GistID != "" {
		fmt.Printf("  Gist ID:       %s\n", cfg.Storage.GistID)
	}

	fmt.Println("\nSync:")
	fmt.Printf("  Auto Sync:     %v\n", cfg.Sync.AutoSync)
	fmt.Printf("  Interval:      %d seconds\n", cfg.Sync.SyncInterval)

	fmt.Println("\nDisplay:")
	fmt.Printf("  Format:        %s\n", cfg.Display.DefaultFormat)
	fmt.Printf("  Sort:          %s\n", cfg.Display.DefaultSort)
	fmt.Printf("  Order:         %s\n", cfg.Display.DefaultOrder)

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Config file: %s\n", config.DefaultConfigPath())

	return nil
}
