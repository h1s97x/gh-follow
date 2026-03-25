package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/h1s97x/gh-follow/internal/config"
	"github.com/h1s97x/gh-follow/internal/models"
	"github.com/h1s97x/gh-follow/internal/storage"
)

// TestAddAndListIntegration tests the full add and list workflow
func TestAddAndListIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	storagePath := filepath.Join(tmpDir, "test-integration.json")

	// Override default storage path
	st := storage.NewStorage(storagePath)

	// Test 1: Add users
	t.Run("AddUsers", func(t *testing.T) {
		err := st.Add("user1", "First user", []string{"developer"})
		if err != nil {
			t.Fatalf("Failed to add user1: %v", err)
		}

		err = st.Add("user2", "Second user", []string{"golang"})
		if err != nil {
			t.Fatalf("Failed to add user2: %v", err)
		}

		err = st.Add("user3", "", []string{"developer", "golang"})
		if err != nil {
			t.Fatalf("Failed to add user3: %v", err)
		}
	})

	// Test 2: List users
	t.Run("ListUsers", func(t *testing.T) {
		list, err := st.Load()
		if err != nil {
			t.Fatalf("Failed to load list: %v", err)
		}

		if len(list.Follows) != 3 {
			t.Errorf("Expected 3 users, got %d", len(list.Follows))
		}

		// Verify metadata
		if list.Metadata.TotalCount != 3 {
			t.Errorf("Expected total count 3, got %d", list.Metadata.TotalCount)
		}
	})

	// Test 3: Remove a user
	t.Run("RemoveUser", func(t *testing.T) {
		err := st.Remove("user2")
		if err != nil {
			t.Fatalf("Failed to remove user2: %v", err)
		}

		list, err := st.Load()
		if err != nil {
			t.Fatalf("Failed to load list: %v", err)
		}

		if len(list.Follows) != 2 {
			t.Errorf("Expected 2 users after removal, got %d", len(list.Follows))
		}

		if list.Contains("user2") {
			t.Error("user2 should have been removed")
		}
	})

	// Test 4: Export
	t.Run("ExportList", func(t *testing.T) {
		exportPath := filepath.Join(tmpDir, "exported.json")
		err := st.Export(exportPath, "json")
		if err != nil {
			t.Fatalf("Failed to export: %v", err)
		}

		if _, err := os.Stat(exportPath); os.IsNotExist(err) {
			t.Error("Export file was not created")
		}
	})

	// Test 5: Import
	t.Run("ImportList", func(t *testing.T) {
		importPath := filepath.Join(tmpDir, "import.json")
		importData := `{
			"version": "1.0.0",
			"follows": [
				{"username": "imported1", "followed_at": "2024-01-01T00:00:00Z"},
				{"username": "imported2", "followed_at": "2024-01-02T00:00:00Z"}
			]
		}`
		if err := os.WriteFile(importPath, []byte(importData), 0644); err != nil {
			t.Fatalf("Failed to write import file: %v", err)
		}

		// Import with merge
		err := st.Import(importPath, true)
		if err != nil {
			t.Fatalf("Failed to import: %v", err)
		}

		list, err := st.Load()
		if err != nil {
			t.Fatalf("Failed to load list: %v", err)
		}

		// Should have original 2 + imported 2 = 4
		if len(list.Follows) != 4 {
			t.Errorf("Expected 4 users after import, got %d", len(list.Follows))
		}
	})
}

// TestFilteringIntegration tests filtering functionality
func TestFilteringIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	storagePath := filepath.Join(tmpDir, "test-filter.json")
	st := storage.NewStorage(storagePath)

	// Setup test data
	if err := st.Add("gopher", "Go developer", []string{"golang", "developer"}); err != nil {
		t.Fatalf("Failed to add gopher: %v", err)
	}
	if err := st.Add("pythonista", "Python developer", []string{"python", "developer"}); err != nil {
		t.Fatalf("Failed to add pythonista: %v", err)
	}
	if err := st.Add("rustacean", "Rust developer", []string{"rust", "developer"}); err != nil {
		t.Fatalf("Failed to add rustacean: %v", err)
	}
	if err := st.Add("golang_org", "Official Go", []string{"golang", "official"}); err != nil {
		t.Fatalf("Failed to add golang_org: %v", err)
	}

	list, _ := st.Load()

	t.Run("FilterByUsername", func(t *testing.T) {
		filtered := applyFilters(list.Follows, "go", "", "", "")
		if len(filtered) != 2 {
			t.Errorf("Expected 2 users matching 'go', got %d", len(filtered))
		}
	})

	t.Run("FilterByTag", func(t *testing.T) {
		filtered := applyFilters(list.Follows, "", "golang", "", "")
		if len(filtered) != 2 {
			t.Errorf("Expected 2 users with 'golang' tag, got %d", len(filtered))
		}
	})

	t.Run("FilterByUsernameAndTag", func(t *testing.T) {
		filtered := applyFilters(list.Follows, "gopher", "golang", "", "")
		if len(filtered) != 1 {
			t.Errorf("Expected 1 user matching both filters, got %d", len(filtered))
		}
	})
}

// TestSortingIntegration tests sorting functionality
func TestSortingIntegration(t *testing.T) {
	list := models.NewFollowList()
	list.Add("charlie", "", nil)
	list.Add("alice", "", nil)
	list.Add("bob", "", nil)

	t.Run("SortByNameAsc", func(t *testing.T) {
		follows := make([]models.Follow, len(list.Follows))
		copy(follows, list.Follows)
		sortFollows(follows, "name", "asc")

		if follows[0].Username != "alice" {
			t.Errorf("Expected first user 'alice', got '%s'", follows[0].Username)
		}
	})

	t.Run("SortByNameDesc", func(t *testing.T) {
		follows := make([]models.Follow, len(list.Follows))
		copy(follows, list.Follows)
		sortFollows(follows, "name", "desc")

		if follows[0].Username != "charlie" {
			t.Errorf("Expected first user 'charlie', got '%s'", follows[0].Username)
		}
	})
}

// TestConfigIntegration tests configuration management
func TestConfigIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.json")

	cm := config.NewConfigManager(configPath)

	t.Run("LoadDefaultConfig", func(t *testing.T) {
		cfg, err := cm.Load()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		if cfg.Version != "1.0.0" {
			t.Errorf("Expected version 1.0.0, got %s", cfg.Version)
		}
	})

	t.Run("UpdateConfig", func(t *testing.T) {
		err := cm.Update(func(c *models.Config) {
			c.Display.DefaultFormat = "json"
			c.Sync.AutoSync = false
		})
		if err != nil {
			t.Fatalf("Failed to update config: %v", err)
		}

		cfg, _ := cm.Load()
		if cfg.Display.DefaultFormat != "json" {
			t.Errorf("Expected format 'json', got '%s'", cfg.Display.DefaultFormat)
		}
		if cfg.Sync.AutoSync {
			t.Error("Expected AutoSync to be false")
		}
	})
}
