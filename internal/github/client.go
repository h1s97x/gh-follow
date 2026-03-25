package github

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/google/go-github/v55/github"

	"github.com/h1s97x/gh-follow/internal/errors"
	"github.com/h1s97x/gh-follow/internal/models"
)

// GitHubClient wraps the GitHub API client with follow-specific operations
type GitHubClient struct {
	client   *github.Client
	token    string
	hostname string
}

// NewGitHubClient creates a new GitHub API client
func NewGitHubClient(token string, hostname string) *GitHubClient {
	if hostname == "" {
		hostname = "github.com"
	}

	client := github.NewClient(nil)
	if hostname != "github.com" {
		// For GitHub Enterprise
		baseURL := fmt.Sprintf("https://%s/api/v3/", hostname)
		client, _ = github.NewClient(nil).WithEnterpriseURLs(baseURL, baseURL)
	}

	if token != "" {
		client = client.WithAuthToken(token)
	}

	return &GitHubClient{
		client:   client,
		token:    token,
		hostname: hostname,
	}
}

// GetTokenFromGH retrieves the GitHub token from gh CLI
func GetTokenFromGH() (string, error) {
	cmd := exec.Command("gh", "auth", "token")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.ErrNoToken
	}
	return strings.TrimSpace(string(output)), nil
}

// Follow follows a GitHub user
func (gc *GitHubClient) Follow(ctx context.Context, username string) error {
	if username == "" {
		return errors.ErrEmptyUsername
	}

	resp, err := gc.client.Users.Follow(ctx, username)
	if err != nil {
		return gc.handleAPIError(err, resp, "follow", username)
	}

	return nil
}

// Unfollow unfollows a GitHub user
func (gc *GitHubClient) Unfollow(ctx context.Context, username string) error {
	if username == "" {
		return errors.ErrEmptyUsername
	}

	resp, err := gc.client.Users.Unfollow(ctx, username)
	if err != nil {
		return gc.handleAPIError(err, resp, "unfollow", username)
	}

	return nil
}

// IsFollowing checks if the authenticated user is following a specific user
func (gc *GitHubClient) IsFollowing(ctx context.Context, username string) (bool, error) {
	if username == "" {
		return false, errors.ErrEmptyUsername
	}

	isFollowing, resp, err := gc.client.Users.IsFollowing(ctx, "", username)
	if err != nil {
		return false, gc.handleAPIError(err, resp, "check following", username)
	}

	return isFollowing, nil
}

// GetFollowing retrieves all users that the authenticated user is following
func (gc *GitHubClient) GetFollowing(ctx context.Context) ([]*github.User, error) {
	var allUsers []*github.User
	opts := &github.ListOptions{PerPage: 100}

	for {
		users, resp, err := gc.client.Users.ListFollowing(ctx, "", opts)
		if err != nil {
			return nil, gc.handleAPIError(err, resp, "get following", "")
		}

		allUsers = append(allUsers, users...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage

		// Rate limiting: wait between requests
		time.Sleep(100 * time.Millisecond)
	}

	return allUsers, nil
}

// GetFollowers retrieves all followers of the authenticated user
func (gc *GitHubClient) GetFollowers(ctx context.Context) ([]*github.User, error) {
	var allUsers []*github.User
	opts := &github.ListOptions{PerPage: 100}

	for {
		users, resp, err := gc.client.Users.ListFollowers(ctx, "", opts)
		if err != nil {
			return nil, gc.handleAPIError(err, resp, "get followers", "")
		}

		allUsers = append(allUsers, users...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage

		// Rate limiting: wait between requests
		time.Sleep(100 * time.Millisecond)
	}

	return allUsers, nil
}

// GetUser retrieves information about a specific user
func (gc *GitHubClient) GetUser(ctx context.Context, username string) (*github.User, error) {
	if username == "" {
		return nil, errors.ErrEmptyUsername
	}

	user, resp, err := gc.client.Users.Get(ctx, username)
	if err != nil {
		return nil, gc.handleAPIError(err, resp, "get user", username)
	}

	return user, nil
}

// GetAuthenticatedUser retrieves information about the authenticated user
func (gc *GitHubClient) GetAuthenticatedUser(ctx context.Context) (*github.User, error) {
	user, resp, err := gc.client.Users.Get(ctx, "")
	if err != nil {
		return nil, gc.handleAPIError(err, resp, "get authenticated user", "")
	}

	return user, nil
}

// handleAPIError handles GitHub API errors and returns appropriate error types
func (gc *GitHubClient) handleAPIError(err error, resp *github.Response, op string, username string) error {
	if resp == nil {
		return errors.NewFollowError(op, username, errors.ErrNetworkError)
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return errors.NewFollowError(op, username, errors.ErrUnauthorized)
	case http.StatusForbidden:
		if strings.Contains(err.Error(), "rate limit") {
			return errors.NewFollowError(op, username, errors.ErrAPILimitExceeded)
		}
		return errors.NewFollowError(op, username, fmt.Errorf("forbidden: %w", err))
	case http.StatusNotFound:
		return errors.NewFollowError(op, username, fmt.Errorf("user not found on GitHub: %s", username))
	default:
		return errors.NewFollowError(op, username, err)
	}
}

// BatchFollow follows multiple users
func (gc *GitHubClient) BatchFollow(ctx context.Context, usernames []string) ([]string, []error) {
	var succeeded []string
	var errs []error

	for _, username := range usernames {
		if err := gc.Follow(ctx, username); err != nil {
			errs = append(errs, errors.NewFollowError("batch follow", username, err))
		} else {
			succeeded = append(succeeded, username)
		}
		// Rate limiting between requests
		time.Sleep(200 * time.Millisecond)
	}

	return succeeded, errs
}

// BatchUnfollow unfollows multiple users
func (gc *GitHubClient) BatchUnfollow(ctx context.Context, usernames []string) ([]string, []error) {
	var succeeded []string
	var errs []error

	for _, username := range usernames {
		if err := gc.Unfollow(ctx, username); err != nil {
			errs = append(errs, errors.NewFollowError("batch unfollow", username, err))
		} else {
			succeeded = append(succeeded, username)
		}
		// Rate limiting between requests
		time.Sleep(200 * time.Millisecond)
	}

	return succeeded, errs
}

// SyncFollowing synchronizes the local follow list with GitHub
func (gc *GitHubClient) SyncFollowing(ctx context.Context, localList *models.FollowList) (*models.FollowList, error) {
	// Get all users we're following on GitHub
	githubFollowing, err := gc.GetFollowing(ctx)
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup
	githubUsers := make(map[string]bool)
	for _, user := range githubFollowing {
		githubUsers[user.GetLogin()] = true
	}

	// Update local list
	newList := models.NewFollowList()
	for _, f := range localList.Follows {
		if githubUsers[f.Username] {
			newList.Follows = append(newList.Follows, f)
		}
	}

	// Add new follows from GitHub
	localUsers := make(map[string]bool)
	for _, f := range localList.Follows {
		localUsers[f.Username] = true
	}

	for _, user := range githubFollowing {
		username := user.GetLogin()
		if !localUsers[username] {
			newList.Add(username, "", nil)
		}
	}

	newList.Metadata.LastSync = time.Now()
	newList.Metadata.SyncStatus = "success"

	return newList, nil
}
