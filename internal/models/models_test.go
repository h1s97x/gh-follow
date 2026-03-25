package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewFollowList(t *testing.T) {
	list := NewFollowList()

	if list.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", list.Version)
	}

	if len(list.Follows) != 0 {
		t.Errorf("Expected empty follows list, got %d", len(list.Follows))
	}

	if list.Metadata.TotalCount != 0 {
		t.Errorf("Expected total count 0, got %d", list.Metadata.TotalCount)
	}
}

func TestFollowListAdd(t *testing.T) {
	list := NewFollowList()

	list.Add("octocat", "GitHub mascot", []string{"developer"})

	if len(list.Follows) != 1 {
		t.Errorf("Expected 1 follow, got %d", len(list.Follows))
	}

	if list.Metadata.TotalCount != 1 {
		t.Errorf("Expected total count 1, got %d", list.Metadata.TotalCount)
	}

	// Test duplicate prevention
	list.Add("octocat", "", nil)
	if len(list.Follows) != 1 {
		t.Errorf("Expected 1 follow after duplicate add, got %d", len(list.Follows))
	}
}

func TestFollowListRemove(t *testing.T) {
	list := NewFollowList()
	list.Add("octocat", "", nil)
	list.Add("torvalds", "", nil)

	removed := list.Remove("octocat")
	if !removed {
		t.Error("Expected octocat to be removed")
	}

	if len(list.Follows) != 1 {
		t.Errorf("Expected 1 follow after removal, got %d", len(list.Follows))
	}

	// Test removing non-existent user
	removed = list.Remove("nonexistent")
	if removed {
		t.Error("Expected false when removing non-existent user")
	}
}

func TestFollowListContains(t *testing.T) {
	list := NewFollowList()
	list.Add("octocat", "", nil)

	if !list.Contains("octocat") {
		t.Error("Expected to contain octocat")
	}

	if list.Contains("nonexistent") {
		t.Error("Expected not to contain nonexistent")
	}
}

func TestFollowListFind(t *testing.T) {
	list := NewFollowList()
	list.Add("octocat", "GitHub mascot", []string{"developer"})

	found := list.Find("octocat")
	if found == nil {
		t.Fatal("Expected to find octocat")
	}

	if found.Notes != "GitHub mascot" {
		t.Errorf("Expected notes 'GitHub mascot', got '%s'", found.Notes)
	}

	notFound := list.Find("nonexistent")
	if notFound != nil {
		t.Error("Expected nil for nonexistent user")
	}
}

func TestFollowListGetStats(t *testing.T) {
	list := NewFollowList()
	list.Add("user1", "", nil)
	list.Add("user2", "", []string{"developer"})
	list.Add("user3", "", []string{"developer", "golang"})

	stats := list.GetStats()

	if stats.TotalFollows != 3 {
		t.Errorf("Expected total follows 3, got %d", stats.TotalFollows)
	}

	if len(stats.RecentFollows) != 3 {
		t.Errorf("Expected 3 recent follows, got %d", len(stats.RecentFollows))
	}

	if len(stats.Tags) != 2 {
		t.Errorf("Expected 2 unique tags, got %d", len(stats.Tags))
	}
}

func TestFollowListJSON(t *testing.T) {
	list := NewFollowList()
	list.Add("octocat", "GitHub mascot", []string{"developer"})

	data, err := json.Marshal(list)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var unmarshaled FollowList
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if unmarshaled.Version != list.Version {
		t.Errorf("Version mismatch after JSON roundtrip")
	}

	if len(unmarshaled.Follows) != 1 {
		t.Errorf("Expected 1 follow after unmarshal, got %d", len(unmarshaled.Follows))
	}
}

func TestNewDefaultConfig(t *testing.T) {
	config := NewDefaultConfig()

	if config.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", config.Version)
	}

	if config.Storage.UseGist {
		t.Error("Expected UseGist to be false by default")
	}

	if config.Display.DefaultFormat != "table" {
		t.Errorf("Expected default format 'table', got '%s'", config.Display.DefaultFormat)
	}
}

func TestFollowTime(t *testing.T) {
	list := NewFollowList()
	before := time.Now()
	list.Add("octocat", "", nil)
	after := time.Now()

	found := list.Find("octocat")
	if found.FollowedAt.Before(before) || found.FollowedAt.After(after) {
		t.Errorf("FollowedAt time is out of expected range: %v", found.FollowedAt)
	}
}
