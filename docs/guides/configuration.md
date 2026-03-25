# Configuration Guide

Complete guide to configuring gh-follow to suit your needs.

## Table of Contents

- [Configuration Overview](#configuration-overview)
- [Configuration File](#configuration-file)
- [Configuration Options](#configuration-options)
- [Setting Configuration](#setting-configuration)
- [Environment Variables](#environment-variables)
- [Advanced Configuration](#advanced-configuration)

## Configuration Overview

gh-follow can be configured through:

1. **Configuration file** (`~/.config/gh/follow-config.json`)
2. **CLI commands** (`gh follow config set/get`)
3. **Environment variables** (for automation)
4. **Command flags** (one-time overrides)

## Configuration File

### Location

The configuration file is located at:

```
~/.config/gh/follow-config.json
```

You can customize this location with the `GH_FOLLOW_CONFIG` environment variable.

### Default Configuration

```json
{
  "version": "1.0.0",
  "storage": {
    "local_path": "~/.config/gh/follow-list.json",
    "use_gist": false,
    "gist_id": ""
  },
  "sync": {
    "auto_sync": true,
    "sync_interval": 3600,
    "last_sync": "0001-01-01T00:00:00Z"
  },
  "display": {
    "default_format": "table",
    "default_sort": "date",
    "default_order": "desc"
  },
  "cache": {
    "enabled": true,
    "ttl": 86400,
    "max_size": 1000
  },
  "suggestions": {
    "enabled": true,
    "max_results": 20
  },
  "batch": {
    "concurrency": 5,
    "rate_limit": 100
  }
}
```

## Configuration Options

### Storage Options (`storage`)

Control how and where your follow list is stored.

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `local_path` | string | `~/.config/gh/follow-list.json` | Path to local follow list |
| `use_gist` | bool | `false` | Enable Gist synchronization |
| `gist_id` | string | `""` | Gist ID for cloud sync |

**Examples:**

```bash
# Enable Gist sync
gh follow config set storage.use_gist true

# Set custom local path
gh follow config set storage.local_path ~/my-data/follows.json

# Set Gist ID after creating one
gh follow config set storage.gist_id abc123def456
```

### Sync Options (`sync`)

Control synchronization behavior.

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `auto_sync` | bool | `true` | Enable automatic synchronization |
| `sync_interval` | int | `3600` | Sync interval in seconds (1 hour) |
| `last_sync` | timestamp | - | Last sync time (auto-updated) |
| `conflict_strategy` | string | `newest-wins` | Conflict resolution strategy |

**Examples:**

```bash
# Enable auto-sync
gh follow config set sync.auto_sync true

# Set sync interval to 30 minutes
gh follow config set sync.sync_interval 1800

# Set conflict resolution strategy
gh follow config set sync.conflict_strategy local-wins
```

### Display Options (`display`)

Control default output format.

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `default_format` | string | `table` | Default output format: table, json, simple |
| `default_sort` | string | `date` | Default sort field: date, name |
| `default_order` | string | `desc` | Default sort order: asc, desc |
| `color_output` | bool | `true` | Enable colored output |
| `show_tags` | bool | `true` | Show tags in list output |

**Examples:**

```bash
# Set default format to JSON
gh follow config set display.default_format json

# Sort by name ascending by default
gh follow config set display.default_sort name
gh follow config set display.default_order asc

# Disable colors
gh follow config set display.color_output false
```

### Cache Options (`cache`)

Control user information caching.

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | bool | `true` | Enable caching |
| `ttl` | int | `86400` | Cache TTL in seconds (24 hours) |
| `max_size` | int | `1000` | Maximum cached users |

**Examples:**

```bash
# Disable cache
gh follow config set cache.enabled false

# Set cache TTL to 12 hours
gh follow config set cache.ttl 43200

# Increase cache size
gh follow config set cache.max_size 2000
```

### Suggestions Options (`suggestions`)

Control follow suggestions.

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | bool | `true` | Enable suggestions |
| `max_results` | int | `20` | Maximum suggestions to show |
| `include_org_members` | bool | `true` | Include org members in suggestions |
| `include_star_contributors` | bool | `true` | Include starred repo contributors |

**Examples:**

```bash
# Increase suggestion limit
gh follow config set suggestions.max_results 50

# Exclude org members from suggestions
gh follow config set suggestions.include_org_members false
```

### Batch Options (`batch`)

Control batch operations.

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `concurrency` | int | `5` | Number of concurrent operations |
| `rate_limit` | int | `100` | Rate limit delay in milliseconds |

**Examples:**

```bash
# Increase concurrency
gh follow config set batch.concurrency 10

# Increase rate limit delay
gh follow config set batch.rate_limit 200
```

## Setting Configuration

### Using CLI Commands

#### View Configuration

```bash
# Show all configuration
gh follow config show

# Get specific value
gh follow config get display.default_format
gh follow config get sync.sync_interval
```

#### Set Configuration

```bash
# Set a single value
gh follow config set display.default_format json

# Set nested value
gh follow config set storage.use_gist true
```

#### Reset Configuration

```bash
# Reset to defaults (with confirmation)
gh follow config reset

# Force reset without confirmation
gh follow config reset --force
```

### Editing Configuration File

You can directly edit the configuration file:

```bash
# Open in editor
vim ~/.config/gh/follow-config.json

# Or use jq for programmatic changes
jq '.display.default_format = "json"' ~/.config/gh/follow-config.json > tmp.json && mv tmp.json ~/.config/gh/follow-config.json
```

## Environment Variables

Override configuration with environment variables.

### Available Variables

| Variable | Corresponding Config | Example |
|----------|---------------------|---------|
| `GH_FOLLOW_CONFIG` | Config file path | `~/custom-config.json` |
| `GH_FOLLOW_DATA` | storage.local_path | `~/custom-data.json` |
| `GH_FOLLOW_GIST_ID` | storage.gist_id | `abc123` |
| `GH_FOLLOW_FORMAT` | display.default_format | `json` |
| `GH_FOLLOW_NO_COLOR` | display.color_output | `true` |
| `GH_FOLLOW_DEBUG` | Enable debug logging | `true` |
| `GH_FOLLOW_CACHE_TTL` | cache.ttl | `43200` |

### Usage Examples

```bash
# Use custom config file
export GH_FOLLOW_CONFIG=~/my-config.json
gh follow list

# Override data path
export GH_FOLLOW_DATA=~/backups/follows.json
gh follow list

# Disable colors
export GH_FOLLOW_NO_COLOR=true
gh follow list

# Enable debug mode
export GH_FOLLOW_DEBUG=true
gh follow list

# Use in scripts
#!/bin/bash
export GH_FOLLOW_FORMAT=json
export GH_FOLLOW_DATA=./data/follows.json
gh follow list | jq '.[].username'
```

## Advanced Configuration

### Profile-Based Configuration

Use different configurations for different use cases:

```bash
# Work profile
cat > ~/.config/gh/follow-config-work.json << 'EOF'
{
  "storage": {
    "local_path": "~/.config/gh/follow-list-work.json"
  },
  "sync": {
    "auto_sync": true
  }
}
EOF

# Personal profile
cat > ~/.config/gh/follow-config-personal.json << 'EOF'
{
  "storage": {
    "local_path": "~/.config/gh/follow-list-personal.json"
  },
  "sync": {
    "auto_sync": false
  }
}
EOF

# Switch profiles
alias gh-follow-work='GH_FOLLOW_CONFIG=~/.config/gh/follow-config-work.json gh follow'
alias gh-follow-personal='GH_FOLLOW_CONFIG=~/.config/gh/follow-config-personal.json gh follow'
```

### Custom Conflict Resolution

Configure advanced conflict resolution:

```json
{
  "sync": {
    "conflict_strategy": "custom",
    "conflict_rules": {
      "prefer_local_tags": true,
      "prefer_remote_notes": false,
      "merge_tags": true
    }
  }
}
```

### Performance Tuning

Optimize for your needs:

```json
{
  "cache": {
    "enabled": true,
    "ttl": 86400,
    "max_size": 5000,
    "preload": true
  },
  "batch": {
    "concurrency": 10,
    "rate_limit": 50,
    "retry_count": 3,
    "retry_delay": 1000
  },
  "api": {
    "timeout": 30,
    "max_retries": 3
  }
}
```

### Logging Configuration

Configure logging behavior:

```json
{
  "logging": {
    "level": "info",
    "file": "~/.config/gh/follow.log",
    "max_size": 10485760,
    "max_backups": 3,
    "compress": true
  }
}
```

## Configuration Examples

### Example 1: Minimal Configuration

```json
{
  "version": "1.0.0",
  "storage": {
    "local_path": "~/.config/gh/follow-list.json"
  }
}
```

### Example 2: Cloud Sync Enabled

```json
{
  "version": "1.0.0",
  "storage": {
    "local_path": "~/.config/gh/follow-list.json",
    "use_gist": true,
    "gist_id": "abc123def456..."
  },
  "sync": {
    "auto_sync": true,
    "sync_interval": 1800
  }
}
```

### Example 3: Performance Optimized

```json
{
  "version": "1.0.0",
  "storage": {
    "local_path": "~/.config/gh/follow-list.json"
  },
  "cache": {
    "enabled": true,
    "ttl": 172800,
    "max_size": 5000
  },
  "batch": {
    "concurrency": 10,
    "rate_limit": 50
  }
}
```

### Example 4: Automation-Friendly

```json
{
  "version": "1.0.0",
  "storage": {
    "local_path": "~/.config/gh/follow-list.json"
  },
  "display": {
    "default_format": "json"
  },
  "sync": {
    "auto_sync": true,
    "sync_interval": 300
  }
}
```

### Example 5: Development Setup

```json
{
  "version": "1.0.0",
  "storage": {
    "local_path": "./data/follow-list.json",
    "use_gist": false
  },
  "display": {
    "default_format": "table",
    "color_output": true
  },
  "logging": {
    "level": "debug",
    "file": "./logs/follow.log"
  }
}
```

## Troubleshooting Configuration

### Issue: Configuration not loaded

```bash
# Check config file exists
ls -la ~/.config/gh/follow-config.json

# Validate JSON syntax
cat ~/.config/gh/follow-config.json | jq .

# Check permissions
chmod 600 ~/.config/gh/follow-config.json
```

### Issue: Changes not taking effect

```bash
# Verify current configuration
gh follow config show

# Force reload (clear any cached config)
rm -rf ~/.config/gh/.cache
gh follow config show
```

### Issue: Invalid configuration value

```bash
# Reset specific option
gh follow config set display.default_format table

# Or reset all
gh follow config reset --force
```

### Issue: Environment variable not working

```bash
# Verify variable is set
echo $GH_FOLLOW_FORMAT

# Check if it's exported
export GH_FOLLOW_FORMAT=json
gh follow config show
```

## Configuration Best Practices

1. **Backup your configuration**
   ```bash
   cp ~/.config/gh/follow-config.json ~/.config/gh/follow-config.json.backup
   ```

2. **Use version control**
   ```bash
   cd ~/.config/gh
   git init
   git add follow-config.json
   git commit -m "Initial config"
   ```

3. **Document custom settings**
   ```json
   {
     "version": "1.0.0",
     "_comments": {
       "sync_interval": "Set to 30 minutes for frequent sync"
     },
     "sync": {
       "sync_interval": 1800
     }
   }
   ```

4. **Test changes before production**
   ```bash
   export GH_FOLLOW_CONFIG=~/.config/gh/follow-config-test.json
   gh follow list
   ```

5. **Use environment variables for sensitive data**
   ```bash
   # Don't store Gist ID in config for public dotfiles
   export GH_FOLLOW_GIST_ID=your-secret-gist-id
   ```

## Next Steps

- [Usage Guide](usage.md) - Learn how to use all features
- [Installation Guide](installation.md) - Install and setup
- [Architecture Overview](../ARCHITECTURE.md) - Understand the internals
