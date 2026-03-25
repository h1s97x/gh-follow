# GH-Follow

A GitHub CLI extension to manage your follow list from the terminal.

[![CI](https://github.com/h1s97x/gh-follow/actions/workflows/ci.yml/badge.svg)](https://github.com/h1s97x/gh-follow/actions/workflows/ci.yml)
[![Release](https://github.com/h1s97x/gh-follow/actions/workflows/release.yml/badge.svg)](https://github.com/h1s97x/gh-follow/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/h1s97x/gh-follow)](https://goreportcard.com/report/github.com/h1s97x/gh-follow)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Overview

**GH-Follow** fills the gap in the official `gh` CLI by providing follow/unfollow functionality. It allows you to:

- Follow and unfollow GitHub users from the command line
- Manage a local follow list with notes and tags
- Sync your local list with GitHub
- **Cloud sync via GitHub Gist** - cross-device synchronization
- **Conflict detection and resolution** - smart merging
- **Auto-sync** - automatic synchronization
- Import and export your follow list

## Installation

### As a GitHub CLI Extension (Recommended)

```bash
gh extension install h1s97x/gh-follow
```

### From Source

```bash
git clone https://github.com/h1s97x/gh-follow.git
cd gh-follow
make build
```

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/h1s97x/gh-follow/releases).

## Prerequisites

- [GitHub CLI](https://cli.github.com/) installed and authenticated
- Go 1.22 or later (for building from source)

## Usage

### Follow a user

```bash
# Follow a single user
gh follow add octocat

# Follow multiple users
gh follow add user1 user2 user3

# Follow with notes and tags
gh follow add octocat --notes "GitHub mascot" --tags developer,go
```

### Unfollow a user

```bash
# Unfollow a single user
gh follow remove octocat

# Unfollow multiple users
gh follow remove user1 user2 user3

# Force without confirmation
gh follow remove octocat --force
```

### List your follows

```bash
# List all follows (table format)
gh follow list

# Output as JSON
gh follow list --format json

# Output as simple list (just usernames)
gh follow list --format simple

# Sort and filter
gh follow list --sort name --order asc
gh follow list --limit 10 --filter go

# Filter by tag
gh follow list --tag developer

# Filter by date range
gh follow list --date-from 2024-01-01 --date-to 2024-03-31
```

### Sync with GitHub

```bash
# Sync both ways (default)
gh follow sync

# Pull from GitHub
gh follow sync --direction pull

# Push to GitHub
gh follow sync --direction push

# Dry run
gh follow sync --dry-run
```

### Gist Sync (Cloud Sync)

```bash
# Create a Gist for sync
gh follow gist create

# Check Gist sync status
gh follow gist status

# Pull from Gist
gh follow gist pull

# Push to Gist
gh follow gist push

# Sync with Gist
gh follow sync --gist
```

### Auto-Sync

```bash
# Check auto-sync status
gh follow autosync status

# Manually trigger sync
gh follow autosync trigger

# Enable/disable auto-sync
gh follow config set sync.auto_sync true
gh follow config set sync.auto_sync false

# Set sync interval (seconds)
gh follow config set sync.sync_interval 3600
```

### Import/Export

```bash
# Export to JSON
gh follow export --output follows.json

# Export to CSV
gh follow export --format csv --output follows.csv

# Import from file
gh follow import --input follows.json

# Import and merge with existing
gh follow import --input follows.json --merge
```

### Statistics

```bash
# Show follow statistics
gh follow stats

# Output as JSON
gh follow stats --format json
```

### Cache Management

```bash
# Show cache status
gh follow cache status

# List cached users
gh follow cache list

# Clear cache
gh follow cache clear --force

# Cleanup expired entries
gh follow cache cleanup

# Refresh cache for all followed users
gh follow cache refresh

# Show cached user details
gh follow cache show octocat
```

### Follow Suggestions

```bash
# Get personalized follow suggestions
gh follow suggest

# Limit suggestions
gh follow suggest --limit 10

# Show trending users
gh follow suggest trending

# Trending users by language
gh follow suggest trending --language go

# Check mutual followers
gh follow suggest mutual

# Find inactive followed users
gh follow suggest inactive --days 365
```

### Batch Operations

```bash
# Follow multiple users
gh follow batch follow user1 user2 user3

# Follow from file
gh follow batch follow --file users.txt

# Dry run (preview only)
gh follow batch follow user1 user2 --dry-run

# Unfollow multiple users
gh follow batch unfollow user1 user2 user3

# Check if users follow you
gh follow batch check user1 user2 user3

# Import usernames from file
gh follow batch import users.txt

# Import and follow all
gh follow batch import users.txt --follow
```

### Configuration

```bash
# Show current configuration
gh follow config

# Get a specific value
gh follow config get display.default_format

# Set a configuration value
gh follow config set display.default_format json

# Reset to defaults
gh follow config reset --force
```

## Configuration Options

| Key | Description | Default |
|-----|-------------|---------|
| `storage.local_path` | Path to local follow list | `~/.config/gh/follow-list.json` |
| `storage.use_gist` | Enable Gist sync | `false` |
| `storage.gist_id` | Gist ID for sync | `""` |
| `sync.auto_sync` | Enable auto sync | `true` |
| `sync.sync_interval` | Sync interval in seconds | `3600` |
| `display.default_format` | Default output format | `table` |
| `display.default_sort` | Default sort field | `date` |
| `display.default_order` | Default sort order | `desc` |

## Conflict Resolution

When syncing between local and remote (GitHub or Gist), conflicts may arise. GH-Follow provides three resolution strategies:

1. **newest-wins** (default): Uses the entry with the most recent timestamp
2. **local-wins**: Local changes take precedence
3. **remote-wins**: Remote changes take precedence

Use `--force` flag to override conflict detection.

## Data Storage

### Local Storage

Follow list is stored at: `~/.config/gh/follow-list.json`

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

### Gist Storage

When Gist sync is enabled, a private Gist is created to store your follow list. This allows:
- Cross-device synchronization
- Version history
- Backup and recovery

## Development

### Build

```bash
make build
```

### Test

```bash
make test
```

### Install locally

```bash
make install
```

## Project Structure

```
gh-follow/
├── main.go              # CLI entry point
├── version.go           # Version information
├── internal/            # Internal packages
│   ├── models.go        # Data models
│   ├── storage.go       # Storage operations
│   ├── github_api.go    # GitHub API client
│   ├── config.go        # Configuration management
│   ├── gist_sync.go     # Gist synchronization
│   ├── gist_cmd.go      # Gist commands
│   ├── auto_sync.go     # Auto-sync functionality
│   ├── cache.go         # User info caching
│   ├── suggest.go       # Follow suggestions
│   ├── batch.go         # Batch operations
│   ├── add.go           # Add command
│   ├── remove.go        # Remove command
│   ├── list.go          # List command
│   ├── sync.go          # Sync command
│   └── *_test.go        # Tests
├── .github/             # GitHub workflows
└── Makefile             # Build automation
```

## Roadmap

### Phase 1: MVP ✅
- [x] Basic add/remove/list commands
- [x] Local storage
- [x] GitHub API integration

### Phase 2: Enhanced Features ✅
- [x] Advanced filtering and sorting
- [x] Configuration management
- [x] Integration tests
- [x] CI/CD workflows

### Phase 3: Cloud Sync ✅
- [x] Gist sync implementation
- [x] Conflict detection and resolution
- [x] Auto-sync functionality

### Phase 4: Advanced Features ✅
- [x] User information caching
- [x] Follow suggestions
- [x] Performance optimization (batch operations)
- [x] Batch operations with concurrency

### Phase 5: Future Enhancements
- [ ] Multi-account support

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Security

See [SECURITY.md](SECURITY.md) for security policy.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Related Projects

- [gh-token](https://github.com/Link-/gh-token) - GitHub App token management
- [GitHub CLI](https://github.com/cli/cli) - The official GitHub CLI
