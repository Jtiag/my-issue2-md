// Package converter provides GitHub data to Markdown conversion functionality.
package converter

import (
	"fmt"
	"strings"

	"github.com/bigwhite/my-issue2md/internal/config"
	"github.com/bigwhite/my-issue2md/internal/github"
)

// PullRequestToMarkdown converts a GitHub Pull Request to Markdown format.
//
// The output includes everything from IssueToMarkdown plus:
//   - Branch information (head and base)
//   - Merge status
//   - File changes statistics
//   - Changed files list
//   - Collapsible diff sections
func PullRequestToMarkdown(pr *github.PullRequest, opts config.Options) (string, error) {
	if pr == nil {
		return "", fmt.Errorf("pull request cannot be nil")
	}

	var b strings.Builder

	// Write YAML Front Matter (reuse from issue)
	b.WriteString(writePRYAMLFrontMatter(pr))

	// Write PR title with link
	b.WriteString(fmt.Sprintf("# %s\n\n", pr.Title))

	// Write PR reference
	b.WriteString(fmt.Sprintf("**Pull Request:** [%s/%s#%d](https://github.com/%s/%s/pull/%d)\n\n",
		pr.Owner, pr.Repo, pr.Number, pr.Owner, pr.Repo, pr.Number))
	b.WriteString(fmt.Sprintf("**Author:** @%s\n\n", pr.Author))
	b.WriteString(fmt.Sprintf("**State:** %s\n\n", pr.State))

	// Write branch information
	b.WriteString(fmt.Sprintf("**Branch:** `%s` → `%s`\n\n", pr.HeadBranch, pr.BaseBranch))

	// Write merge status
	if pr.Merged {
		b.WriteString("**Status:** ✅ Merged\n\n")
		if pr.MergedAt != nil {
			b.WriteString(fmt.Sprintf("**Merged at:** %s\n\n", pr.MergedAt.Format("2006-01-02 15:04:05 UTC")))
		}
	} else if pr.Mergeable {
		b.WriteString("**Status:** 🟢 Mergeable\n\n")
	} else {
		b.WriteString("**Status:** 🔴 Not mergeable\n\n")
	}

	// Write PR statistics
	b.WriteString("**Changes:**\n\n")
	b.WriteString(fmt.Sprintf("- **+%d** additions\n", pr.Additions))
	b.WriteString(fmt.Sprintf("- **-%d** deletions\n", pr.Deletions))
	b.WriteString(fmt.Sprintf("- **%d** files changed\n\n", pr.ChangedFiles))

	// Write labels if present
	if len(pr.Labels) > 0 {
		b.WriteString("**Labels:**")
		for _, label := range pr.Labels {
			b.WriteString(fmt.Sprintf(" `%s`", label.Name))
		}
		b.WriteString("\n\n")
	}

	// Write milestone if present
	if pr.Milestone != nil {
		b.WriteString(fmt.Sprintf("**Milestone:** %s\n\n", pr.Milestone.Title))
	}

	// Write horizontal rule
	b.WriteString("---\n\n")

	// Write PR body (process for user links if enabled)
	body := pr.Body
	if opts.EnableUserLinks && body != "" {
		body = writeUserLinks(body)
	}
	if body != "" {
		b.WriteString(body)
		b.WriteString("\n\n")
	}

	// Write files section
	if len(pr.Files) > 0 {
		b.WriteString("---\n\n")
		b.WriteString("## Changed Files\n\n")

		for _, file := range pr.Files {
			b.WriteString(fmt.Sprintf("### %s\n\n", file.Path))
			b.WriteString(fmt.Sprintf("+%d -%d\n\n", file.Additions, file.Deletions))

			// Add collapsible diff if patch exists
			if file.Patch != "" {
				b.WriteString("<details>\n")
				b.WriteString(fmt.Sprintf("<summary>View diff</summary>\n\n"))
				b.WriteString("```diff\n")
				b.WriteString(file.Patch)
				b.WriteString("\n```\n")
				b.WriteString("</details>\n\n")
			}
		}
	}

	// Write comments section
	convOpts := converterOptions{
		enableReactions: opts.EnableReactions,
		enableUserLinks: opts.EnableUserLinks,
	}
	comments := writeComments(pr.Comments, convOpts)
	if comments != "" {
		b.WriteString("---\n\n")
		b.WriteString(comments)
	}

	return b.String(), nil
}

// writePRYAMLFrontMatter generates the YAML Front Matter section for a PR.
func writePRYAMLFrontMatter(pr *github.PullRequest) string {
	var b strings.Builder

	b.WriteString("---\n")
	b.WriteString(fmt.Sprintf("title: %q\n", pr.Title))
	b.WriteString(fmt.Sprintf("number: %d\n", pr.Number))
	b.WriteString(fmt.Sprintf("type: \"pull_request\"\n"))
	b.WriteString(fmt.Sprintf("state: %q\n", pr.State))
	b.WriteString(fmt.Sprintf("author: %q\n", pr.Author))
	b.WriteString(fmt.Sprintf("created_at: %s\n", pr.CreatedAt.Format("2006-01-02T15:04:05Z")))
	b.WriteString(fmt.Sprintf("updated_at: %s\n", pr.UpdatedAt.Format("2006-01-02T15:04:05Z")))
	b.WriteString(fmt.Sprintf("head_branch: %q\n", pr.HeadBranch))
	b.WriteString(fmt.Sprintf("base_branch: %q\n", pr.BaseBranch))
	b.WriteString(fmt.Sprintf("additions: %d\n", pr.Additions))
	b.WriteString(fmt.Sprintf("deletions: %d\n", pr.Deletions))
	b.WriteString(fmt.Sprintf("changed_files: %d\n", pr.ChangedFiles))
	b.WriteString(fmt.Sprintf("merged: %v\n", pr.Merged))
	b.WriteString(fmt.Sprintf("mergeable: %v\n", pr.Mergeable))

	// Add labels if present
	if len(pr.Labels) > 0 {
		b.WriteString("labels:\n")
		for _, label := range pr.Labels {
			b.WriteString(fmt.Sprintf("  - name: %q\n", label.Name))
			b.WriteString(fmt.Sprintf("    color: %q\n", label.Color))
		}
	}

	// Add milestone if present
	if pr.Milestone != nil {
		b.WriteString(fmt.Sprintf("milestone:\n"))
		b.WriteString(fmt.Sprintf("  title: %q\n", pr.Milestone.Title))
		b.WriteString(fmt.Sprintf("  state: %q\n", pr.Milestone.State))
	}

	b.WriteString("---\n\n")

	return b.String()
}
