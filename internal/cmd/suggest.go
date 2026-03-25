package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	gh_client "github.com/h1s97x/gh-follow/internal/github"
	"github.com/h1s97x/gh-follow/internal/storage"
	"github.com/h1s97x/gh-follow/internal/suggest"
)

// Suggest handles the suggest command
func Suggest(c *cli.Context) error {
	return SuggestGenerate(c)
}

// SuggestGenerate generates follow suggestions
func SuggestGenerate(c *cli.Context) error {
	token, err := gh_client.GetTokenFromGH()
	if err != nil {
		return fmt.Errorf("failed to get GitHub token: %w", err)
	}

	gc := gh_client.NewGitHubClient(token, "github.com")
	st := storage.NewStorage(storage.DefaultStoragePath())
	engine := suggest.NewSuggestionEngine(gc, st)

	opts := &suggest.SuggestionOptions{
		Limit:          c.Int("limit"),
		IncludeReasons: true,
	}

	if opts.Limit == 0 {
		opts.Limit = 20
	}

	format := c.String("format")

	fmt.Println("\n🔍 Generating follow suggestions...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	ctx := context.Background()
	suggestions, err := engine.GenerateSuggestions(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to generate suggestions: %w", err)
	}

	if len(suggestions) == 0 {
		fmt.Println("No suggestions found. Try following more users first!")
		return nil
	}

	switch format {
	case "json":
		for _, s := range suggestions {
			fmt.Printf("%s|%s|%d|%s\n", s.Username, s.Reason, s.Score, s.Source)
		}
	default:
		displaySuggestions(suggestions)
	}

	return nil
}

// displaySuggestions displays suggestions in a formatted table
func displaySuggestions(suggestions []*suggest.Suggestion) {
	fmt.Printf("\n📋 Suggested Users (%d)\n", len(suggestions))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for i, s := range suggestions {
		fmt.Printf("\n%d. \033[1;36m%s\033[0m", i+1, s.Username)
		if s.Name != "" {
			fmt.Printf(" - %s", s.Name)
		}
		fmt.Println()

		fmt.Printf("   Reason: %s\n", s.Reason)
		fmt.Printf("   Score:  %d | Source: %s", s.Score, s.Source)
		if s.Followers > 0 {
			fmt.Printf(" | Followers: %d", s.Followers)
		}
		fmt.Println()

		if len(s.FollowedBy) > 0 {
			fmt.Printf("   Followed by: %s", s.FollowedBy[0])
			if len(s.FollowedBy) > 1 {
				fmt.Printf(" and %d others", len(s.FollowedBy)-1)
			}
			fmt.Println()
		}

		if len(s.NotableRepos) > 0 {
			fmt.Printf("   Notable repos: %s", s.NotableRepos[0])
			if len(s.NotableRepos) > 1 {
				fmt.Printf(", %s", s.NotableRepos[1])
			}
			fmt.Println()
		}

		if s.Bio != "" && len(s.Bio) > 60 {
			fmt.Printf("   Bio: %s...\n", s.Bio[:60])
		} else if s.Bio != "" {
			fmt.Printf("   Bio: %s\n", s.Bio)
		}
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Use 'gh follow add <username>' to follow a suggested user\n")
}

// SuggestTrending shows trending users (placeholder)
func SuggestTrending(c *cli.Context) error {
	fmt.Println("\n🔥 Trending Users feature requires additional API implementation")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("\nThis feature would show trending GitHub users.")
	fmt.Println("Use 'gh follow suggest' for personalized suggestions based on your network.")
	return nil
}

// SuggestMutual shows mutual follow suggestions
func SuggestMutual(c *cli.Context) error {
	token, err := gh_client.GetTokenFromGH()
	if err != nil {
		return fmt.Errorf("failed to get GitHub token: %w", err)
	}

	gc := gh_client.NewGitHubClient(token, "github.com")
	st := storage.NewStorage(storage.DefaultStoragePath())

	list, err := st.Load()
	if err != nil {
		return fmt.Errorf("failed to load follow list: %w", err)
	}

	fmt.Println("\n🤝 Checking mutual followers...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	ctx := context.Background()
	mutuals := make([]string, 0)
	notFollowingBack := make([]string, 0)

	for _, f := range list.Follows {
		// Check if they follow you back
		isFollowing, err := gc.IsFollowing(ctx, f.Username)
		if err != nil {
			continue
		}

		if isFollowing {
			mutuals = append(mutuals, f.Username)
		} else {
			notFollowingBack = append(notFollowingBack, f.Username)
		}
	}

	fmt.Printf("\n✅ Mutual followers (%d):\n", len(mutuals))
	if len(mutuals) > 0 {
		for _, u := range mutuals {
			fmt.Printf("   %s\n", u)
		}
	} else {
		fmt.Println("   None")
	}

	fmt.Printf("\n❌ Not following you back (%d):\n", len(notFollowingBack))
	if len(notFollowingBack) > 0 {
		for _, u := range notFollowingBack {
			fmt.Printf("   %s\n", u)
		}
	} else {
		fmt.Println("   None")
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	return nil
}

// SuggestInactive shows suggestions for inactive followed users
func SuggestInactive(c *cli.Context) error {
	_, err := gh_client.GetTokenFromGH()
	if err != nil {
		return fmt.Errorf("failed to get GitHub token: %w", err)
	}

	st := storage.NewStorage(storage.DefaultStoragePath())

	list, err := st.Load()
	if err != nil {
		return fmt.Errorf("failed to load follow list: %w", err)
	}

	days := c.Int("days")
	if days == 0 {
		days = 365
	}

	fmt.Printf("\n😴 Inactive users detection requires additional implementation\n")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Would check for users inactive for >%d days\n", days)
	fmt.Printf("Total users to check: %d\n", len(list.Follows))

	return nil
}

