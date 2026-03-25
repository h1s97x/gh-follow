package main

import "fmt"

// Version is the current version of gh-follow
// This is set at build time via ldflags
var Version = "1.0.0"

// Build information (set at build time)
var (
	BuildDate   = "unknown"
	BuildCommit = "unknown"
)

// GetVersion returns the full version string
func GetVersion() string {
	return Version
}

// GetFullVersion returns the version with build information
func GetFullVersion() string {
	return fmt.Sprintf("%s (commit: %s, built: %s)", Version, BuildCommit, BuildDate)
}
