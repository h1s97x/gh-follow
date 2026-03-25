package sync

import (
	"testing"
	"time"

	"github.com/h1s97x/gh-follow/internal/models"
)

func TestDetectConflicts(t *testing.T) {
	// Create test lists with same timestamps
	now := time.Now()

	t.Run("DifferentNotes", func(t *testing.T) {
		localList := models.NewFollowList()
		localList.Follows = []models.Follow{
			{Username: "user1", FollowedAt: now, Notes: "Local note", Tags: []string{"tag1"}},
			{Username: "user2", FollowedAt: now, Notes: "Same note", Tags: []string{"tag2"}},
		}
		localList.UpdatedAt = now

		remoteList := models.NewFollowList()
		remoteList.Follows = []models.Follow{
			{Username: "user1", FollowedAt: now, Notes: "Remote note", Tags: []string{"tag1"}},
			{Username: "user2", FollowedAt: now, Notes: "Same note", Tags: []string{"tag2"}},
		}
		remoteList.UpdatedAt = now.Add(-time.Hour)

		conflicts := detectConflicts(localList, remoteList)
		// user1 has different notes which is a conflict
		if len(conflicts) != 1 {
			t.Errorf("Expected 1 conflict, got %d: %v", len(conflicts), conflicts)
		}
	})

	t.Run("UserOnlyInLocal", func(t *testing.T) {
		local := models.NewFollowList()
		local.Follows = []models.Follow{
			{Username: "unique_local", FollowedAt: now, Notes: ""},
		}
		local.UpdatedAt = now

		remote := models.NewFollowList()
		remote.Follows = []models.Follow{}
		remote.UpdatedAt = now.Add(-time.Hour) // Remote was updated before local added user

		conflicts := detectConflicts(local, remote)
		// Should not be a conflict if remote wasn't modified after local
		if len(conflicts) != 0 {
			t.Errorf("Expected 0 conflicts, got %d: %v", len(conflicts), conflicts)
		}
	})

	t.Run("UserOnlyInRemote", func(t *testing.T) {
		local := models.NewFollowList()
		local.Follows = []models.Follow{}
		local.UpdatedAt = now

		remote := models.NewFollowList()
		remote.Follows = []models.Follow{
			{Username: "unique_remote", FollowedAt: now, Notes: ""},
		}
		remote.UpdatedAt = now

		conflicts := detectConflicts(local, remote)
		// Should not be a conflict if local wasn't modified after remote added user
		if len(conflicts) != 0 {
			t.Errorf("Expected 0 conflicts, got %d: %v", len(conflicts), conflicts)
		}
	})

	t.Run("RemovedFromRemote", func(t *testing.T) {
		now := time.Now()

		local := models.NewFollowList()
		local.Follows = []models.Follow{
			{Username: "user1", FollowedAt: now.Add(-time.Hour), Notes: ""},
		}
		local.UpdatedAt = now.Add(-time.Hour)

		remote := models.NewFollowList()
		remote.Follows = []models.Follow{} // User was removed from remote
		remote.UpdatedAt = now             // Remote was updated after local

		conflicts := detectConflicts(local, remote)
		// This IS a conflict - user removed from remote after local was modified
		if len(conflicts) == 0 {
			t.Error("Expected conflict when user removed from remote after local modification")
		}
	})
}

func TestMergeLists(t *testing.T) {
	t.Run("LocalWins", func(t *testing.T) {
		local := models.NewFollowList()
		local.Add("user1", "Local note", nil)

		remote := models.NewFollowList()
		remote.Add("user1", "Remote note", nil)
		remote.Add("user2", "User 2", nil)

		merged := mergeLists(local, remote, "local-wins")

		// Local note should be preserved
		if merged.Find("user1").Notes != "Local note" {
			t.Error("Expected local note to win")
		}

		// Remote-only user should be included
		if !merged.Contains("user2") {
			t.Error("Expected remote-only user to be included")
		}
	})

	t.Run("RemoteWins", func(t *testing.T) {
		local := models.NewFollowList()
		local.Add("user1", "Local note", nil)

		remote := models.NewFollowList()
		remote.Add("user1", "Remote note", nil)
		remote.Add("user2", "User 2", nil)

		merged := mergeLists(local, remote, "remote-wins")

		// Remote note should win
		if merged.Find("user1").Notes != "Remote note" {
			t.Error("Expected remote note to win")
		}

		// Local-only user should be included
		if !merged.Contains("user1") {
			t.Error("Expected local user to be included")
		}
	})

	t.Run("NewestWins", func(t *testing.T) {
		local := models.NewFollowList()
		local.Add("user1", "Local note", nil)

		remote := models.NewFollowList()
		remote.Add("user1", "Remote note", nil)
		remote.Add("user2", "User 2", nil)

		merged := mergeLists(local, remote, "newest-wins")

		// Both users should be present
		if !merged.Contains("user1") || !merged.Contains("user2") {
			t.Error("Expected both users to be present")
		}

		// Total count should be 2
		if len(merged.Follows) != 2 {
			t.Errorf("Expected 2 users, got %d", len(merged.Follows))
		}
	})
}

func TestGistSyncCreation(t *testing.T) {
	// Test that GistSync can be created
	gs := NewGistSync(nil, "test-gist-id")

	if gs.GetGistID() != "test-gist-id" {
		t.Errorf("Expected gist ID 'test-gist-id', got '%s'", gs.GetGistID())
	}

	if gs.filename != "gh-follow-list.json" {
		t.Errorf("Expected filename 'gh-follow-list.json', got '%s'", gs.filename)
	}
}

func TestAutoSyncCreation(t *testing.T) {
	as := NewAutoSync()

	if as.config == nil {
		t.Error("Expected config manager to be initialized")
	}

	if as.storage == nil {
		t.Error("Expected storage to be initialized")
	}
}

func TestConflictResolutionStrategies(t *testing.T) {
	// Setup test data
	now := time.Now()
	earlier := now.Add(-time.Hour)

	t.Run("LocalWinsPreservesLocalData", func(t *testing.T) {
		local := &models.FollowList{
			Follows: []models.Follow{
				{Username: "test", FollowedAt: now, Notes: "local"},
			},
		}

		remote := &models.FollowList{
			Follows: []models.Follow{
				{Username: "test", FollowedAt: earlier, Notes: "remote"},
			},
		}

		merged := mergeLists(local, remote, "local-wins")
		if merged.Find("test").Notes != "local" {
			t.Error("Local wins strategy should preserve local data")
		}
	})

	t.Run("RemoteWinsPreservesRemoteData", func(t *testing.T) {
		local := &models.FollowList{
			Follows: []models.Follow{
				{Username: "test", FollowedAt: now, Notes: "local"},
			},
		}

		remote := &models.FollowList{
			Follows: []models.Follow{
				{Username: "test", FollowedAt: earlier, Notes: "remote"},
			},
		}

		merged := mergeLists(local, remote, "remote-wins")
		if merged.Find("test").Notes != "remote" {
			t.Error("Remote wins strategy should preserve remote data")
		}
	})

	t.Run("NewestWinsUsesNewerTimestamp", func(t *testing.T) {
		local := &models.FollowList{
			Follows: []models.Follow{
				{Username: "test", FollowedAt: now, Notes: "newer"},
			},
		}

		remote := &models.FollowList{
			Follows: []models.Follow{
				{Username: "test", FollowedAt: earlier, Notes: "older"},
			},
		}

		merged := mergeLists(local, remote, "newest-wins")
		if merged.Find("test").Notes != "newer" {
			t.Error("Newest wins should use the entry with newer timestamp")
		}
	})
}
