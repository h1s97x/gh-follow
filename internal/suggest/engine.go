package suggest

import (
	"context"
	"fmt"
	"sort"

	"github.com/h1s97x/gh-follow/internal/github"
	"github.com/h1s97x/gh-follow/internal/models"
	"github.com/h1s97x/gh-follow/internal/storage"
)

// Suggestion represents a suggested user to follow
type Suggestion struct {
	Username     string   `json:"username"`
	Name         string   `json:"name,omitempty"`
	Bio          string   `json:"bio,omitempty"`
	Reason       string   `json:"reason"`
	Score        int      `json:"score"`
	Source       string   `json:"source"`
	Followers    int      `json:"followers"`
	FollowedBy   []string `json:"followed_by,omitempty"`
	NotableRepos []string `json:"notable_repos,omitempty"`
}

// SuggestionOptions configures the suggestion generation
type SuggestionOptions struct {
	Limit          int
	MinScore       int
	IncludeReasons bool
}

// SuggestionEngine generates follow suggestions
type SuggestionEngine struct {
	client  *github.GitHubClient
	storage *storage.Storage
}

// NewSuggestionEngine creates a new suggestion engine
func NewSuggestionEngine(client *github.GitHubClient, store *storage.Storage) *SuggestionEngine {
	return &SuggestionEngine{
		client:  client,
		storage: store,
	}
}

// GenerateSuggestions generates follow suggestions based on various strategies
func (se *SuggestionEngine) GenerateSuggestions(ctx context.Context, opts *SuggestionOptions) ([]*Suggestion, error) {
	if opts == nil {
		opts = &SuggestionOptions{Limit: 20}
	}

	// Load current follow list
	list, err := se.storage.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load follow list: %w", err)
	}

	// Get who you're following
	following, err := se.client.GetFollowing(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}

	// Build a set of who you follow
	followingSet := make(map[string]bool)
	for _, u := range following {
		followingSet[u.GetLogin()] = true
	}

	// Collect suggestions from multiple sources
	suggestionMap := make(map[string]*Suggestion)

	// Strategy 1: Users followed by people you follow
	if err := se.suggestFromFollowing(ctx, list, followingSet, suggestionMap, opts); err != nil {
		// Log error but continue
		fmt.Printf("Warning: failed to generate suggestions from following: %v\n", err)
	}

	// Strategy 2: Users you interact with (stargazers, etc.)
	if err := se.suggestFromInteractions(ctx, list, followingSet, suggestionMap, opts); err != nil {
		// Log error but continue
		fmt.Printf("Warning: failed to generate suggestions from interactions: %v\n", err)
	}

	// Convert map to slice and sort by score
	suggestions := make([]*Suggestion, 0, len(suggestionMap))
	for _, s := range suggestionMap {
		if s.Score >= opts.MinScore {
			suggestions = append(suggestions, s)
		}
	}

	// Sort by score (descending)
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Score > suggestions[j].Score
	})

	// Apply limit
	if len(suggestions) > opts.Limit {
		suggestions = suggestions[:opts.Limit]
	}

	return suggestions, nil
}

// suggestFromFollowing generates suggestions based on who you follow
func (se *SuggestionEngine) suggestFromFollowing(ctx context.Context, list *models.FollowList, followingSet map[string]bool, suggestionMap map[string]*Suggestion, opts *SuggestionOptions) error {
	// Get a sample of who you follow
	count := 0
	maxToCheck := 10

	for _, entry := range list.Follows {
		if count >= maxToCheck {
			break
		}

		// Get who this user follows
		theirFollowing, err := se.client.GetFollowing(ctx)
		if err != nil {
			continue
		}

		for _, user := range theirFollowing {
			username := user.GetLogin()
			// Skip if already following
			if followingSet[username] {
				continue
			}

			// Add or update suggestion
			if s, exists := suggestionMap[username]; exists {
				s.Score += 1
				s.FollowedBy = append(s.FollowedBy, entry.Username)
			} else {
				suggestionMap[username] = &Suggestion{
					Username:   username,
					Reason:     "Followed by people you follow",
					Score:      1,
					Source:     "following-network",
					FollowedBy: []string{entry.Username},
				}
			}
		}

		count++
	}

	return nil
}

// suggestFromInteractions generates suggestions based on interactions
func (se *SuggestionEngine) suggestFromInteractions(ctx context.Context, list *models.FollowList, followingSet map[string]bool, suggestionMap map[string]*Suggestion, opts *SuggestionOptions) error {
	// This is a placeholder for interaction-based suggestions
	// In a real implementation, this would analyze stargazers, collaborators, etc.
	return nil
}

// GetUserDetails gets details for a suggested user
func (se *SuggestionEngine) GetUserDetails(ctx context.Context, username string) (*Suggestion, error) {
	user, err := se.client.GetUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	suggestion := &Suggestion{
		Username:  user.GetLogin(),
		Name:      user.GetName(),
		Bio:       user.GetBio(),
		Followers: user.GetFollowers(),
	}

	return suggestion, nil
}
