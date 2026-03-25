package cmd

import (
	"context"
	"fmt"
	"time"

	go_github "github.com/google/go-github/v55/github"
	"github.com/urfave/cli/v2"

	"github.com/h1s97x/gh-follow/internal/cache"
	gh_client "github.com/h1s97x/gh-follow/internal/github"
	"github.com/h1s97x/gh-follow/internal/models"
	"github.com/h1s97x/gh-follow/internal/storage"
)

// CacheCmd handles the cache command
func CacheCmd(c *cli.Context) error {
	return CacheStatus(c)
}

// CacheStatus shows the cache status
func CacheStatus(c *cli.Context) error {
	uc := cache.NewCache(cache.DefaultCachePath(), 24*time.Hour)
	stats := uc.GetStats()

	fmt.Println("\n📦 User Cache Status")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Cached Users: %d / %d\n", stats.TotalUsers, stats.MaxSize)
	fmt.Printf("Cache TTL:    %s\n", stats.TTL)

	if stats.TotalUsers > 0 {
		fmt.Printf("Oldest Entry: %s\n", stats.Oldest.Format("2006-01-02 15:04:05"))
		fmt.Printf("Newest Entry: %s\n", stats.Newest.Format("2006-01-02 15:04:05"))
	}

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("\nCache file: %s\n", cache.DefaultCachePath())

	return nil
}

// CacheList lists all cached users
func CacheList(c *cli.Context) error {
	format := c.String("format")
	uc := cache.NewCache(cache.DefaultCachePath(), 24*time.Hour)
	users := uc.List()

	if len(users) == 0 {
		fmt.Println("Cache is empty")
		return nil
	}

	switch format {
	case "json":
		for _, user := range users {
			fmt.Printf("%s\n", user.Login)
		}
	default:
		fmt.Println("\n📦 Cached Users")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("%-25s %-20s %-15s\n", "Username", "Name", "Cached At")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		for _, user := range users {
			name := user.Name
			if name == "" {
				name = "-"
			}
			if len(name) > 18 {
				name = name[:15] + "..."
			}
			fmt.Printf("%-25s %-20s %-15s\n",
				user.Login,
				name,
				user.CachedAt.Format("2006-01-02"),
			)
		}

		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("Total: %d users\n", len(users))
	}

	return nil
}

// CacheClear clears the cache
func CacheClear(c *cli.Context) error {
	force := c.Bool("force")

	if !force {
		fmt.Print("Are you sure you want to clear the cache? [y/N]: ")
		var confirm string
		if _, err := fmt.Scanln(&confirm); err != nil {
			fmt.Println("Cancelled")
			return nil
		}
		if confirm != "y" && confirm != "Y" {
			fmt.Println("Cancelled")
			return nil
		}
	}

	uc := cache.NewCache(cache.DefaultCachePath(), 24*time.Hour)
	if err := uc.Clear(); err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	fmt.Println("✅ Cache cleared")
	return nil
}

// CacheCleanup removes expired entries from cache
func CacheCleanup(c *cli.Context) error {
	uc := cache.NewCache(cache.DefaultCachePath(), 24*time.Hour)
	removed := uc.Cleanup()

	fmt.Printf("✅ Removed %d expired entries\n", removed)
	return nil
}

