package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/h1s97x/gh-follow/internal/models"
)

// ConfigManager manages the configuration
type ConfigManager struct {
	configPath string
}

// DefaultConfigPath returns the default config path
func DefaultConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "gh-follow", "config.json")
}

// NewConfigManager creates a new ConfigManager
func NewConfigManager(path string) *ConfigManager {
	return &ConfigManager{
		configPath: path,
	}
}

// Load loads the configuration
func (cm *ConfigManager) Load() (*models.Config, error) {
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config
			return models.NewDefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config models.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// Save saves the configuration
func (cm *ConfigManager) Save(config *models.Config) error {
	// Ensure directory exists
	dir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// Update updates the configuration using a function
func (cm *ConfigManager) Update(fn func(*models.Config)) error {
	config, err := cm.Load()
	if err != nil {
		return err
	}

	fn(config)

	return cm.Save(config)
}

// Get gets a configuration value by path
func (cm *ConfigManager) Get(path string) (interface{}, error) {
	config, err := cm.Load()
	if err != nil {
		return nil, err
	}

	// Simple path-based access
	switch path {
	case "sync.auto_sync":
		return config.Sync.AutoSync, nil
	case "sync.sync_interval":
		return config.Sync.SyncInterval, nil
	case "storage.gist_id":
		return config.Storage.GistID, nil
	case "storage.use_gist":
		return config.Storage.UseGist, nil
	case "sync.last_sync":
		return config.Sync.LastSync, nil
	case "display.default_format":
		return config.Display.DefaultFormat, nil
	case "display.default_sort":
		return config.Display.DefaultSort, nil
	case "display.default_order":
		return config.Display.DefaultOrder, nil
	default:
		return nil, fmt.Errorf("unknown config path: %s", path)
	}
}

// Set sets a configuration value by path
func (cm *ConfigManager) Set(path string, value interface{}) error {
	return cm.Update(func(c *models.Config) {
		switch path {
		case "sync.auto_sync":
			if v, ok := value.(bool); ok {
				c.Sync.AutoSync = v
			}
		case "sync.sync_interval":
			if v, ok := value.(int); ok {
				c.Sync.SyncInterval = v
			}
		case "storage.gist_id":
			if v, ok := value.(string); ok {
				c.Storage.GistID = v
			}
		case "storage.use_gist":
			if v, ok := value.(bool); ok {
				c.Storage.UseGist = v
			}
		case "sync.last_sync":
			if v, ok := value.(time.Time); ok {
				c.Sync.LastSync = v
			}
		case "display.default_format":
			if v, ok := value.(string); ok {
				c.Display.DefaultFormat = v
			}
		case "display.default_sort":
			if v, ok := value.(string); ok {
				c.Display.DefaultSort = v
			}
		case "display.default_order":
			if v, ok := value.(string); ok {
				c.Display.DefaultOrder = v
			}
		}
	})
}

// GetGistID gets the Gist ID for sync
func (cm *ConfigManager) GetGistID() string {
	config, err := cm.Load()
	if err != nil {
		return ""
	}
	return config.Storage.GistID
}

// SetGistID sets the Gist ID for sync
func (cm *ConfigManager) SetGistID(id string) error {
	return cm.Set("storage.gist_id", id)
}
