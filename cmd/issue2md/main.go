// Package main provides the CLI entry point for issue2md.
//
// issue2md converts GitHub Issues, Pull Requests, and Discussions
// to GitHub Flavored Markdown format.
package main

import (
	"os"

	"github.com/bigwhite/my-issue2md/internal/cli"
)

func main() {
	exitCode := cli.Run()
	os.Exit(exitCode)
}
