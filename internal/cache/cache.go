package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// UserCache represents cached user information
type UserCache struct {
	Login       string    `json:"login"`
	ID          int64     `json:"id"`
	AvatarURL   string    `json:"avatar_url"`
	Name        string    `json:"name,omitempty"`
	Bio         string    `json:"bio,omitempty"`
	Blog        string    `json:"blog,omitempty"`
	Location    string    `json:"location,omitempty"`
	Company     string    `json:"company,omitempty"`
	PublicRepos int       `json:"public_repos"`
	Followers   int       `json:"followers"`
	Following   int       `json:"following"`
	CachedAt    time.Time `json:"cached_at"`
}

// CacheStats represents cache statistics
type CacheStats struct {
	TotalUsers int       `json:"total_users"`
	MaxSize    int       `json:"max_size"`
	TTL        string    `json:"ttl"`
	Oldest     time.Time `json:"oldest"`
	Newest     time.Time `json:"newest"`
}

// Cache represents a simple file-based cache
type Cache struct {
	path     string
	duration time.Duration
}

// Entry represents a cache entry
type Entry struct {
	Data      []byte    `json:"data"`
	ExpiresAt time.Time `json:"expires_at"`
}

// DefaultCachePath returns the default cache path
func DefaultCachePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cache", "gh-follow")
}

// NewCache creates a new Cache instance
func NewCache(path string, duration time.Duration) *Cache {
	if duration == 0 {
		duration = 1 * time.Hour
	}
	return &Cache{
		path:     path,
		duration: duration,
	}
}

// Get retrieves a user from the cache
func (c *Cache) Get(username string) *UserCache {
	filename := filepath.Join(c.path, username+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	var user UserCache
	if err := json.Unmarshal(data, &user); err != nil {
		return nil
	}

	// Check if cached data is too old
	if time.Since(user.CachedAt) > c.duration {
		return nil
	}

	return &user
}

// Set stores a user in the cache
func (c *Cache) Set(user *UserCache) error {
	// Ensure directory exists
	if err := os.MkdirAll(c.path, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	user.CachedAt = time.Now()
	data, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal user cache: %w", err)
	}

	filename := filepath.Join(c.path, user.Login+".json")
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache: %w", err)
	}

	return nil
}

// Delete removes a user from the cache
func (c *Cache) Delete(username string) error {
	filename := filepath.Join(c.path, username+".json")
	return os.Remove(filename)
}

// Clear clears all cache entries
func (c *Cache) Clear() error {
	return os.RemoveAll(c.path)
}

// Cleanup removes expired entries and returns count of removed entries
func (c *Cache) Cleanup() int {
	entries, err := os.ReadDir(c.path)
	if err != nil {
		return 0
	}

	removed := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := filepath.Join(c.path, entry.Name())
		data, err := os.ReadFile(filename)
		if err != nil {
			continue
		}

		var user UserCache
		if err := json.Unmarshal(data, &user); err != nil {
			continue
		}

		// Delete if expired
		if time.Since(user.CachedAt) > c.duration {
			os.Remove(filename)
			removed++
		}
	}

	return removed
}

// List returns all cached users
func (c *Cache) List() []*UserCache {
	entries, err := os.ReadDir(c.path)
	if err != nil {
		return nil
	}

	var users []*UserCache
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := filepath.Join(c.path, entry.Name())
		data, err := os.ReadFile(filename)
		if err != nil {
			continue
		}

		var user UserCache
		if err := json.Unmarshal(data, &user); err != nil {
			continue
		}

		users = append(users, &user)
	}

	return users
}

// GetStats returns cache statistics
func (c *Cache) GetStats() *CacheStats {
	users := c.List()

	stats := &CacheStats{
		TotalUsers: len(users),
		MaxSize:    1000,
		TTL:        c.duration.String(),
	}

	if len(users) > 0 {
		stats.Oldest = users[0].CachedAt
		stats.Newest = users[0].CachedAt

		for _, user := range users {
			if user.CachedAt.Before(stats.Oldest) {
				stats.Oldest = user.CachedAt
			}
			if user.CachedAt.After(stats.Newest) {
				stats.Newest = user.CachedAt
			}
		}
	}

	return stats
}

// Save saves the cache (no-op for file-based cache)
func (c *Cache) Save() error {
	// No-op for file-based cache
	return nil
}

// GetRaw retrieves raw data from the cache
func (c *Cache) GetRaw(key string) ([]byte, bool) {
	filename := filepath.Join(c.path, key+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, false
	}

	return data, true
}

// SetRaw stores raw data in the cache
func (c *Cache) SetRaw(key string, value []byte) error {
	// Ensure directory exists
	if err := os.MkdirAll(c.path, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	filename := filepath.Join(c.path, key+".json")
	if err := os.WriteFile(filename, value, 0644); err != nil {
		return fmt.Errorf("failed to write cache: %w", err)
	}

	return nil
}
