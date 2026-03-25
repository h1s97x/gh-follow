# Installation Guide

Complete guide to installing and setting up gh-follow.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation Methods](#installation-methods)
- [Post-Installation Setup](#post-installation-setup)
- [Verification](#verification)
- [Troubleshooting](#troubleshooting)
- [Uninstallation](#uninstallation)

## Prerequisites

Before installing gh-follow, ensure you have the following:

### Required

| Requirement | Version | How to Check | Install |
|-------------|---------|--------------|---------|
| GitHub CLI | 2.0+ | `gh --version` | [Install Guide](https://cli.github.com/) |
| GitHub Account | - | - | [Sign Up](https://github.com) |

### For Building from Source

| Requirement | Version | How to Check | Install |
|-------------|---------|--------------|---------|
| Go | 1.22+ | `go version` | [Install Guide](https://golang.org/doc/install) |
| Git | 2.0+ | `git --version` | [Install Guide](https://git-scm.com/) |
| Make | - | `make --version` | Package manager |

## Installation Methods

### Method 1: GitHub CLI Extension (Recommended)

The easiest way to install gh-follow is as a GitHub CLI extension.

```bash
# Install the extension
gh extension install h1s97x/gh-follow

# Verify installation
gh follow --version
```

**Advantages:**
- Automatic updates via `gh extension upgrade`
- Integrated with GitHub CLI
- Simple installation

### Method 2: Download Binary

Download the latest binary for your platform from [GitHub Releases](https://github.com/h1s97x/gh-follow/releases).

#### Linux

```bash
# Download (replace VERSION and ARCH with appropriate values)
VERSION="1.0.0"
ARCH="amd64"  # or "arm64"

curl -sSL https://github.com/h1s97x/gh-follow/releases/download/v${VERSION}/gh-follow-linux-${ARCH}.tar.gz | tar xz

# Move to PATH
sudo mv gh-follow /usr/local/bin/

# Make executable
sudo chmod +x /usr/local/bin/gh-follow
```

#### macOS

```bash
# Using Homebrew (if available)
brew install h1s97x/tap/gh-follow

# Or download binary
VERSION="1.0.0"
ARCH="arm64"  # Use "amd64" for Intel Macs

curl -sSL https://github.com/h1s97x/gh-follow/releases/download/v${VERSION}/gh-follow-darwin-${ARCH}.tar.gz | tar xz

# Move to PATH
sudo mv gh-follow /usr/local/bin/
```

#### Windows

```powershell
# Using PowerShell
$VERSION = "1.0.0"
$ARCH = "amd64"

# Download
Invoke-WebRequest -Uri "https://github.com/h1s97x/gh-follow/releases/download/v$VERSION/gh-follow-windows-$ARCH.zip" -OutFile "gh-follow.zip"

# Extract
Expand-Archive -Path "gh-follow.zip" -DestinationPath "."

# Move to PATH (requires admin)
Move-Item -Path "gh-follow.exe" -Destination "C:\Program Files\gh-follow\"
```

### Method 3: Build from Source

```bash
# Clone the repository
git clone https://github.com/h1s97x/gh-follow.git
cd gh-follow

# Build
make build

# Install locally
make install

# Or install to system PATH
sudo make install PREFIX=/usr/local
```

### Method 4: Go Install

```bash
# Install directly
go install github.com/h1s97x/gh-follow@latest

# Binary will be at: $(go env GOPATH)/bin/gh-follow
```

## Post-Installation Setup

### 1. GitHub CLI Authentication

gh-follow uses GitHub CLI's authentication, so ensure you're logged in:

```bash
# Check authentication status
gh auth status

# If not authenticated, log in
gh auth login

# Follow the prompts:
# ? What account do you want to log into? GitHub.com
# ? What is your preferred protocol for Git operations? HTTPS
# ? Authenticate Git with your GitHub credentials? Yes
# ? How would you like to authenticate GitHub CLI? Login with a web browser
```

### 2. Initialize Configuration

```bash
# Initialize with default settings
gh follow config show

# Optionally customize configuration
gh follow config set display.default_format json
gh follow config set sync.auto_sync true
```

### 3. Verify Setup

```bash
# Check version
gh follow --version

# Test basic functionality
gh follow list

# Show help
gh follow --help
```

## Verification

### Check Installation

```bash
# Verify binary exists
which gh-follow

# Check version
gh-follow --version

# Test API connectivity
gh follow stats
```

### Expected Output

```
$ gh follow --version
gh-follow version 1.0.0

$ gh follow stats

📊 Follow List Statistics
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Total Follows:    0
Last Updated:    2024-03-25 10:30:00
Tags Used:       0
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## Troubleshooting

### Common Issues

#### Issue: "gh: command not found"

**Cause:** GitHub CLI is not installed or not in PATH.

**Solution:**
```bash
# Install GitHub CLI
# macOS
brew install gh

# Linux (Ubuntu/Debian)
sudo apt install gh

# Linux (Fedora/RHEL)
sudo dnf install gh

# Windows
winget install GitHub.cli
```

#### Issue: "gh-follow: command not found"

**Cause:** gh-follow is not installed or not in PATH.

**Solution:**
```bash
# Check if installed
gh extension list | grep follow

# If not, install it
gh extension install h1s97x/gh-follow

# Or check PATH
echo $PATH | grep -o "$(go env GOPATH)/bin"
```

#### Issue: "token not found" or "authentication required"

**Cause:** GitHub CLI is not authenticated.

**Solution:**
```bash
# Authenticate
gh auth login

# Verify authentication
gh auth status
```

#### Issue: "permission denied" when writing files

**Cause:** Insufficient permissions for config directory.

**Solution:**
```bash
# Check permissions
ls -la ~/.config/gh/

# Fix permissions
chmod 755 ~/.config/gh/
chmod 600 ~/.config/gh/follow-list.json
```

#### Issue: "rate limit exceeded"

**Cause:** GitHub API rate limit reached.

**Solution:**
```bash
# Check rate limit status
gh api rate_limit

# Wait for reset (authenticated: 5000/hr, unauthenticated: 60/hr)
# Or use authenticated requests
gh auth refresh -h github.com
```

#### Issue: "connection refused" or network errors

**Cause:** Network connectivity issues.

**Solution:**
```bash
# Check network
ping github.com

# Check proxy settings
echo $HTTP_PROXY
echo $HTTPS_PROXY

# Test API
gh api user
```

### Debug Mode

Enable debug logging for more information:

```bash
# Enable debug
export GH_FOLLOW_DEBUG=true

# Run with verbose output
gh follow --help
```

### Verbose Output

```bash
# Run with verbose flag (if available)
gh follow list --verbose
```

### Log Files

Check log files for errors:

```bash
# GitHub CLI logs
cat ~/.config/gh/logs/gh-follow.log

# System logs (Linux/macOS)
journalctl --user -u gh-follow
```

## Platform-Specific Notes

### Linux

- Install via package manager or download binary
- Ensure `~/.local/bin` is in PATH for user installations
- May need to install `ca-certificates` for HTTPS

```bash
# Ubuntu/Debian
sudo apt update && sudo apt install ca-certificates

# RHEL/Fedora
sudo dnf install ca-certificates
```

### macOS

- Homebrew installation recommended
- For Apple Silicon (M1/M2), use `arm64` binaries
- May need to allow the binary in Security settings

```bash
# If blocked by security
xattr -d com.apple.quarantine /usr/local/bin/gh-follow
```

### Windows

- PowerShell 5.1+ or PowerShell Core recommended
- May need to unblock downloaded files

```powershell
# Unblock downloaded file
Unblock-File -Path gh-follow.exe

# Or via Properties: Right-click > Properties > Unblock
```

## Uninstallation

### Remove Extension

```bash
# Remove GitHub CLI extension
gh extension remove h1s97x/gh-follow
```

### Remove Binary

```bash
# Remove installed binary
sudo rm /usr/local/bin/gh-follow

# Or from Go bin
rm $(go env GOPATH)/bin/gh-follow
```

### Clean Up Data

```bash
# Remove configuration and data files
rm -rf ~/.config/gh/follow-*.json

# Or keep your data and just remove the binary
```

## Upgrading

### GitHub CLI Extension

```bash
# Upgrade extension
gh extension upgrade h1s97x/gh-follow

# Upgrade all extensions
gh extension upgrade --all
```

### Binary Installation

```bash
# Download and replace the binary
# Follow the installation steps for your platform
```

### Build from Source

```bash
cd gh-follow
git pull origin main
make build
make install
```

## Next Steps

After installation, continue with:

- [Usage Guide](usage.md) - Learn how to use all features
- [Configuration Guide](configuration.md) - Customize settings
- [Architecture Overview](../ARCHITECTURE.md) - Understand the internals

## Getting Help

- **Documentation**: [GitHub Wiki](https://github.com/h1s97x/gh-follow/wiki)
- **Issues**: [GitHub Issues](https://github.com/h1s97x/gh-follow/issues)
- **Discussions**: [GitHub Discussions](https://github.com/h1s97x/gh-follow/discussions)
