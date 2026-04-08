// Package converter provides GitHub data to Markdown conversion functionality.
package converter

import (
	"fmt"
	"strings"

	"github.com/bigwhite/my-issue2md/internal/config"
	"github.com/bigwhite/my-issue2md/internal/github"
)

// DiscussionToMarkdown converts a GitHub Discussion to Markdown format.
//
// The output includes:
//   - YAML Front Matter with metadata
//   - Discussion title and content
//   - Category and upvotes
//   - Replies with nested indentation using blockquotes
//
// Nested replies are indented using increasing levels of blockquotes (>).
func DiscussionToMarkdown(disc *github.Discussion, opts config.Options) (string, error) {
	if disc == nil {
		return "", fmt.Errorf("discussion cannot be nil")
	}

	var b strings.Builder

	// Write YAML Front Matter
	b.WriteString("---\n")
	b.WriteString(fmt.Sprintf("title: %q\n", disc.Title))
	b.WriteString(fmt.Sprintf("number: %d\n", disc.Number))
	b.WriteString(fmt.Sprintf("type: \"discussion\"\n"))
	b.WriteString(fmt.Sprintf("author: %q\n", disc.Author))
	b.WriteString(fmt.Sprintf("created_at: %s\n", disc.CreatedAt.Format("2006-01-02T15:04:05Z")))
	b.WriteString(fmt.Sprintf("category: %q\n", disc.Category))
	b.WriteString(fmt.Sprintf("upvotes: %d\n", disc.Upvotes))
	b.WriteString("---\n\n")

	// Write title
	b.WriteString(fmt.Sprintf("# %s\n\n", disc.Title))

	// Write metadata
	b.WriteString(fmt.Sprintf("**Author:** @%s\n\n", disc.Author))
	b.WriteString(fmt.Sprintf("**Category:** %s\n\n", disc.Category))
	b.WriteString(fmt.Sprintf("**Upvotes:** 👍 %d\n\n", disc.Upvotes))
	b.WriteString(fmt.Sprintf("**Created:** %s\n\n", disc.CreatedAt.Format("2006-01-02 15:04:05 UTC")))

	// Write horizontal rule
	b.WriteString("---\n\n")

	// Write discussion body (process for user links if enabled)
	body := disc.Body
	if opts.EnableUserLinks && body != "" {
		body = writeUserLinks(body)
	}
	if body != "" {
		b.WriteString(body)
		b.WriteString("\n\n")
	}

	// Write replies section
	if len(disc.Replies) > 0 {
		b.WriteString("---\n\n")
		b.WriteString("## Replies\n\n")
		writeReplies(&b, disc.Replies, 0, opts.EnableUserLinks)
	}

	return b.String(), nil
}

// writeReplies recursively writes discussion replies with proper nesting.
// Each level of nesting adds a blockquote prefix (>).
func writeReplies(b *strings.Builder, replies []github.DiscussionReply, level int, enableUserLinks bool) {
	for i, reply := range replies {
		// Create blockquote prefix based on nesting level
		prefix := strings.Repeat("> ", level)

		// Write reply header
		b.WriteString(fmt.Sprintf("%s**Reply %d**\n\n", prefix, i+1))
		b.WriteString(fmt.Sprintf("%s**@%s** - %s\n\n", prefix, reply.Author,
			reply.CreatedAt.Format("2006-01-02 15:04:05 UTC")))

		// Write reply body (process for user links if enabled)
		body := reply.Body
		if enableUserLinks && body != "" {
			body = writeUserLinks(body)
		}
		if body != "" {
			// Each line of body needs the prefix
			lines := strings.Split(body, "\n")
			for _, line := range lines {
				b.WriteString(prefix)
				b.WriteString(line)
				b.WriteString("\n")
			}
			b.WriteString("\n")
		}

		// Recursively write nested replies
		if len(reply.Replies) > 0 {
			writeReplies(b, reply.Replies, level+1, enableUserLinks)
		}
	}
}
