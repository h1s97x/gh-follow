package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/h1s97x/gh-follow/internal/errors"
	"github.com/h1s97x/gh-follow/internal/models"
)

func TestStorageNew(t *testing.T) {
	st := NewStorage("/tmp/test-follow-list.json")
	if st.Path() != "/tmp/test-follow-list.json" {
		t.Errorf("Unexpected path: %s", st.Path())
	}
}

func TestStorageHomeExpansion(t *testing.T) {
	st := NewStorage("~/test-follow-list.json")
	// Just check it doesn't start with ~
	if len(st.Path()) > 0 && st.Path()[0] == '~' {
		t.Error("Expected home directory to be expanded")
	}
}

func TestStorageLoadNonExistent(t *testing.T) {
	// Use a temp file that doesn't exist
	tmpDir := t.TempDir()
	st := NewStorage(filepath.Join(tmpDir, "nonexistent.json"))

	list, err := st.Load()
	if err != nil {
		t.Fatalf("Expected no error for non-existent file, got: %v", err)
	}

	if list == nil {
		t.Fatal("Expected non-nil list")
	}

	if len(list.Follows) != 0 {
		t.Error("Expected empty list for non-existent file")
	}
}

func TestStorageSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	st := NewStorage(filepath.Join(tmpDir, "test-follow-list.json"))

	// Create and save a list
	list := models.NewFollowList()
	list.Add("octocat", "GitHub mascot", []string{"developer"})

	if err := st.Save(list); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Check file was created
	if !st.Exists() {
		t.Error("Expected file to exist after save")
	}

	// Load and verify
	loaded, err := st.Load()
	if err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	if len(loaded.Follows) != 1 {
		t.Errorf("Expected 1 follow, got %d", len(loaded.Follows))
	}

	if !loaded.Contains("octocat") {
		t.Error("Expected to contain octocat")
	}
}

func TestStorageAdd(t *testing.T) {
	tmpDir := t.TempDir()
	st := NewStorage(filepath.Join(tmpDir, "test-add.json"))

	// Add a user
	if err := st.Add("octocat", "Test note", []string{"tag1"}); err != nil {
		t.Fatalf("Failed to add: %v", err)
	}

	// Verify
	list, err := st.Load()
	if err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	if !list.Contains("octocat") {
		t.Error("Expected to contain octocat")
	}

	// Test duplicate
	if err := st.Add("octocat", "", nil); err != errors.ErrUserAlreadyFollowed {
		t.Errorf("Expected ErrUserAlreadyFollowed, got: %v", err)
	}
}

func TestStorageRemove(t *testing.T) {
	tmpDir := t.TempDir()
	st := NewStorage(filepath.Join(tmpDir, "test-remove.json"))

	// Add then remove
	if err := st.Add("octocat", "", nil); err != nil {
		t.Fatalf("Failed to add: %v", err)
	}
	if err := st.Remove("octocat"); err != nil {
		t.Fatalf("Failed to remove: %v", err)
	}

	// Verify
	list, err := st.Load()
	if err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	if list.Contains("octocat") {
		t.Error("Expected octocat to be removed")
	}

	// Test remove non-existent
	if err := st.Remove("nonexistent"); err != errors.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got: %v", err)
	}
}

func TestStorageEmptyUsername(t *testing.T) {
	st := NewStorage("/tmp/test-empty.json")

	if err := st.Add("", "", nil); err != errors.ErrEmptyUsername {
		t.Errorf("Expected ErrEmptyUsername for add, got: %v", err)
	}

	if err := st.Remove(""); err != errors.ErrEmptyUsername {
		t.Errorf("Expected ErrEmptyUsername for remove, got: %v", err)
	}
}

func TestStorageDelete(t *testing.T) {
	tmpDir := t.TempDir()
	st := NewStorage(filepath.Join(tmpDir, "test-delete.json"))

	// Create and save
	list := models.NewFollowList()
	if err := st.Save(list); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Delete
	if err := st.Delete(); err != nil {
		t.Fatalf("Failed to delete: %v", err)
	}

	// Verify
	if st.Exists() {
		t.Error("Expected file to not exist after delete")
	}
}

func TestStorageExportImport(t *testing.T) {
	tmpDir := t.TempDir()
	st := NewStorage(filepath.Join(tmpDir, "test-export.json"))
	exportPath := filepath.Join(tmpDir, "exported.json")

	// Add some data
	if err := st.Add("octocat", "Note", []string{"tag1"}); err != nil {
		t.Fatalf("Failed to add octocat: %v", err)
	}
	if err := st.Add("torvalds", "", []string{"linux"}); err != nil {
		t.Fatalf("Failed to add torvalds: %v", err)
	}

	// Export
	if err := st.Export(exportPath, "json"); err != nil {
		t.Fatalf("Failed to export: %v", err)
	}

	// Verify export file exists
	if _, err := os.Stat(exportPath); os.IsNotExist(err) {
		t.Fatal("Export file was not created")
	}

	// Create new storage and import
	st2 := NewStorage(filepath.Join(tmpDir, "test-import.json"))
	if err := st2.Import(exportPath, false); err != nil {
		t.Fatalf("Failed to import: %v", err)
	}

	// Verify imported data
	list, err := st2.Load()
	if err != nil {
		t.Fatalf("Failed to load imported: %v", err)
	}

	if len(list.Follows) != 2 {
		t.Errorf("Expected 2 follows after import, got %d", len(list.Follows))
	}
}

func TestStorageFilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	st := NewStorage(filepath.Join(tmpDir, "test-perms.json"))

	// Save
	list := models.NewFollowList()
	if err := st.Save(list); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Check permissions
	info, err := os.Stat(st.Path())
	if err != nil {
		t.Fatalf("Failed to stat: %v", err)
	}

	// Check that file is not world-readable (0600)
	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected file permissions 0600, got %o", info.Mode().Perm())
	}
}
