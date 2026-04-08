// Package config provides configuration management for the CLI.
//
// It handles command-line flag parsing and provides options
// for the markdown converter.
package config

import "fmt"

// Config holds the application configuration parsed from command-line flags.
type Config struct {
	// EnableReactions enables GitHub reactions statistics in output.
	EnableReactions bool
	// EnableUserLinks converts @username to clickable GitHub profile links.
	EnableUserLinks bool
	// OutputFile is the path to the output file.
	// Empty string means write to stdout.
	OutputFile string
	// versionRequested is true when -v or -version flag is provided.
	versionRequested bool
	// helpRequested is true when -h or -help flag is provided.
	helpRequested bool
}

// Options holds converter options derived from Config.
type Options struct {
	// EnableReactions enables GitHub reactions statistics.
	EnableReactions bool
	// EnableUserLinks enables @username to links conversion.
	EnableUserLinks bool
}

// ParseFlags parses command-line arguments and returns a Config.
// It handles flags and positional arguments without using the flag package.
//
// Flags:
//   - -h, -help: Show help text
//   - -v, -version: Show version information
//   - -enable-reactions: Enable GitHub reactions in output
//   - -enable-user-links: Convert @mentions to profile links
//
// Positional arguments:
//  1. GitHub URL (required, unless help/version is requested)
//  2. Output file path (optional, defaults to stdout)
func ParseFlags(args []string) (*Config, error) {
	cfg := &Config{}

	// Parse flags first
	positional := make([]string, 0)
	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch arg {
		case "-h", "-help":
			cfg.helpRequested = true
			return cfg, nil
		case "-v", "-version":
			cfg.versionRequested = true
			return cfg, nil
		case "-enable-reactions":
			cfg.EnableReactions = true
		case "-enable-user-links":
			cfg.EnableUserLinks = true
		default:
			// Check if it's a flag (starts with -)
			if len(arg) > 0 && arg[0] == '-' {
				return nil, fmt.Errorf("unknown flag: %s", arg)
			}
			// It's a positional argument
			positional = append(positional, arg)
		}
	}

	// Validate positional arguments
	if len(positional) == 0 {
		return nil, fmt.Errorf("GitHub URL is required")
	}
	if len(positional) > 2 {
		return nil, fmt.Errorf("too many arguments: expected URL [output_file], got %d arguments", len(positional))
	}

	// Set output file if provided
	if len(positional) == 2 {
		cfg.OutputFile = positional[1]
	}

	return cfg, nil
}

// Validate checks if the configuration is valid.
// Returns an error if required fields are missing.
func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("config cannot be nil")
	}

	// Help and version requests are always valid
	if c.helpRequested || c.versionRequested {
		return nil
	}

	// URL validation would be done by the parser
	// Here we just check the config itself
	return nil
}

// OutputOptions returns converter options derived from the Config.
func (c *Config) OutputOptions() Options {
	return Options{
		EnableReactions: c.EnableReactions,
		EnableUserLinks: c.EnableUserLinks,
	}
}

// VersionInfo returns version information for the CLI.
func VersionInfo() string {
	return "issue2md version 1.0.0"
}

// HelpText returns the help text for the CLI.
func HelpText() string {
	return `issue2md - Convert GitHub Issues, Pull Requests, and Discussions to Markdown

Usage:
  issue2md [flags] <github_url> [output_file]

Arguments:
  github_url    The GitHub URL for an issue, PR, or discussion
  output_file   Optional output file path (defaults to stdout)

Flags:
  -h, -help              Show this help text
  -v, -version           Show version information
  -enable-reactions      Include GitHub reactions in output
  -enable-user-links     Convert @mentions to clickable profile links

Examples:
  issue2md https://github.com/owner/repo/issues/123
  issue2md https://github.com/owner/repo/pull/456 output.md
  issue2md -enable-reactions https://github.com/owner/repo/issues/789

Exit Codes:
  0 - Success
  1 - Invalid URL or missing arguments
  2 - GitHub API error
  3 - Resource not found
  4 - Authentication error (missing GITHUB_TOKEN)
  5 - Network error
  6 - Timeout

Environment Variables:
  GITHUB_TOKEN    GitHub personal access token (required for private repos)
`
}

// IsHelpRequested returns true if help was requested.
func (c *Config) IsHelpRequested() bool {
	return c.helpRequested
}

// IsVersionRequested returns true if version was requested.
func (c *Config) IsVersionRequested() bool {
	return c.versionRequested
}
