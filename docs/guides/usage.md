# Usage Guide

Complete guide to using gh-follow for managing your GitHub follow list.

## Table of Contents

- [Quick Start](#quick-start)
- [Command Reference](#command-reference)
- [Common Workflows](#common-workflows)
- [Advanced Usage](#advanced-usage)
- [Examples](#examples)

## Quick Start

```bash
# Follow a user
gh follow add octocat

# List your follows
gh follow list

# Unfollow a user
gh follow remove octocat

# Get suggestions
gh follow suggest
```

## Command Reference

### add / follow

Follow one or more GitHub users.

```bash
gh follow add <username> [username...]
gh follow <username>  # alias
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--notes`, `-n` | string | "" | Add notes for this user |
| `--tags`, `-t` | string | "" | Comma-separated tags |
| `--no-sync` | bool | false | Skip sync after add |

**Examples:**

```bash
# Follow a single user
gh follow add octocat

# Follow multiple users
gh follow add user1 user2 user3

# Follow with notes
gh follow add octocat --notes "GitHub mascot"

# Follow with tags
gh follow add torvalds --tags "linux,creator"

# Follow with both
gh follow add github --notes "GitHub official" --tags "official,platform"
```

### remove / unfollow / rm / delete

Unfollow one or more GitHub users.

```bash
gh follow remove <username> [username...]
gh follow unfollow <username>  # alias
gh follow rm <username>        # alias
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--force`, `-f` | bool | false | Skip confirmation |
| `--no-sync` | bool | false | Skip sync after remove |

**Examples:**

```bash
# Unfollow a user (with confirmation)
gh follow remove octocat

# Force unfollow (no confirmation)
gh follow remove octocat --force

# Unfollow multiple users
gh follow remove user1 user2 user3
```

### list / ls

List all users you follow.

```bash
gh follow list
gh follow ls  # alias
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--format`, `-f` | string | "table" | Output format: table, json, simple |
| `--sort`, `-s` | string | "date" | Sort by: date, name |
| `--order`, `-o` | string | "desc" | Sort order: asc, desc |
| `--limit`, `-l` | int | 0 | Limit number of results (0 = all) |
| `--tag` | string | "" | Filter by tag |
| `--filter` | string | "" | Filter by username pattern |
| `--date-from` | string | "" | Filter from date (YYYY-MM-DD) |
| `--date-to` | string | "" | Filter to date (YYYY-MM-DD) |

**Examples:**

```bash
# List all follows (table format)
gh follow list

# List as JSON
gh follow list --format json

# List simple format (just usernames)
gh follow list --format simple

# Sort by name ascending
gh follow list --sort name --order asc

# Limit to 10 results
gh follow list --limit 10

# Filter by tag
gh follow list --tag developer

# Filter by username pattern
gh follow list --filter "go"

# Filter by date range
gh follow list --date-from 2024-01-01 --date-to 2024-03-31
```

### sync

Synchronize follow list with GitHub or Gist.

```bash
gh follow sync
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--direction`, `-d` | string | "both" | Sync direction: pull, push, both |
| `--gist` | bool | false | Sync with Gist instead of GitHub |
| `--strategy`, `-s` | string | "newest-wins" | Conflict resolution: newest-wins, local-wins, remote-wins |
| `--dry-run` | bool | false | Preview changes without applying |
| `--force` | bool | false | Skip conflict detection |

**Examples:**

```bash
# Sync both ways (default)
gh follow sync

# Pull from GitHub
gh follow sync --direction pull

# Push to GitHub
gh follow sync --direction push

# Sync with Gist
gh follow sync --gist

# Preview changes
gh follow sync --dry-run

# Use different conflict strategy
gh follow sync --strategy local-wins
```

### gist

Manage Gist synchronization.

```bash
gh follow gist <subcommand>
```

**Subcommands:**

| Command | Description |
|---------|-------------|
| `create` | Create a new Gist for sync |
| `status` | Show Gist sync status |
| `pull` | Pull follow list from Gist |
| `push` | Push follow list to Gist |

**Examples:**

```bash
# Create a Gist for sync
gh follow gist create

# Check Gist sync status
gh follow gist status

# Pull from Gist
gh follow gist pull

# Push to Gist
gh follow gist push

# Pull with conflict strategy
gh follow gist pull --strategy remote-wins
```

### autosync

Manage auto-sync settings.

```bash
gh follow autosync <subcommand>
```

**Subcommands:**

| Command | Description |
|---------|-------------|
| `status` | Show auto-sync status |
| `trigger` | Manually trigger sync |

**Examples:**

```bash
# Check auto-sync status
gh follow autosync status

# Manually trigger sync
gh follow autosync trigger

# Trigger sync with Gist
gh follow autosync trigger --gist
```

### export / exp

Export follow list to a file.

```bash
gh follow export
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--output`, `-o` | string | "follows.json" | Output file path |
| `--format`, `-f` | string | "json" | Export format: json, csv |

**Examples:**

```bash
# Export to JSON (default)
gh follow export --output follows.json

# Export to CSV
gh follow export --format csv --output follows.csv

# Export to specific path
gh follow export --output ~/backups/follows.json
```

### import / imp

Import follow list from a file.

```bash
gh follow import --input <file>
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--input`, `-i` | string | "" | Input file path (required) |
| `--merge`, `-m` | bool | false | Merge with existing list |
| `--follow` | bool | false | Follow imported users on GitHub |

**Examples:**

```bash
# Import (replace existing)
gh follow import --input follows.json

# Import and merge
gh follow import --input follows.json --merge

# Import and follow users
gh follow import --input follows.json --follow
```

### stats

Show follow list statistics.

```bash
gh follow stats
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--format`, `-f` | string | "table" | Output format: table, json |

**Examples:**

```bash
# Show statistics
gh follow stats

# Output as JSON
gh follow stats --format json
```

### config

Manage configuration.

```bash
gh follow config <subcommand>
```

**Subcommands:**

| Command | Description |
|---------|-------------|
| `show` | Show current configuration |
| `get` | Get a configuration value |
| `set` | Set a configuration value |
| `reset` | Reset to defaults |

**Examples:**

```bash
# Show all configuration
gh follow config show

# Get specific value
gh follow config get display.default_format

# Set a value
gh follow config set display.default_format json

# Reset to defaults
gh follow config reset --force
```

### cache

Manage user information cache.

```bash
gh follow cache <subcommand>
```

**Subcommands:**

| Command | Description |
|---------|-------------|
| `status` | Show cache status |
| `list` | List cached users |
| `clear` | Clear the cache |
| `cleanup` | Remove expired entries |
| `refresh` | Refresh cache for all followed users |
| `show` | Show cached user details |

**Examples:**

```bash
# Show cache status
gh follow cache status

# List cached users
gh follow cache list

# Clear cache
gh follow cache clear --force

# Cleanup expired entries
gh follow cache cleanup

# Refresh all cache
gh follow cache refresh

# Show user details
gh follow cache show octocat
```

### suggest / recommend

Get follow suggestions.

```bash
gh follow suggest
gh follow recommend  # alias
```

**Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit`, `-l` | int | 20 | Maximum suggestions |
| `--format`, `-f` | string | "table" | Output format |
| `--no-follow-of-follow` | bool | false | Exclude follow-of-follow |
| `--no-org-members` | bool | false | Exclude org members |
| `--no-star-contributors` | bool | false | Exclude star contributors |

**Subcommands:**

| Command | Description |
|---------|-------------|
| `trending` | Show trending users |
| `mutual` | Check mutual followers |
| `inactive` | Find inactive followed users |

**Examples:**

```bash
# Get suggestions
gh follow suggest

# Limit to 10 suggestions
gh follow suggest --limit 10

# Trending users
gh follow suggest trending

# Trending by language
gh follow suggest trending --language go

# Check mutual followers
gh follow suggest mutual

# Find inactive users
gh follow suggest inactive --days 365
```

### batch

Perform batch operations.

```bash
gh follow batch <subcommand>
```

**Subcommands:**

| Command | Description |
|---------|-------------|
| `follow` | Follow multiple users |
| `unfollow` | Unfollow multiple users |
| `check` | Check if users follow you |
| `import` | Import from file |

**Examples:**

```bash
# Batch follow
gh follow batch follow user1 user2 user3

# Follow from file
gh follow batch follow --file users.txt

# Dry run (preview only)
gh follow batch follow user1 user2 --dry-run

# Batch unfollow
gh follow batch unfollow user1 user2 user3

# Check followers
gh follow batch check user1 user2 user3

# Import from file and follow
gh follow batch import users.txt --follow
```

## Common Workflows

### Workflow 1: Initial Setup

```bash
# 1. Install
gh extension install h1s97x/gh-follow

# 2. Verify authentication
gh auth status

# 3. Check configuration
gh follow config show

# 4. Sync from GitHub
gh follow sync --direction pull
```

### Workflow 2: Daily Usage

```bash
# Follow users you discover
gh follow add interesting-user --notes "Met at conference" --tags "conference,go"

# Check your list
gh follow list --limit 10

# Get suggestions for new users
gh follow suggest --limit 5
```

### Workflow 3: Cross-Device Sync

```bash
# On device A: Setup Gist sync
gh follow gist create
gh follow gist push

# On device B: Pull from Gist
gh follow gist pull

# Both devices: Enable auto-sync
gh follow config set sync.auto_sync true
```

### Workflow 4: Bulk Management

```bash
# Export current list
gh follow export --output backup.json

# Import from backup
gh follow import --input backup.json --merge

# Batch follow from file
gh follow batch follow --file new-users.txt

# Find and remove inactive users
gh follow suggest inactive --days 365
gh follow batch unfollow inactive1 inactive2
```

### Workflow 5: Analysis and Cleanup

```bash
# Check statistics
gh follow stats

# Find users not following back
gh follow suggest mutual

# Clear old cache
gh follow cache cleanup

# Refresh user info
gh follow cache refresh
```

## Advanced Usage

### Using with Scripts

```bash
# Export usernames only for scripting
gh follow list --format simple > usernames.txt

# Count follows
gh follow list --format simple | wc -l

# Find users with specific tag
gh follow list --tag developer --format simple

# Bulk operations with xargs
gh follow list --format simple | xargs -I {} gh api users/{}/followers
```

### JSON Output for Tools

```bash
# Get JSON output for processing
gh follow list --format json | jq '.[].username'

# Export and process
gh follow export --format json --output - | jq '.follows | length'

# Get specific fields
gh follow stats --format json | jq '.total_follows'
```

### Integration with Other Tools

```bash
# Combine with gh api
gh follow list --format simple | while read user; do
  repos=$(gh api users/$user/repos --paginate -q '.[].name' | head -5)
  echo "$user: $repos"
done

# Combine with fzf for interactive selection
gh follow list --format simple | fzf --preview 'gh api users/{} | jq ".name, .bio"'

# Use with parallel for concurrent operations
gh follow list --format simple | parallel -j 5 gh api users/{}/followers -q 'length'
```

### Custom Output Formatting

```bash
# Custom format with jq
gh follow list --format json | jq -r '.[] | "\(.username)\t\(.tags | join(","))"'

# CSV export with custom columns
gh follow list --format json | jq -r '["username","date","tags"], (.[] | [.username, .followed_at, (.tags // []) | join(";")]) | @csv'
```

## Examples

### Example 1: Conference Networking

```bash
# After a conference, add all new contacts
gh follow batch follow \
  john-speaker jane-speaker bob-speaker \
  --notes "Met at GopherCon 2024" \
  --tags "gophercon,conference,speaker"

# Tag existing follows by topic
gh follow add existing-dev --tags "golang,kubernetes"
```

### Example 2: Open Source Project Team

```bash
# Follow project maintainers
gh follow batch follow maintainer1 maintainer2 maintainer3 \
  --tags "project-x,maintainer"

# Track contributors
gh follow suggest trending --language go
```

### Example 3: Monthly Cleanup

```bash
# Monthly routine
# 1. Backup
gh follow export --output ~/backups/follows-$(date +%Y%m).json

# 2. Find inactive users
gh follow suggest inactive --days 90

# 3. Remove unwanted follows
gh follow remove inactive-user --force

# 4. Sync
gh follow sync --gist
```

### Example 4: Team Onboarding

```bash
# New team member setup
# 1. Import team list
gh follow batch import team-follows.txt --follow

# 2. Tag by team
gh follow add team-lead --tags "team,lead"

# 3. Setup cloud sync
gh follow gist create
```

### Example 5: Interest-Based Curation

```bash
# Organize by interests
gh follow add go-expert --tags "golang,expert,blog"
gh follow add rust-dev --tags "rust,developer"
gh follow add ml-researcher --tags "ml,ai,research"

# List by interest
gh follow list --tag golang
gh follow list --tag ml
```

## Tips & Tricks

### Keyboard Shortcuts

Use shell aliases for common commands:

```bash
# Add to ~/.bashrc or ~/.zshrc
alias gfa='gh follow add'
alias gfr='gh follow remove'
alias gfl='gh follow list'
alias gfs='gh follow sync'
alias gfS='gh follow stats'
```

### Tab Completion

Add completion for usernames:

```bash
# Bash
complete -W "$(gh follow list --format simple 2>/dev/null)" gh-follow

# Zsh
compadd $(gh follow list --format simple 2>/dev/null)
```

### Cron Jobs

Set up automatic sync:

```bash
# Add to crontab (crontab -e)
# Sync every hour
0 * * * * /usr/local/bin/gh-follow autosync trigger --gist >> /tmp/gh-follow.log 2>&1
```

## Next Steps

- [Configuration Guide](configuration.md) - Customize settings
- [Architecture Overview](../ARCHITECTURE.md) - Understand internals
- [Development Guide](../DEVELOPMENT.md) - Contribute to the project
