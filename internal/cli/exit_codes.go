// Package cli provides the CLI entry point for issue2md.
//
// It coordinates between other internal packages to provide
// the complete command-line interface functionality.
package cli

const (
	// ExitSuccess indicates successful execution or help/version request.
	ExitSuccess = 0
	// ExitInvalidURL indicates the URL format is invalid.
	ExitInvalidURL = 1
	// ExitNotFound indicates the GitHub resource was not found.
	ExitNotFound = 2
	// ExitAPIFailed indicates the GitHub API request failed.
	ExitAPIFailed = 3
	// ExitTokenMissing indicates GITHUB_TOKEN environment variable is not set.
	ExitTokenMissing = 4
	// ExitWriteFailed indicates file write operation failed.
	ExitWriteFailed = 5
	// ExitTimeout indicates network timeout occurred.
	ExitTimeout = 6
)
