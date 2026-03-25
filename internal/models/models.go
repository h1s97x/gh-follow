package models

import "time"

// Follow represents a single follow record
type Follow struct {
	Username  string    `json:"username"`
	FollowedAt time.Time `json:"followed_at"`
	Notes     string    `json:"notes,omitempty"`
	Tags      []string  `json:"tags,omitempty"`
}

// FollowList represents the complete follow list with metadata
type FollowList struct {
	Version   string    `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
	Follows   []Follow  `json:"follows"`
	Metadata  Metadata  `json:"metadata"`
}

// Metadata contains statistics about the follow list
type Metadata struct {
	TotalCount int       `json:"total_count"`
	LastSync   time.Time `json:"last_sync"`
	SyncStatus string    `json:"sync_status"`
}

// Config represents the application configuration
type Config struct {
	Version string      `json:"version"`
	Storage StorageConf `json:"storage"`
	Sync    SyncConf    `json:"sync"`
	Display DisplayConf `json:"display"`
}

// StorageConf contains storage-related configuration
type StorageConf struct {
	LocalPath string `json:"local_path"`
	UseGist   bool   `json:"use_gist"`
	GistID    string `json:"gist_id"`
}

// SyncConf contains sync-related configuration
type SyncConf struct {
	AutoSync     bool      `json:"auto_sync"`
	SyncInterval int       `json:"sync_interval"`
	LastSync     time.Time `json:"last_sync"`
}

// DisplayConf contains display-related configuration
type DisplayConf struct {
	DefaultFormat string `json:"default_format"`
	DefaultSort   string `json:"default_sort"`
	DefaultOrder  string `json:"default_order"`
}

// UserInfo represents basic GitHub user information
type UserInfo struct {
	Login     string    `json:"login"`
	ID        int64     `json:"id"`
	AvatarURL string    `json:"avatar_url"`
	Type      string    `json:"type"`
	Name      string    `json:"name,omitempty"`
	Bio       string    `json:"bio,omitempty"`
	Blog      string    `json:"blog,omitempty"`
	Location  string    `json:"location,omitempty"`
}

// Stats represents follow statistics
type Stats struct {
	TotalFollows  int       `json:"total_follows"`
	RecentFollows []Follow  `json:"recent_follows"`
	OldestFollow  *Follow   `json:"oldest_follow,omitempty"`
	Tags          []string  `json:"tags,omitempty"`
	LastUpdated   time.Time `json:"last_updated"`
}

// NewFollowList creates a new empty follow list with default values
func NewFollowList() *FollowList {
	return &FollowList{
		Version:   "1.0.0",
		UpdatedAt: time.Now(),
		Follows:   []Follow{},
		Metadata: Metadata{
			TotalCount: 0,
			SyncStatus: "never",
		},
	}
}

// NewDefaultConfig creates a new config with default values
func NewDefaultConfig() *Config {
	return &Config{
		Version: "1.0.0",
		Storage: StorageConf{
			LocalPath: "~/.config/gh/follow-list.json",
			UseGist:   false,
			GistID:    "",
		},
		Sync: SyncConf{
			AutoSync:     true,
			SyncInterval: 3600,
		},
		Display: DisplayConf{
			DefaultFormat: "table",
			DefaultSort:   "date",
			DefaultOrder:  "desc",
		},
	}
}

// Contains checks if a username is already in the follow list
func (fl *FollowList) Contains(username string) bool {
	for _, f := range fl.Follows {
		if f.Username == username {
			return true
		}
	}
	return false
}

// Find returns the follow record for a given username
func (fl *FollowList) Find(username string) *Follow {
	for i := range fl.Follows {
		if fl.Follows[i].Username == username {
			return &fl.Follows[i]
		}
	}
	return nil
}

// Add adds a new follow record
func (fl *FollowList) Add(username string, notes string, tags []string) {
	if fl.Contains(username) {
		return
	}

	follow := Follow{
		Username:  username,
		FollowedAt: time.Now(),
		Notes:     notes,
		Tags:      tags,
	}

	fl.Follows = append(fl.Follows, follow)
	fl.Metadata.TotalCount = len(fl.Follows)
	fl.UpdatedAt = time.Now()
}

// Remove removes a follow record by username
func (fl *FollowList) Remove(username string) bool {
	for i, f := range fl.Follows {
		if f.Username == username {
			fl.Follows = append(fl.Follows[:i], fl.Follows[i+1:]...)
			fl.Metadata.TotalCount = len(fl.Follows)
			fl.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

// GetStats returns statistics about the follow list
func (fl *FollowList) GetStats() *Stats {
	stats := &Stats{
		TotalFollows: len(fl.Follows),
		LastUpdated:  fl.UpdatedAt,
	}

	if len(fl.Follows) > 0 {
		// Get recent follows (last 5)
		recentCount := 5
		if len(fl.Follows) < recentCount {
			recentCount = len(fl.Follows)
		}
		stats.RecentFollows = fl.Follows[:recentCount]

		// Get oldest follow
		stats.OldestFollow = &fl.Follows[len(fl.Follows)-1]

		// Collect unique tags
		tagSet := make(map[string]bool)
		for _, f := range fl.Follows {
			for _, tag := range f.Tags {
				tagSet[tag] = true
			}
		}
		for tag := range tagSet {
			stats.Tags = append(stats.Tags, tag)
		}
	}

	return stats
}
