# Contributing to GH-Follow

Thank you for your interest in contributing to GH-Follow! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with:

1. A clear title and description
2. Steps to reproduce the issue
3. Expected behavior
4. Actual behavior
5. Your environment (OS, Go version, etc.)

### Suggesting Features

We welcome feature suggestions! Please create an issue with:

1. A clear description of the feature
2. Use case and benefits
3. Possible implementation approach (optional)

### Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`go test ./...`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## Development Setup

```bash
# Clone your fork
git clone https://github.com/h1s97x/gh-follow.git
cd gh-follow

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o gh-follow .
```

## Code Style

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Add comments for exported functions
- Write tests for new functionality

## Project Structure

```
gh-follow/
├── main.go              # CLI entry point
├── version.go           # Version information
├── internal/            # Internal packages
│   ├── models.go        # Data models
│   ├── storage.go       # Storage operations
│   ├── github_api.go    # GitHub API client
│   ├── config.go        # Configuration
│   └── *_test.go        # Tests
├── .github/             # GitHub workflows
├── docs/                # Documentation
└── Makefile             # Build automation
```

## Commit Messages

We follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation changes
- `test:` Test additions/changes
- `refactor:` Code refactoring
- `chore:` Maintenance tasks

Example: `feat: add tag filtering for list command`

## Testing

- Write unit tests for all new functionality
- Aim for high test coverage
- Run tests before submitting PRs: `go test -race -cover ./...`

## Questions?

Feel free to open an issue for any questions or discussions.

Thank you for contributing!
