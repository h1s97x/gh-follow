package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/h1s97x/gh-follow/internal/errors"
	"github.com/h1s97x/gh-follow/internal/models"
)

// Storage handles local storage operations for the follow list
type Storage struct {
	path string
}

// NewStorage creates a new Storage instance
func NewStorage(path string) *Storage {
	// Expand home directory
	if strings.HasPrefix(path, "~/") {
		homeDir, _ := os.UserHomeDir()
		path = filepath.Join(homeDir, path[2:])
	}
	return &Storage{path: path}
}

// DefaultStoragePath returns the default storage file path
func DefaultStoragePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "gh", "follow-list.json")
}

// Load loads the follow list from disk
func (s *Storage) Load() (*models.FollowList, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return a new empty list if file doesn't exist
			return models.NewFollowList(), nil
		}
		return nil, errors.NewFollowError("load", "", err)
	}

	var list models.FollowList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, errors.NewFollowError("load", "", fmt.Errorf("failed to parse follow list: %w", err))
	}

	return &list, nil
}

// Save saves the follow list to disk
func (s *Storage) Save(list *models.FollowList) error {
	// Ensure directory exists
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.NewFollowError("save", "", fmt.Errorf("failed to create directory: %w", err))
	}

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return errors.NewFollowError("save", "", fmt.Errorf("failed to marshal follow list: %w", err))
	}

	// Write with secure permissions
	if err := os.WriteFile(s.path, data, 0600); err != nil {
		return errors.NewFollowError("save", "", fmt.Errorf("failed to write follow list: %w", err))
	}

	return nil
}

// Add adds a username to the follow list
func (s *Storage) Add(username string, notes string, tags []string) error {
	if username == "" {
		return errors.ErrEmptyUsername
	}

	list, err := s.Load()
	if err != nil {
		return err
	}

	if list.Contains(username) {
		return errors.ErrUserAlreadyFollowed
	}

	list.Add(username, notes, tags)
	return s.Save(list)
}

// Remove removes a username from the follow list
func (s *Storage) Remove(username string) error {
	if username == "" {
		return errors.ErrEmptyUsername
	}

	list, err := s.Load()
	if err != nil {
		return err
	}

	if !list.Remove(username) {
		return errors.ErrUserNotFound
	}

	return s.Save(list)
}

// Exists checks if the storage file exists
func (s *Storage) Exists() bool {
	_, err := os.Stat(s.path)
	return err == nil
}

// Path returns the storage file path
func (s *Storage) Path() string {
	return s.path
}

// Delete removes the storage file
func (s *Storage) Delete() error {
	if err := os.Remove(s.path); err != nil && !os.IsNotExist(err) {
		return errors.NewFollowError("delete", "", err)
	}
	return nil
}

// Export exports the follow list to a file
func (s *Storage) Export(outputPath string, format string) error {
	list, err := s.Load()
	if err != nil {
		return err
	}

	var data []byte
	switch format {
	case "json":
		data, err = json.MarshalIndent(list, "", "  ")
	case "csv":
		data, err = s.toCSV(list)
	default:
		return errors.ErrInvalidFormat
	}

	if err != nil {
		return errors.NewFollowError("export", "", err)
	}

	// Expand home directory in output path
	if strings.HasPrefix(outputPath, "~/") {
		homeDir, _ := os.UserHomeDir()
		outputPath = filepath.Join(homeDir, outputPath[2:])
	}

	return os.WriteFile(outputPath, data, 0644)
}

// Import imports a follow list from a file
func (s *Storage) Import(inputPath string, merge bool) error {
	// Expand home directory in input path
	if strings.HasPrefix(inputPath, "~/") {
		homeDir, _ := os.UserHomeDir()
		inputPath = filepath.Join(homeDir, inputPath[2:])
	}

	data, err := os.ReadFile(inputPath)
	if err != nil {
		return errors.NewFollowError("import", "", fmt.Errorf("failed to read import file: %w", err))
	}

	var importedList models.FollowList
	if err := json.Unmarshal(data, &importedList); err != nil {
		return errors.NewFollowError("import", "", fmt.Errorf("failed to parse import file: %w", err))
	}

	if merge {
		// Merge with existing list
		existingList, err := s.Load()
		if err != nil {
			return err
		}

		for _, f := range importedList.Follows {
			if !existingList.Contains(f.Username) {
				existingList.Follows = append(existingList.Follows, f)
			}
		}
		existingList.Metadata.TotalCount = len(existingList.Follows)
		existingList.UpdatedAt = s.now()

		return s.Save(existingList)
	}

	// Replace existing list
	importedList.UpdatedAt = s.now()
	return s.Save(&importedList)
}

// toCSV converts follow list to CSV format
func (s *Storage) toCSV(list *models.FollowList) ([]byte, error) {
	var csv strings.Builder
	csv.WriteString("username,followed_at,notes,tags\n")

	for _, f := range list.Follows {
		tags := strings.Join(f.Tags, ";")
		csv.WriteString(fmt.Sprintf("%s,%s,%s,%s\n",
			f.Username,
			f.FollowedAt.Format("2006-01-02T15:04:05Z"),
			f.Notes,
			tags,
		))
	}

	return []byte(csv.String()), nil
}

// now returns current time (extracted for testing)
func (s *Storage) now() time.Time {
	return time.Now()
}
