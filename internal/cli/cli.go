// Package cli provides the command-line interface for issue2md.
package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bigwhite/my-issue2md/internal/config"
	"github.com/bigwhite/my-issue2md/internal/converter"
	"github.com/bigwhite/my-issue2md/internal/github"
	"github.com/bigwhite/my-issue2md/internal/parser"
)

// Execute is the main entry point for the CLI.
// It reads from stdin, writes to stdout and stderr, and processes the given arguments.
// Returns an exit code (0-6 as defined in exit_codes.go).
func Execute(stdin io.Reader, stdout, stderr io.Writer, args []string) int {
	// Step 1: Parse flags
	cfg, err := config.ParseFlags(args)
	if err != nil {
		fmt.Fprintln(stderr, "Error:", err)
		return ExitInvalidURL
	}

	// Step 2: Handle help/version
	if cfg.IsHelpRequested() {
		fmt.Fprintln(stdout, config.HelpText())
		return ExitSuccess
	}
	if cfg.IsVersionRequested() {
		fmt.Fprintln(stdout, config.VersionInfo())
		return ExitSuccess
	}

	// Step 3: Validate config
	if err := cfg.Validate(); err != nil {
		fmt.Fprintln(stderr, "Error:", err)
		return ExitInvalidURL
	}

	// Get the URL from positional args
	// ParseFlags already validated we have at least one positional arg
	// We need to get the raw URL - reparse to get it
	url, err := extractURL(args)
	if err != nil {
		fmt.Fprintln(stderr, "Error:", err)
		return ExitInvalidURL
	}

	// Step 4: Parse URL
	parsed, err := parser.ParseURL(url)
	if err != nil {
		fmt.Fprintln(stderr, "Error: invalid URL:", err)
		return ExitInvalidURL
	}

	// Step 5: Create GitHub client and check token
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Fprintln(stderr, "Error: GITHUB_TOKEN environment variable is required")
		return ExitTokenMissing
	}
	client := github.NewClient(token)

	// Step 6: Fetch resource from GitHub
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var markdown string
	opts := cfg.OutputOptions()

	switch parsed.Type {
	case parser.TypeIssue:
		issue, err := client.FetchIssue(ctx, parsed.Owner, parsed.Repo, parsed.Number)
		if err != nil {
			fmt.Fprintln(stderr, "Error fetching issue:", err)
			return ExitAPIFailed
		}
		markdown, err = converter.IssueToMarkdown(issue, opts)
		if err != nil {
			fmt.Fprintln(stderr, "Error converting to markdown:", err)
			return ExitAPIFailed
		}

	case parser.TypePullRequest:
		pr, err := client.FetchPullRequest(ctx, parsed.Owner, parsed.Repo, parsed.Number)
		if err != nil {
			fmt.Fprintln(stderr, "Error fetching pull request:", err)
			return ExitAPIFailed
		}
		markdown, err = converter.PullRequestToMarkdown(pr, opts)
		if err != nil {
			fmt.Fprintln(stderr, "Error converting to markdown:", err)
			return ExitAPIFailed
		}

	case parser.TypeDiscussion:
		// Discussions require GraphQL API - not yet implemented
		fmt.Fprintln(stderr, "Error: discussions are not yet supported")
		return ExitAPIFailed

	default:
		fmt.Fprintln(stderr, "Error: unsupported resource type:", parsed.Type)
		return ExitInvalidURL
	}

	// Step 7: Write output
	if cfg.OutputFile != "" {
		err := os.WriteFile(cfg.OutputFile, []byte(markdown), 0644)
		if err != nil {
			fmt.Fprintln(stderr, "Error writing to file:", err)
			return ExitWriteFailed
		}
	} else {
		fmt.Fprint(stdout, markdown)
	}

	return ExitSuccess
}

// extractURL extracts the GitHub URL from command-line arguments.
func extractURL(args []string) (string, error) {
	for _, arg := range args {
		if len(arg) > 0 && arg[0] != '-' {
			return arg, nil
		}
	}
	return "", fmt.Errorf("GitHub URL not found in arguments")
}

// Run executes the CLI with standard OS streams.
// It calls os.Exit() with the appropriate exit code.
func Run() int {
	return Execute(os.Stdin, os.Stdout, os.Stderr, os.Args[1:])
}
