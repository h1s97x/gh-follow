# Architecture Overview

This document describes the architecture of gh-follow, its modules, data flow, and extension points.

## Table of Contents

- [High-Level Architecture](#high-level-architecture)
- [Core Modules](#core-modules)
- [Data Flow](#data-flow)
- [Storage Architecture](#storage-architecture)
- [Sync Mechanism](#sync-mechanism)
- [Extension Guide](#extension-guide)

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        CLI Layer (main.go)                       │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐   │
│  │   add   │ │  remove │ │   list  │ │  sync   │ │  batch  │   │
│  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘   │
└───────┼──────────┼──────────┼──────────┼──────────┼─────────────┘
        │          │          │          │          │
┌───────┴──────────┴──────────┴──────────┴──────────┴─────────────┐
│                     Command Handlers Layer                       │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐   │
│  │ add.go  │ │remove.go│ │ list.go │ │sync.go  │ │batch.go │   │
│  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘   │
└───────┼──────────┼──────────┼──────────┼──────────┼─────────────┘
        │          │          │          │          │
┌───────┴──────────┴──────────┴──────────┴──────────┴─────────────┐
│                     Core Services Layer                          │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────────────┐   │
│  │   Storage   │  │ GitHubClient │  │   SuggestionEngine   │   │
│  │  (local)    │  │   (API)      │  │                      │   │
│  └──────┬──────┘  └──────┬───────┘  └──────────────────────┘   │
│         │                │                                       │
│  ┌──────┴──────┐  ┌──────┴───────┐  ┌──────────────────────┐   │
│  │   Cache     │  │  GistSync    │  │   BatchProcessor     │   │
│  │  (memory)   │  │  (cloud)     │  │                      │   │
│  └─────────────┘  └──────────────┘  └──────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────┴───────────────────────────────────┐
│                     External Services                            │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────────────┐   │
│  │ GitHub API  │  │  GitHub Gist │  │   GitHub CLI Auth    │   │
│  └─────────────┘  └──────────────┘  └──────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

## Core Modules

### 1. Models (`internal/models.go`)

Core data structures for the application.

```go
// Primary data model
type FollowList struct {
    Version   string    `json:"version"`
    UpdatedAt time.Time `json:"updated_at"`
    Follows   []Follow  `json:"follows"`
    Metadata  Metadata  `json:"metadata"`
}

// Individual follow record
type Follow struct {
    Username  string    `json:"username"`
    FollowedAt time.Time `json:"followed_at"`
    Notes     string    `json:"notes,omitempty"`
    Tags      []string  `json:"tags,omitempty"`
}

// Configuration model
type Config struct {
    Version string
    Storage StorageConf
    Sync    SyncConf
    Display DisplayConf
}
```

**Responsibilities:**
- Define data structures
- Provide CRUD operations on FollowList
- Calculate statistics

### 2. Storage (`internal/storage.go`)

Local file storage management.

```go
type Storage struct {
    path string
}

func (s *Storage) Load() (*FollowList, error)
func (s *Storage) Save(list *FollowList) error
func (s *Storage) Add(username string, notes string, tags []string) error
func (s *Storage) Remove(username string) error
func (s *Storage) Export(outputPath string, format string) error
func (s *Storage) Import(inputPath string, merge bool) error
```

**Responsibilities:**
- Read/write follow list to disk
- Handle file permissions
- Import/export functionality

### 3. GitHub Client (`internal/github_api.go`)

GitHub API interaction layer.

```go
type GitHubClient struct {
    client *github.Client
}

func (gc *GitHubClient) Follow(ctx context.Context, username string) error
func (gc *GitHubClient) Unfollow(ctx context.Context, username string) error
func (gc *GitHubClient) GetUser(ctx context.Context, username string) (*github.User, error)
func (gc *GitHubClient) IsFollowing(ctx context.Context, username string) (bool, error)
```

**Responsibilities:**
- GitHub API authentication
- Follow/unfollow operations
- User information retrieval
- Rate limit handling

### 4. Gist Sync (`internal/gist_sync.go`)

Cloud synchronization via GitHub Gist.

```go
type GistSync struct {
    client  *github.Client
    gistID  string
    storage *Storage
}

func (gs *GistSync) Create(ctx context.Context) (string, error)
func (gs *GistSync) Upload(ctx context.Context) error
func (gs *GistSync) Download(ctx context.Context) (*FollowList, error)
func (gs *GistSync) DetectConflicts(local, remote *FollowList) []Conflict
func (gs *GistSync) MergeLists(local, remote *FollowList, strategy string) *FollowList
```

**Responsibilities:**
- Create sync Gist
- Upload/download follow list
- Conflict detection and resolution
- Merge strategies (newest-wins, local-wins, remote-wins)

### 5. Cache (`internal/cache.go`)

User information caching system.

```go
type Cache struct {
    path    string
    users   map[string]*UserCache
    ttl     time.Duration
    maxSize int
}

func (c *Cache) Get(username string) *UserCache
func (c *Cache) Set(user *UserCache)
func (c *Cache) Cleanup() int
func (c *Cache) GetStats() *CacheStats
```

**Responsibilities:**
- Cache user information
- TTL-based expiration
- LRU eviction
- Reduce API calls

### 6. Suggestion Engine (`internal/suggest.go`)

Follow suggestion algorithms.

```go
type SuggestionEngine struct {
    gc      *GitHubClient
    storage *Storage
    cache   *Cache
}

func (se *SuggestionEngine) GenerateSuggestions(ctx context.Context, opts *SuggestionOptions) ([]*Suggestion, error)
func (se *SuggestionEngine) GetTrendingUsers(ctx context.Context, language string, limit int) ([]*Suggestion, error)
```

**Suggestion Strategies:**
- **Follow-of-follow**: Users followed by people you follow
- **Org members**: Members of your organizations
- **Star contributors**: Contributors to repos you starred
- **Similar users**: Users in the same location

### 7. Batch Processor (`internal/batch.go`)

Concurrent batch operations.

```go
type BatchProcessor struct {
    gc          *GitHubClient
    concurrency int
    rateLimit   time.Duration
}

func (bp *BatchProcessor) BatchFollow(ctx context.Context, usernames []string, opts *BatchOptions) []*BatchOperation
func (bp *BatchProcessor) BatchUnfollow(ctx context.Context, usernames []string, opts *BatchOptions) []*BatchOperation
func (bp *BatchProcessor) BatchCheckFollowers(ctx context.Context, usernames []string) map[string]bool
```

**Features:**
- Concurrent execution
- Rate limiting
- Progress tracking
- Error aggregation

### 8. Configuration (`internal/config.go`)

Configuration management.

```go
type ConfigManager struct {
    path string
}

func (cm *ConfigManager) Load() (*Config, error)
func (cm *ConfigManager) Save(config *Config) error
func (cm *ConfigManager) Get(key string) (interface{}, error)
func (cm *ConfigManager) Set(key string, value interface{}) error
```

## Data Flow

### Follow a User

```
User Input (CLI)
      │
      ▼
┌─────────────┐
│  add.go     │ Parse flags & validate
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ GitHubClient│ Check if user exists
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ GitHubClient│ Call GitHub API to follow
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Storage   │ Add to local list
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   GistSync  │ Sync to cloud (if enabled)
└─────────────┘
```

### Sync with Gist

```
User Input (gh follow sync --gist)
              │
              ▼
┌─────────────────────┐
│     sync.go         │ Determine direction
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│     Storage         │ Load local list
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│     GistSync        │ Download remote list
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│     GistSync        │ Detect conflicts
└──────────┬──────────┘
           │
           ▼
      ┌────┴────┐
      │Conflicts?│
      └────┬────┘
           │
    ┌──────┴──────┐
    │             │
   Yes            No
    │             │
    ▼             │
┌─────────┐       │
│  Merge  │       │
│ (strategy)│      │
└────┬────┘       │
     │            │
     └─────┬──────┘
           │
           ▼
┌─────────────────────┐
│     Storage         │ Save merged list
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│     GistSync        │ Upload to Gist
└─────────────────────┘
```

### Generate Suggestions

```
User Input (gh follow suggest)
              │
              ▼
┌─────────────────────┐
│ SuggestionEngine    │ Initialize engine
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│     Storage         │ Load follow list
└──────────┬──────────┘
           │
           ├──────────────────┬──────────────────┐
           │                  │                  │
           ▼                  ▼                  ▼
    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
    │Follow-of-Follow│   │ Org Members │    │Starred Repos│
    └──────┬──────┘    └──────┬──────┘    └──────┬──────┘
           │                  │                  │
           ▼                  ▼                  ▼
    ┌─────────────────────────────────────────────────┐
    │              Aggregate & Score                   │
    └────────────────────────┬────────────────────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │  Sort & Filter  │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │  Return Results │
                    └─────────────────┘
```

## Storage Architecture

### Local Storage

```
~/.config/gh/
├── follow-list.json       # Main follow list
├── follow-config.json     # Configuration
└── follow-cache.json      # User info cache
```

### File Formats

**follow-list.json:**
```json
{
  "version": "1.0.0",
  "updated_at": "2024-03-25T10:30:00Z",
  "follows": [
    {
      "username": "octocat",
      "followed_at": "2024-03-20T15:30:00Z",
      "notes": "GitHub mascot",
      "tags": ["developer", "go"]
    }
  ],
  "metadata": {
    "total_count": 1,
    "last_sync": "2024-03-25T10:30:00Z",
    "sync_status": "success"
  }
}
```

**follow-config.json:**
```json
{
  "version": "1.0.0",
  "storage": {
    "local_path": "~/.config/gh/follow-list.json",
    "use_gist": true,
    "gist_id": "abc123..."
  },
  "sync": {
    "auto_sync": true,
    "sync_interval": 3600
  },
  "display": {
    "default_format": "table",
    "default_sort": "date",
    "default_order": "desc"
  }
}
```

### Gist Storage

When Gist sync is enabled:

```
Gist ID: abc123...
├── follow-list.json    # Same format as local
└── metadata.json       # Sync metadata
```

## Sync Mechanism

### Conflict Detection

```go
type Conflict struct {
    Username    string
    LocalEntry  *Follow
    RemoteEntry *Follow
    Type        ConflictType
}

type ConflictType int

const (
    ConflictModified ConflictType = iota  // Both modified
    ConflictAdded                         // Added in both
    ConflictRemoved                       // Removed in one
)
```

### Resolution Strategies

| Strategy | Behavior |
|----------|----------|
| `newest-wins` | Use entry with latest timestamp |
| `local-wins` | Always prefer local changes |
| `remote-wins` | Always prefer remote changes |

### Auto-Sync

```go
type AutoSync struct {
    interval    time.Duration
    lastSync    time.Time
    syncManager *SyncManager
}

func (as *AutoSync) ShouldSync() bool {
    return time.Since(as.lastSync) > as.interval
}

func (as *AutoSync) Trigger() error {
    // Perform sync
    as.lastSync = time.Now()
    return nil
}
```

## Extension Guide

### Adding a New Command

1. **Create command file** (`internal/new_cmd.go`):

```go
package internal

import "github.com/urfave/cli/v2"

func NewCommand(c *cli.Context) error {
    // Implementation
    return nil
}
```

2. **Add flags** (`internal/new_cmd_flags.go`):

```go
package internal

import "github.com/urfave/cli/v2"

func NewCommandFlags() []cli.Flag {
    return []cli.Flag{
        &cli.StringFlag{
            Name:  "example",
            Usage: "Example flag",
        },
    }
}
```

3. **Register in `main.go`**:

```go
{
    Name:   "new",
    Usage:  "New command",
    Flags:  internal.NewCommandFlags(),
    Action: internal.NewCommand,
},
```

### Adding a New Storage Backend

1. **Define interface**:

```go
type StorageBackend interface {
    Load() (*FollowList, error)
    Save(list *FollowList) error
}
```

2. **Implement backend**:

```go
type S3Storage struct {
    bucket string
    key    string
}

func (s *S3Storage) Load() (*FollowList, error) {
    // Implementation
}

func (s *S3Storage) Save(list *FollowList) error {
    // Implementation
}
```

3. **Add configuration**:

```go
type StorageConf struct {
    LocalPath string `json:"local_path"`
    UseGist   bool   `json:"use_gist"`
    GistID    string `json:"gist_id"`
    // Add new backend config
    UseS3     bool   `json:"use_s3"`
    S3Bucket  string `json:"s3_bucket"`
}
```

### Adding a New Suggestion Strategy

1. **Add to SuggestionEngine**:

```go
func (se *SuggestionEngine) suggestFromNewStrategy(ctx context.Context) ([]*Suggestion, error) {
    // Implementation
    return suggestions, nil
}
```

2. **Integrate into GenerateSuggestions**:

```go
func (se *SuggestionEngine) GenerateSuggestions(ctx context.Context, opts *SuggestionOptions) ([]*Suggestion, error) {
    // ...existing strategies...
    
    // New strategy
    new, err := se.suggestFromNewStrategy(ctx)
    if err == nil {
        for _, s := range new {
            suggestions[s.Username] = s
        }
    }
    
    // ...
}
```

## Performance Considerations

### API Rate Limits

| API | Limit | With Auth |
|-----|-------|-----------|
| REST API | 60/hr | 5000/hr |
| GraphQL | - | 5000/hr |

### Optimization Strategies

1. **Caching**: Reduce redundant API calls
2. **Batching**: Concurrent operations with rate limiting
3. **Pagination**: Handle large lists efficiently
4. **Lazy Loading**: Load data on demand

### Memory Management

```go
// Use pagination for large lists
func (gc *GitHubClient) GetAllFollowers(ctx context.Context, username string) ([]string, error) {
    var allFollowers []string
    opts := &github.ListOptions{PerPage: 100}
    
    for {
        followers, resp, err := gc.client.Users.ListFollowers(ctx, username, opts)
        if err != nil {
            return nil, err
        }
        
        for _, f := range followers {
            allFollowers = append(allFollowers, f.GetLogin())
        }
        
        if resp.NextPage == 0 {
            break
        }
        opts.Page = resp.NextPage
    }
    
    return allFollowers, nil
}
```

## Security Considerations

- **Token Storage**: Use GitHub CLI's token storage
- **File Permissions**: 600 for sensitive files
- **Gist Privacy**: Create private gists by default
- **Input Validation**: Sanitize all user inputs
- **Rate Limiting**: Prevent API abuse
