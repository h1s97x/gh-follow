# Development Guide

This guide covers everything you need to know about developing gh-follow.

## Table of Contents

- [Development Environment Setup](#development-environment-setup)
- [Building the Project](#building-the-project)
- [Testing](#testing)
- [Code Style Guide](#code-style-guide)
- [Debugging Tips](#debugging-tips)
- [Contributing Workflow](#contributing-workflow)

## Development Environment Setup

### Prerequisites

- **Go 1.22+**: Install from [golang.org](https://golang.org/doc/install)
- **GitHub CLI**: Install from [cli.github.com](https://cli.github.com/)
- **Make**: For build automation
- **Git**: For version control

### Initial Setup

1. **Fork and Clone**
   ```bash
   # Fork the repository on GitHub, then:
   git clone https://github.com/h1s97x/gh-follow.git
   cd gh-follow
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   ```

3. **Verify Setup**
   ```bash
   make build
   ./gh-follow --version
   ```

### Development Tools

Recommended VS Code extensions:

```json
{
  "recommendations": [
    "golang.go",
    "ms-azuretools.vscode-docker",
    "redhat.vscode-yaml"
  ]
}
```

## Building the Project

### Using Make

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install locally
make install

# Clean build artifacts
make clean
```

### Manual Build

```bash
# Build binary
go build -o gh-follow .

# Build with version info
go build -ldflags "-X main.Version=$(git describe --tags)" -o gh-follow .
```

### Cross-Platform Builds

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o gh-follow-linux-amd64 .

# macOS
GOOS=darwin GOARCH=amd64 go build -o gh-follow-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -o gh-follow-darwin-arm64 .

# Windows
GOOS=windows GOARCH=amd64 go build -o gh-follow-windows-amd64.exe .
```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run with verbose output
go test ./... -v

# Run specific package
go test ./internal -v

# Run specific test
go test -run TestFollowListAdd ./internal

# Run with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Categories

| Category | Location | Purpose |
|----------|----------|---------|
| Unit Tests | `internal/*_test.go` | Test individual functions |
| Integration Tests | `internal/integration_test.go` | Test workflows |
| Sync Tests | `internal/gist_sync_test.go` | Test sync logic |

### Writing Tests

```go
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"case1", "input1", "output1"},
        {"case2", "input2", "output2"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := MyFunction(tt.input)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## Code Style Guide

### Go Conventions

Follow [Effective Go](https://golang.org/doc/effective_go) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

### Formatting

```bash
# Format all code
gofmt -s -w .

# Or use goimports (recommended)
goimports -w .
```

### Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Packages | lowercase, single word | `internal` |
| Types | PascalCase | `FollowList` |
| Functions | PascalCase (exported) | `NewStorage()` |
| Variables | camelCase | `followList` |
| Constants | PascalCase or UPPER_SNAKE_CASE | `DefaultPath` |

### Error Handling

```go
// Good: Wrap errors with context
if err := storage.Save(list); err != nil {
    return fmt.Errorf("failed to save follow list: %w", err)
}

// Good: Use custom errors
if username == "" {
    return ErrEmptyUsername
}

// Bad: Ignore errors
storage.Save(list) // Don't do this
```

### Documentation

```go
// FollowList represents a collection of follow records.
// It maintains metadata about the list including version,
// last update time, and sync status.
type FollowList struct {
    Version   string    `json:"version"`
    UpdatedAt time.Time `json:"updated_at"`
    Follows   []Follow  `json:"follows"`
    Metadata  Metadata  `json:"metadata"`
}

// Add adds a new follow record to the list.
// If the username already exists, it does nothing.
func (fl *FollowList) Add(username string, notes string, tags []string) {
    // ...
}
```

### Project Structure

```
gh-follow/
├── main.go              # Entry point, CLI definitions
├── version.go           # Version information
├── internal/            # Private packages
│   ├── models.go        # Data models
│   ├── storage.go       # Storage operations
│   ├── github_api.go    # GitHub API client
│   ├── config.go        # Configuration
│   ├── gist_sync.go     # Gist sync logic
│   ├── cache.go         # User caching
│   ├── suggest.go       # Suggestions engine
│   ├── batch.go         # Batch operations
│   └── *_cmd.go         # Command handlers
├── docs/                # Documentation
└── .github/             # Workflows, templates
```

## Debugging Tips

### Enable Debug Logging

```bash
# Set debug environment variable
export GH_FOLLOW_DEBUG=true

# Run with verbose output
./gh-follow --verbose list
```

### Using Delve Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug build
go build -gcflags="all=-N -l" -o gh-follow-debug .

# Start debugger
dlv exec ./gh-follow-debug -- list
```

### Common Issues

#### Issue: "token not found"

```bash
# Check GitHub CLI authentication
gh auth status

# Re-authenticate if needed
gh auth login
```

#### Issue: "rate limit exceeded"

```bash
# Wait for rate limit reset
gh api rate_limit

# Use authenticated requests (higher limits)
gh auth refresh -h github.com
```

#### Issue: "permission denied"

```bash
# Check file permissions
ls -la ~/.config/gh/

# Fix permissions
chmod 600 ~/.config/gh/follow-list.json
```

### Verbose API Calls

```go
// In github_api.go, add logging
func (gc *GitHubClient) GetUser(ctx context.Context, username string) (*github.User, *github.Response, error) {
    if os.Getenv("GH_FOLLOW_DEBUG") == "true" {
        log.Printf("[DEBUG] Fetching user: %s", username)
    }
    return gc.client.Users.Get(ctx, username)
}
```

## Contributing Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Changes

- Write code following the style guide
- Add/update tests
- Update documentation

### 3. Run Checks

```bash
# Format code
gofmt -s -w .

# Run linter
go vet ./...

# Run tests
make test

# Build check
make build
```

### 4. Commit Changes

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```bash
# Features
git commit -m "feat: add batch follow command"

# Bug fixes
git commit -m "fix: handle empty username correctly"

# Documentation
git commit -m "docs: update installation guide"

# Refactoring
git commit -m "refactor: simplify sync logic"
```

### 5. Push and Create PR

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub.

### PR Checklist

- [ ] Code follows style guidelines
- [ ] Tests pass locally
- [ ] Documentation updated
- [ ] Commit messages follow conventions
- [ ] No breaking changes (or documented)

## Release Process

### Creating a Release

1. Update version in `version.go`
2. Update `CHANGELOG.md`
3. Create tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
4. GitHub Actions will build and publish

### Version Naming

- **Major (X.0.0)**: Breaking changes
- **Minor (1.X.0)**: New features
- **Patch (1.0.X)**: Bug fixes

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [urfave/cli Guide](https://cli.urfave.org/)
- [go-github Documentation](https://pkg.go.dev/github.com/google/go-github/v55)
- [GitHub API Reference](https://docs.github.com/en/rest)
