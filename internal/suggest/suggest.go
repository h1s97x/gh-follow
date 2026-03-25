package suggest

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/h1s97x/gh-follow/internal/github"
	"github.com/h1s97x/gh-follow/internal/storage"
)

// RunSuggestions runs the suggestion command
func RunSuggestions(c *cli.Context) error {
	// Get GitHub token
	token, err := github.GetTokenFromGH()
	if err != nil {
		return fmt.Errorf("failed to get GitHub token: %w", err)
	}

	// Create clients
	client := github.NewGitHubClient(token, "github.com")
	store := storage.NewStorage(storage.DefaultStoragePath())

	// Create suggestion engine
	engine := NewSuggestionEngine(client, store)

	// Generate suggestions
	ctx := c.Context
	opts := &SuggestionOptions{
		Limit:          c.Int("limit"),
		MinScore:       c.Int("min-score"),
		IncludeReasons: c.Bool("reasons"),
	}

	if opts.Limit == 0 {
		opts.Limit = 20
	}

	suggestions, err := engine.GenerateSuggestions(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to generate suggestions: %w", err)
	}

	// Display suggestions
	if len(suggestions) == 0 {
		fmt.Println("No suggestions found.")
		return nil
	}

	fmt.Printf("\n💡 Suggested Users to Follow (%d)\n", len(suggestions))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, s := range suggestions {
		fmt.Printf("\n@%s", s.Username)
		if s.Name != "" {
			fmt.Printf(" (%s)", s.Name)
		}
		fmt.Println()

		if s.Bio != "" {
			fmt.Printf("  %s\n", s.Bio)
		}

		fmt.Printf("  Score: %d | Followers: %d\n", s.Score, s.Followers)

		if opts.IncludeReasons && s.Reason != "" {
			fmt.Printf("  Reason: %s\n", s.Reason)
		}

		if len(s.FollowedBy) > 0 {
			fmt.Printf("  Followed by: %s\n", strings.Join(s.FollowedBy, ", "))
		}
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("\nCommands:")
	fmt.Println("  gh follow suggest --limit 10     # Show top 10 suggestions")
	fmt.Println("  gh follow suggest --min-score 3  # Only show high-score suggestions")
	fmt.Println("  gh follow suggest --reasons      # Include reasons")

	return nil
}

// InteractiveSuggestion allows interactive selection of suggestions
func InteractiveSuggestion(c *cli.Context) error {
	// Get GitHub token
	token, err := github.GetTokenFromGH()
	if err != nil {
		return fmt.Errorf("failed to get GitHub token: %w", err)
	}

	// Create clients
	client := github.NewGitHubClient(token, "github.com")
	store := storage.NewStorage(storage.DefaultStoragePath())

	// Create suggestion engine
	engine := NewSuggestionEngine(client, store)

	// Generate suggestions
	ctx := c.Context
	opts := &SuggestionOptions{
		Limit:    10,
		MinScore: 2,
	}

	suggestions, err := engine.GenerateSuggestions(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to generate suggestions: %w", err)
	}

	if len(suggestions) == 0 {
		fmt.Println("No suggestions found.")
		return nil
	}

	// Display suggestions for selection
	fmt.Printf("\n💡 Suggested Users to Follow (%d)\n", len(suggestions))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for i, s := range suggestions {
		fmt.Printf("%d. @%s", i+1, s.Username)
		if s.Name != "" {
			fmt.Printf(" (%s)", s.Name)
		}
		fmt.Println()
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Use fzf for selection if available
	if _, err := exec.LookPath("fzf"); err == nil {
		return selectWithFzf(suggestions, client)
	}

	// Fallback to simple selection
	fmt.Println("\nTo follow a user: gh follow add <username>")

	return nil
}

// selectWithFzf uses fzf for interactive selection
func selectWithFzf(suggestions []*Suggestion, client *github.GitHubClient) error {
	// Build input for fzf
	var input strings.Builder
	for _, s := range suggestions {
		input.WriteString(s.Username)
		input.WriteString("\n")
	}

	// Run fzf
	cmd := exec.Command("fzf", "--prompt=Select users to follow: ")
	cmd.Stdin = strings.NewReader(input.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = strings.NewReader(input.String())

	if err := cmd.Run(); err != nil {
		// User cancelled
		return nil
	}

	return nil
}
