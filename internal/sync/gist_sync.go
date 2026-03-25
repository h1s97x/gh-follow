package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	go_github "github.com/google/go-github/v55/github"

	"github.com/h1s97x/gh-follow/internal/github"
	"github.com/h1s97x/gh-follow/internal/models"
)

// GistSync handles synchronization with GitHub Gist
type GistSync struct {
	client   *github.GitHubClient
	gistID   string
	filename string
}

// NewGistSync creates a new GistSync instance
func NewGistSync(client *github.GitHubClient, gistID string) *GistSync {
	return &GistSync{
		client:   client,
		gistID:   gistID,
		filename: "gh-follow-list.json",
	}
}

// CreateGist creates a new Gist with the follow list
func (gs *GistSync) CreateGist(ctx context.Context, list *models.FollowList) (*go_github.Gist, error) {
	// Serialize the follow list
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal follow list: %w", err)
	}

	// Create Gist
	gist := &go_github.Gist{
		Description: go_github.String("GH-Follow sync data"),
		Public:      go_github.Bool(false), // Always private
		Files: map[go_github.GistFilename]go_github.GistFile{
			go_github.GistFilename(gs.filename): {
				Content: go_github.String(string(data)),
			},
		},
	}

	// Note: We need to access the underlying client
	// For now, return an error indicating this needs the raw client
	_ = gist // Silence unused variable error until we implement proper API access
	return nil, fmt.Errorf("CreateGist requires direct API access - use sync manager")
}

// Download downloads the follow list from Gist
func (gs *GistSync) Download(ctx context.Context) (*models.FollowList, error) {
	if gs.gistID == "" {
		return nil, fmt.Errorf("Gist ID not configured")
	}

	// This is a simplified version - in real implementation,
	// you would use the GitHub client to fetch the gist
	return nil, fmt.Errorf("Download requires direct API access - use sync manager")
}

// Upload uploads the follow list to Gist
func (gs *GistSync) Upload(ctx context.Context, list *models.FollowList) error {
	// Update timestamp
	list.UpdatedAt = time.Now()

	// This is a simplified version
	return fmt.Errorf("Upload requires direct API access - use sync manager")
}

// GetGistID returns the current Gist ID
func (gs *GistSync) GetGistID() string {
	return gs.gistID
}

// SetGistID sets the Gist ID
func (gs *GistSync) SetGistID(id string) {
	gs.gistID = id
}

// detectConflicts detects conflicts between local and remote lists
func detectConflicts(local, remote *models.FollowList) []string {
	var conflicts []string

	// Build maps for comparison
	localMap := make(map[string]models.Follow)
	for _, f := range local.Follows {
		localMap[f.Username] = f
	}

	remoteMap := make(map[string]models.Follow)
	for _, f := range remote.Follows {
		remoteMap[f.Username] = f
	}

	// Check for users modified in both places
	for username, localFollow := range localMap {
		if remoteFollow, exists := remoteMap[username]; exists {
			// User exists in both - check if modified
			if !localFollow.FollowedAt.Equal(remoteFollow.FollowedAt) {
				conflicts = append(conflicts, fmt.Sprintf("%s: different follow dates", username))
			}
			if localFollow.Notes != remoteFollow.Notes {
				conflicts = append(conflicts, fmt.Sprintf("%s: different notes", username))
			}
		}
	}

	// Check for users only in one place
	for username := range localMap {
		if _, exists := remoteMap[username]; !exists {
			// Only in local - potential conflict if remote was modified
			if remote.UpdatedAt.After(local.UpdatedAt) {
				conflicts = append(conflicts, fmt.Sprintf("%s: removed from remote", username))
			}
		}
	}

	for username := range remoteMap {
		if _, exists := localMap[username]; !exists {
			// Only in remote - potential conflict if local was modified
			if local.UpdatedAt.After(remote.UpdatedAt) {
				conflicts = append(conflicts, fmt.Sprintf("%s: removed from local", username))
			}
		}
	}

	return conflicts
}

// mergeLists merges local and remote lists with conflict resolution
func mergeLists(local, remote *models.FollowList, strategy string) *models.FollowList {
	merged := models.NewFollowList()

	// Build maps
	localMap := make(map[string]models.Follow)
	for _, f := range local.Follows {
		localMap[f.Username] = f
	}

	remoteMap := make(map[string]models.Follow)
	for _, f := range remote.Follows {
		remoteMap[f.Username] = f
	}

	// Merge based on strategy
	switch strategy {
	case "local-wins":
		// Local changes take precedence
		for _, f := range local.Follows {
			merged.Follows = append(merged.Follows, f)
		}
		// Add remote-only follows
		for username, f := range remoteMap {
			if _, exists := localMap[username]; !exists {
				merged.Follows = append(merged.Follows, f)
			}
		}

	case "remote-wins":
		// Remote changes take precedence
		for _, f := range remote.Follows {
			merged.Follows = append(merged.Follows, f)
		}
		// Add local-only follows
		for username, f := range localMap {
			if _, exists := remoteMap[username]; !exists {
				merged.Follows = append(merged.Follows, f)
			}
		}

	default: // "newest-wins"
		// Use the most recently modified version
		for username, localFollow := range localMap {
			if remoteFollow, exists := remoteMap[username]; exists {
				// User exists in both - use newer
				if localFollow.FollowedAt.After(remoteFollow.FollowedAt) {
					merged.Follows = append(merged.Follows, localFollow)
				} else {
					merged.Follows = append(merged.Follows, remoteFollow)
				}
			} else {
				merged.Follows = append(merged.Follows, localFollow)
			}
		}
		// Add remote-only follows
		for username, f := range remoteMap {
			if _, exists := localMap[username]; !exists {
				merged.Follows = append(merged.Follows, f)
			}
		}
	}

	merged.Metadata.TotalCount = len(merged.Follows)
	merged.UpdatedAt = time.Now()

	return merged
}
