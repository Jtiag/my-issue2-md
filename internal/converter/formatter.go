// Package converter provides GitHub data to Markdown conversion functionality.
package converter

import (
	"fmt"
	"strings"

	"github.com/bigwhite/my-issue2md/internal/github"
)

// writeYAMLFrontMatter generates the YAML Front Matter section.
// It includes issue metadata like title, number, state, author, dates, labels, and milestone.
func writeYAMLFrontMatter(issue *github.Issue) string {
	var b strings.Builder

	b.WriteString("---\n")
	b.WriteString(fmt.Sprintf("title: %q\n", issue.Title))
	b.WriteString(fmt.Sprintf("number: %d\n", issue.Number))
	b.WriteString(fmt.Sprintf("state: %q\n", issue.State))
	b.WriteString(fmt.Sprintf("author: %q\n", issue.Author))
	b.WriteString(fmt.Sprintf("created_at: %s\n", issue.CreatedAt.Format("2006-01-02T15:04:05Z")))
	b.WriteString(fmt.Sprintf("updated_at: %s\n", issue.UpdatedAt.Format("2006-01-02T15:04:05Z")))

	// Add labels if present
	if len(issue.Labels) > 0 {
		b.WriteString("labels:\n")
		for _, label := range issue.Labels {
			b.WriteString(fmt.Sprintf("  - name: %q\n", label.Name))
			b.WriteString(fmt.Sprintf("    color: %q\n", label.Color))
		}
	}

	// Add milestone if present
	if issue.Milestone != nil {
		b.WriteString(fmt.Sprintf("milestone:\n"))
		b.WriteString(fmt.Sprintf("  title: %q\n", issue.Milestone.Title))
		b.WriteString(fmt.Sprintf("  state: %q\n", issue.Milestone.State))
	}

	b.WriteString("---\n\n")

	return b.String()
}

// writeComments formats the comments section.
// Comments are written in chronological order.
func writeComments(comments []github.Comment, opts converterOptions) string {
	if len(comments) == 0 {
		return ""
	}

	var b strings.Builder

	b.WriteString("## Comments\n\n")

	for i, comment := range comments {
		b.WriteString(fmt.Sprintf("### Comment %d\n\n", i+1))
		b.WriteString(fmt.Sprintf("**@%s** commented on %s\n\n", comment.Author, comment.CreatedAt.Format("2006-01-02 15:04:05 UTC")))

		// Process body with options
		body := comment.Body
		if opts.enableUserLinks {
			body = writeUserLinks(body)
		}

		b.WriteString(body)
		b.WriteString("\n\n")

		// Add reactions if enabled
		if opts.enableReactions {
			if r := writeReactions(comment.Reactions); r != "" {
				b.WriteString(r)
				b.WriteString("\n")
			}
		}
	}

	return b.String()
}

// writeReactions formats the reactions summary.
// Returns empty string if there are no reactions.
func writeReactions(reactions github.Reactions) string {
	if reactions.Total() == 0 {
		return ""
	}

	var parts []string

	if reactions.ThumbsUp > 0 {
		parts = append(parts, fmt.Sprintf("👍 %d", reactions.ThumbsUp))
	}
	if reactions.ThumbsDown > 0 {
		parts = append(parts, fmt.Sprintf("👎 %d", reactions.ThumbsDown))
	}
	if reactions.Laugh > 0 {
		parts = append(parts, fmt.Sprintf("😄 %d", reactions.Laugh))
	}
	if reactions.Hooray > 0 {
		parts = append(parts, fmt.Sprintf("🎉 %d", reactions.Hooray))
	}
	if reactions.Confused > 0 {
		parts = append(parts, fmt.Sprintf("😕 %d", reactions.Confused))
	}
	if reactions.Heart > 0 {
		parts = append(parts, fmt.Sprintf("❤️ %d", reactions.Heart))
	}
	if reactions.Rocket > 0 {
		parts = append(parts, fmt.Sprintf("🚀 %d", reactions.Rocket))
	}
	if reactions.Eyes > 0 {
		parts = append(parts, fmt.Sprintf("👀 %d", reactions.Eyes))
	}

	if len(parts) == 0 {
		return ""
	}

	return "**Reactions:** " + strings.Join(parts, " | ")
}

// writeUserLink converts a username to a markdown link.
// Example: "username" -> "[@username](https://github.com/username)"
func writeUserLink(username string) string {
	if username == "" {
		return ""
	}
	return fmt.Sprintf("[@%s](https://github.com/%s)", username, username)
}

// writeUserLinks converts all @mentions in text to markdown links.
// Example: "Hey @user, please review" -> "Hey [@user](https://github.com/user), please review"
func writeUserLinks(text string) string {
	if text == "" {
		return text
	}

	// Split by @ and process each part
	parts := strings.Split(text, "@")
	var result strings.Builder

	result.WriteString(parts[0]) // First part has no @

	for i := 1; i < len(parts); i++ {
		part := parts[i]
		// Extract username (alphanumeric, hyphens, underscores)
		endIdx := 0
		for endIdx < len(part) {
			c := part[endIdx]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' {
				endIdx++
			} else {
				break
			}
		}

		if endIdx > 0 {
			username := part[:endIdx]
			result.WriteString(writeUserLink(username))
			result.WriteString(part[endIdx:])
		} else {
			// Not a valid mention, keep the @
			result.WriteString("@")
			result.WriteString(part)
		}
	}

	return result.String()
}

// converterOptions holds conversion options for internal use.
type converterOptions struct {
	enableReactions bool
	enableUserLinks bool
}