// CacheRefresh refreshes the cache for all followed users
func CacheRefresh(c *cli.Context) error {
	token, err := gh_client.GetTokenFromGH()
	if err != nil {
		return fmt.Errorf("failed to get GitHub token: %w", err)
	}

	gc := gh_client.NewGitHubClient(token, "github.com")
	st := storage.NewStorage(storage.DefaultStoragePath())
	uc := cache.NewCache(cache.DefaultCachePath(), 24*time.Hour)

	// Load follow list
	list, err := st.Load()
	if err != nil {
		return fmt.Errorf("failed to load follow list: %w", err)
	}

	if len(list.Follows) == 0 {
		fmt.Println("No users to refresh")
		return nil
	}

	fmt.Printf("Refreshing cache for %d users...\n\n", len(list.Follows))

	ctx := context.Background()
	successCount := 0
	failCount := 0

	for _, f := range list.Follows {
		user, err := gc.GetUser(ctx, f.Username)
		if err != nil {
			fmt.Printf("⚠️  Failed to fetch %s: %v\n", f.Username, err)
			failCount++
			continue
		}

		cachedUser := &cache.UserCache{
			Login:       user.GetLogin(),
			ID:          user.GetID(),
			AvatarURL:   user.GetAvatarURL(),
			Name:        user.GetName(),
			Bio:         user.GetBio(),
			Blog:        user.GetBlog(),
			Location:    user.GetLocation(),
			Company:     user.GetCompany(),
			PublicRepos: user.GetPublicRepos(),
			Followers:   user.GetFollowers(),
			Following:   user.GetFollowing(),
		}

		if err := uc.Set(cachedUser); err != nil {
			fmt.Printf("⚠️  Failed to cache %s: %v\n", f.Username, err)
		}
		fmt.Printf("✅ Refreshed %s\n", f.Username)
		successCount++

		// Rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("\n✅ Refreshed %d users, %d failed\n", successCount, failCount)

	return nil
}

// CacheShow shows detailed info for a cached user
func CacheShow(c *cli.Context) error {
	if c.Args().Len() == 0 {
		return fmt.Errorf("please provide a username")
	}

	username := c.Args().Get(0)
	uc := cache.NewCache(cache.DefaultCachePath(), 24*time.Hour)

	user := uc.Get(username)
	if user == nil {
		fmt.Printf("User %s not found in cache\n", username)
		fmt.Println("Use 'gh follow cache refresh' to populate cache")
		return nil
	}

	fmt.Printf("\n👤 %s\n", user.Login)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if user.Name != "" {
		fmt.Printf("Name:         %s\n", user.Name)
	}
	if user.Bio != "" {
		fmt.Printf("Bio:          %s\n", user.Bio)
	}
	if user.Company != "" {
		fmt.Printf("Company:      %s\n", user.Company)
	}
	if user.Location != "" {
		fmt.Printf("Location:     %s\n", user.Location)
	}
	if user.Blog != "" {
		fmt.Printf("Blog:         %s\n", user.Blog)
	}

	fmt.Printf("\nRepositories: %d\n", user.PublicRepos)
	fmt.Printf("Followers:    %d\n", user.Followers)
	fmt.Printf("Following:    %d\n", user.Following)
	fmt.Printf("\nCached at:    %s\n", user.CachedAt.Format(time.RFC3339))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	return nil
}

// FetchUserWithCache fetches user info with cache support
func FetchUserWithCache(ctx context.Context, gc *gh_client.GitHubClient, username string) (*cache.UserCache, error) {
	uc := cache.NewCache(cache.DefaultCachePath(), 24*time.Hour)

	// Check cache first
	cached := uc.Get(username)
	if cached != nil {
		return cached, nil
	}

	// Fetch from GitHub
	user, err := gc.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}

	// Cache the result
	cachedUser := &cache.UserCache{
		Login:       user.GetLogin(),
		ID:          user.GetID(),
		AvatarURL:   user.GetAvatarURL(),
		Name:        user.GetName(),
		Bio:         user.GetBio(),
		Blog:        user.GetBlog(),
		Location:    user.GetLocation(),
		Company:     user.GetCompany(),
		PublicRepos: user.GetPublicRepos(),
		Followers:   user.GetFollowers(),
		Following:   user.GetFollowing(),
	}

	// Ignore cache write errors - the data will still be returned
	_ = uc.Set(cachedUser)

	return cachedUser, nil
}

// GetCachedUsersMap returns a map of cached users for quick lookup
func GetCachedUsersMap() map[string]*cache.UserCache {
	uc := cache.NewCache(cache.DefaultCachePath(), 24*time.Hour)
	users := uc.List()

	result := make(map[string]*cache.UserCache)
	for _, user := range users {
		result[user.Login] = user
	}

	return result
}

// ConvertGitHubUser converts a GitHub user to a cached user
func ConvertGitHubUser(user *go_github.User) *cache.UserCache {
	return &cache.UserCache{
		Login:       user.GetLogin(),
		ID:          user.GetID(),
		AvatarURL:   user.GetAvatarURL(),
		Name:        user.GetName(),
		Bio:         user.GetBio(),
		Blog:        user.GetBlog(),
		Location:    user.GetLocation(),
		Company:     user.GetCompany(),
		PublicRepos: user.GetPublicRepos(),
		Followers:   user.GetFollowers(),
		Following:   user.GetFollowing(),
	}
}

// ConvertModelToCachedUser converts a models.Follow to cached user (minimal)
func ConvertModelToCachedUser(f models.Follow) *cache.UserCache {
	return &cache.UserCache{
		Login: f.Username,
	}
}
