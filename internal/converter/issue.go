// Package converter provides GitHub data to Markdown conversion functionality.
//
// It converts GitHub issues, pull requests, and discussions
// into GitHub Flavored Markdown format.
package converter

import (
	"fmt"
	"strings"

	"github.com/bigwhite/my-issue2md/internal/config"
	"github.com/bigwhite/my-issue2md/internal/github"
)

// IssueToMarkdown converts a GitHub Issue to Markdown format.
//
// The output includes:
//   - YAML Front Matter with metadata
//   - Issue title and description
//   - Labels and milestone (if present)
//   - Reactions (if EnableReactions is true)
//   - Comments section (if comments exist)
//
// Usernames are converted to profile links if EnableUserLinks is true.
func IssueToMarkdown(issue *github.Issue, opts config.Options) (string, error) {
	if issue == nil {
		return "", fmt.Errorf("issue cannot be nil")
	}

	var b strings.Builder

	// Write YAML Front Matter
	b.WriteString(writeYAMLFrontMatter(issue))

	// Write issue title with link
	b.WriteString(fmt.Sprintf("# %s\n\n", issue.Title))

	// Write issue reference
	b.WriteString(fmt.Sprintf("**Issue:** [%s/%s#%d](https://github.com/%s/%s/issues/%d)\n\n",
		issue.Owner, issue.Repo, issue.Number, issue.Owner, issue.Repo, issue.Number))
	b.WriteString(fmt.Sprintf("**Author:** @%s\n\n", issue.Author))
	b.WriteString(fmt.Sprintf("**State:** %s\n\n", issue.State))
	b.WriteString(fmt.Sprintf("**Created:** %s\n\n", issue.CreatedAt.Format("2006-01-02 15:04:05 UTC")))
	b.WriteString(fmt.Sprintf("**Updated:** %s\n\n", issue.UpdatedAt.Format("2006-01-02 15:04:05 UTC")))

	// Write labels if present
	if len(issue.Labels) > 0 {
		b.WriteString("**Labels:**")
		for _, label := range issue.Labels {
			b.WriteString(fmt.Sprintf(" `%s`", label.Name))
		}
		b.WriteString("\n\n")
	}

	// Write milestone if present
	if issue.Milestone != nil {
		b.WriteString(fmt.Sprintf("**Milestone:** %s\n\n", issue.Milestone.Title))
	}

	// Write reactions if enabled
	if opts.EnableReactions {
		if r := writeReactions(issue.Reactions); r != "" {
			b.WriteString(r)
			b.WriteString("\n\n")
		}
	}

	// Write horizontal rule before body
	b.WriteString("---\n\n")

	// Write issue body (process for user links if enabled)
	body := issue.Body
	if opts.EnableUserLinks && body != "" {
		body = writeUserLinks(body)
	}
	if body != "" {
		b.WriteString(body)
		b.WriteString("\n\n")
	}

	// Write comments section
	convOpts := converterOptions{
		enableReactions: opts.EnableReactions,
		enableUserLinks: opts.EnableUserLinks,
	}
	comments := writeComments(issue.Comments, convOpts)
	if comments != "" {
		b.WriteString("---\n\n")
		b.WriteString(comments)
	}

	return b.String(), nil
}
