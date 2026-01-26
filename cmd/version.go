package cmd

import (
	"fmt"
	"runtime"
)

var (
	// Version is the semantic version number
	Version = "0.2.0"

	// Commit is the git commit hash (set at build time via -ldflags)
	Commit = "unknown"

	// Date is the build date (set at build time via -ldflags)
	Date = "unknown"
)

// VersionInfo returns a formatted version string with build information
func VersionInfo() string {
	info := fmt.Sprintf("cmd version %s", Version)
	if Commit != "unknown" {
		info += fmt.Sprintf("\ncommit: %s", Commit)
	}
	if Date != "unknown" {
		info += fmt.Sprintf("\nbuilt at: %s", Date)
	}
	info += fmt.Sprintf("\ngo: %s", runtime.Version())
	return info
}
