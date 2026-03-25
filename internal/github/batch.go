package github

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/go-github/v55/github"
)

// BatchOperation represents a batch operation result
type BatchOperation struct {
	Username string
	Success  bool
	Error    error
}

// BatchProcessor handles concurrent batch operations
type BatchProcessor struct {
	gc          *GitHubClient
	concurrency int
	rateLimit   time.Duration
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(gc *GitHubClient, opts *BatchOptions) *BatchProcessor {
	if opts == nil {
		opts = &BatchOptions{
			Concurrency: 5,
			RateLimit:   100 * time.Millisecond,
		}
	}

	return &BatchProcessor{
		gc:          gc,
		concurrency: opts.Concurrency,
		rateLimit:   opts.RateLimit,
	}
}

// BatchOptions represents options for batch processing
type BatchOptions struct {
	Concurrency int
	RateLimit   time.Duration
	DryRun      bool
	Progress    bool
}

// BatchFollow follows multiple users concurrently
func (bp *BatchProcessor) BatchFollow(ctx context.Context, usernames []string, opts *BatchOptions) []*BatchOperation {
	results := make([]*BatchOperation, len(usernames))
	var wg sync.WaitGroup
	sem := make(chan struct{}, bp.concurrency)

	successCount := int32(0)
	failCount := int32(0)

	for i, username := range usernames {
		wg.Add(1)
		go func(idx int, user string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			// Rate limiting
			if bp.rateLimit > 0 {
				time.Sleep(bp.rateLimit)
			}

			results[idx] = &BatchOperation{
				Username: user,
			}

			if opts != nil && opts.DryRun {
				results[idx].Success = true
				atomic.AddInt32(&successCount, 1)
				return
			}

			_, err := bp.gc.client.Users.Follow(ctx, user)
			if err != nil {
				results[idx].Error = err
				results[idx].Success = false
				atomic.AddInt32(&failCount, 1)
			} else {
				results[idx].Success = true
				atomic.AddInt32(&successCount, 1)
			}

			// Progress update
			if opts != nil && opts.Progress {
				fmt.Printf("\r[Progress] Success: %d | Failed: %d | Total: %d/%d",
					successCount, failCount, idx+1, len(usernames))
			}
		}(i, username)
	}

	wg.Wait()

	if opts != nil && opts.Progress {
		fmt.Println()
	}

	return results
}

// BatchUnfollow unfollows multiple users concurrently
func (bp *BatchProcessor) BatchUnfollow(ctx context.Context, usernames []string, opts *BatchOptions) []*BatchOperation {
	results := make([]*BatchOperation, len(usernames))
	var wg sync.WaitGroup
	sem := make(chan struct{}, bp.concurrency)

	successCount := int32(0)
	failCount := int32(0)

	for i, username := range usernames {
		wg.Add(1)
		go func(idx int, user string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			// Rate limiting
			if bp.rateLimit > 0 {
				time.Sleep(bp.rateLimit)
			}

			results[idx] = &BatchOperation{
				Username: user,
			}

			if opts != nil && opts.DryRun {
				results[idx].Success = true
				atomic.AddInt32(&successCount, 1)
				return
			}

			_, err := bp.gc.client.Users.Unfollow(ctx, user)
			if err != nil {
				results[idx].Error = err
				results[idx].Success = false
				atomic.AddInt32(&failCount, 1)
			} else {
				results[idx].Success = true
				atomic.AddInt32(&successCount, 1)
			}

			// Progress update
			if opts != nil && opts.Progress {
				fmt.Printf("\r[Progress] Success: %d | Failed: %d | Total: %d/%d",
					successCount, failCount, idx+1, len(usernames))
			}
		}(i, username)
	}

	wg.Wait()

	if opts != nil && opts.Progress {
		fmt.Println()
	}

	return results
}

// BatchCheckFollowers checks if users follow you concurrently
func (bp *BatchProcessor) BatchCheckFollowers(ctx context.Context, usernames []string) map[string]bool {
	results := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, bp.concurrency)

	for _, username := range usernames {
		wg.Add(1)
		go func(user string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			if bp.rateLimit > 0 {
				time.Sleep(bp.rateLimit)
			}

			isFollowing, _, err := bp.gc.client.Users.IsFollowing(ctx, user, "")
			if err != nil {
				mu.Lock()
				results[user] = false
				mu.Unlock()
				return
			}

			mu.Lock()
			results[user] = isFollowing
			mu.Unlock()
		}(username)
	}

	wg.Wait()
	return results
}

// BatchFetchUsers fetches multiple users concurrently
func (bp *BatchProcessor) BatchFetchUsers(ctx context.Context, usernames []string) map[string]*github.User {
	results := make(map[string]*github.User)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, bp.concurrency)

	for _, username := range usernames {
		wg.Add(1)
		go func(user string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			if bp.rateLimit > 0 {
				time.Sleep(bp.rateLimit)
			}

			u, _, err := bp.gc.client.Users.Get(ctx, user)
			if err != nil {
				return
			}

			mu.Lock()
			results[user] = u
			mu.Unlock()
		}(username)
	}

	wg.Wait()
	return results
}

// BatchResult represents batch operation results
type BatchResult struct {
	Total     int
	Success   int
	Failed    int
	Errors    []error
	Duration  time.Duration
	ErrorsByUser map[string]error
}

// GetBatchStats returns statistics from batch results
func GetBatchStats(results []*BatchOperation) *BatchResult {
	stats := &BatchResult{
		Total:        len(results),
		Errors:       make([]error, 0),
		ErrorsByUser: make(map[string]error),
	}

	for _, r := range results {
		if r.Success {
			stats.Success++
		} else {
			stats.Failed++
			if r.Error != nil {
				stats.Errors = append(stats.Errors, r.Error)
				stats.ErrorsByUser[r.Username] = r.Error
			}
		}
	}

	return stats
}

// DisplayBatchResults displays batch operation results
func DisplayBatchResults(results []*BatchOperation, operation string) {
	stats := GetBatchStats(results)

	fmt.Printf("\n📊 Batch %s Results\n", operation)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Total:   %d\n", stats.Total)
	fmt.Printf("Success: %d\n", stats.Success)
	fmt.Printf("Failed:  %d\n", stats.Failed)

	if stats.Failed > 0 {
		fmt.Println("\nFailed users:")
		for username, err := range stats.ErrorsByUser {
			fmt.Printf("  %s: %v\n", username, err)
		}
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
