// Package converter provides GitHub data to Markdown conversion functionality.
package converter

import (
	"strings"
	"testing"
	"time"

	"github.com/bigwhite/my-issue2md/internal/config"
	"github.com/bigwhite/my-issue2md/internal/github"
)

// TestIssueToMarkdown is a table-driven test for the IssueToMarkdown function.
// It covers various issue conversion scenarios as specified in spec.md section 2.2.
func TestIssueToMarkdown(t *testing.T) {
	// Create a fixed time for consistent testing
	fixedTime := time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name    string
		issue   *github.Issue
		opts    config.Options
		want    string
		wantErr bool
	}{
		{
			name: "basic issue conversion",
			issue: &github.Issue{
				Number:    123,
				Title:     "Test Issue",
				Body:      "This is a test issue",
				State:     "open",
				Author:    "testuser",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				Labels:    []github.Label{},
				Comments:  []github.Comment{},
			},
			opts:    config.Options{},
			wantErr: false,
		},
		{
			name: "issue with labels and milestone",
			issue: &github.Issue{
				Number:    456,
				Title:     "Feature Request",
				Body:      "Please add this feature",
				State:     "open",
				Author:    "featureuser",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				Labels: []github.Label{
					{Name: "enhancement", Color: "a2eeef"},
					{Name: "good first issue", Color: "7057ff"},
				},
				Milestone: &github.Milestone{
					Title: "v1.0.0",
					State: "open",
				},
				Comments: []github.Comment{},
			},
			opts:    config.Options{},
			wantErr: false,
		},
		{
			name: "issue with reactions enabled",
			issue: &github.Issue{
				Number:    789,
				Title:     "Issue with Reactions",
				Body:      "This issue has reactions",
				State:     "open",
				Author:    "reactionuser",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				Labels:    []github.Label{},
				Reactions: github.Reactions{
					ThumbsUp:   5,
					ThumbsDown: 1,
					Laugh:      2,
					Hooray:     3,
					Confused:   0,
					Heart:      10,
					Rocket:     1,
					Eyes:       2,
				},
				Comments: []github.Comment{},
			},
			opts:    config.Options{EnableReactions: true},
			wantErr: false,
		},
		{
			name: "issue with user links enabled",
			issue: &github.Issue{
				Number:    101,
				Title:     "Issue with User Links",
				Body:      "@mention1 please review this",
				State:     "open",
				Author:    "linkuser",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				Labels:    []github.Label{},
				Comments: []github.Comment{
					{
						ID:        1,
						Author:    "commenter1",
						Body:      "@author I agree",
						CreatedAt: fixedTime,
					},
				},
			},
			opts:    config.Options{EnableUserLinks: true},
			wantErr: false,
		},
		{
			name: "issue with multiple comments",
			issue: &github.Issue{
				Number:    202,
				Title:     "Issue with Comments",
				Body:      "Main issue description",
				State:     "open",
				Author:    "commentuser",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				Labels:    []github.Label{},
				Comments: []github.Comment{
					{
						ID:        1,
						Author:    "user1",
						Body:      "First comment",
						CreatedAt: fixedTime,
					},
					{
						ID:        2,
						Author:    "user2",
						Body:      "Second comment",
						CreatedAt: fixedTime.Add(time.Hour),
					},
					{
						ID:        3,
						Author:    "user3",
						Body:      "Third comment",
						CreatedAt: fixedTime.Add(2 * time.Hour),
					},
				},
			},
			opts:    config.Options{},
			wantErr: false,
		},
		{
			name: "issue with all options enabled",
			issue: &github.Issue{
				Number:    303,
				Title:     "Complete Issue",
				Body:      "Full featured issue @team",
				State:     "closed",
				Author:    "completeuser",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime.Add(24 * time.Hour),
				Labels: []github.Label{
					{Name: "bug", Color: "d73a4a"},
				},
				Milestone: &github.Milestone{
					Title: "v2.0.0",
					State: "closed",
				},
				Reactions: github.Reactions{
					ThumbsUp: 3,
					Heart:    5,
				},
				Comments: []github.Comment{
					{
						ID:        1,
						Author:    "commenter",
						Body:      "Fixed by @developer",
						CreatedAt: fixedTime.Add(time.Hour),
						Reactions: github.Reactions{
							ThumbsUp: 2,
						},
					},
				},
			},
			opts:    config.Options{EnableReactions: true, EnableUserLinks: true},
			wantErr: false,
		},
		{
			name:    "nil issue",
			issue:   nil,
			opts:    config.Options{},
			want:    "",
			wantErr: true,
		},
		{
			name: "issue with empty body",
			issue: &github.Issue{
				Number:    404,
				Title:     "Empty Body Issue",
				Body:      "",
				State:     "open",
				Author:    "emptyuser",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				Labels:    []github.Label{},
				Comments:  []github.Comment{},
			},
			opts:    config.Options{},
			wantErr: false,
		},
		{
			name: "issue with special characters in body",
			issue: &github.Issue{
				Number:    505,
				Title:     "Special Characters Issue",
				Body:      "Body with **bold**, *italic*, `code`, and [links](https://example.com)",
				State:     "open",
				Author:    "specialuser",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				Labels:    []github.Label{},
				Comments:  []github.Comment{},
			},
			opts:    config.Options{},
			wantErr: false,
		},
		{
			name: "issue with multiline body",
			issue: &github.Issue{
				Number:    606,
				Title:     "Multiline Issue",
				Body:      "Line 1\nLine 2\n\nLine 4 (after blank)\n- List item 1\n- List item 2",
				State:     "open",
				Author:    "multilineuser",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				Labels:    []github.Label{},
				Comments:  []github.Comment{},
			},
			opts:    config.Options{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IssueToMarkdown(tt.issue, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("IssueToMarkdown() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// For successful conversions, check that output contains expected elements
			if !tt.wantErr && tt.issue != nil {
				// Check YAML front matter
				if !strings.Contains(got, "---") {
					t.Error("IssueToMarkdown() output missing YAML front matter markers")
				}

				// Check title
				if !strings.Contains(got, tt.issue.Title) {
					t.Errorf("IssueToMarkdown() output missing title %q", tt.issue.Title)
				}

				// Check author
				if !strings.Contains(got, tt.issue.Author) {
					t.Errorf("IssueToMarkdown() output missing author %q", tt.issue.Author)
				}

				// Check state
				if !strings.Contains(got, tt.issue.State) {
					t.Errorf("IssueToMarkdown() output missing state %q", tt.issue.State)
				}

				// Check number
				if !strings.Contains(got, "#") && !strings.Contains(got, "123") {
					t.Error("IssueToMarkdown() output missing issue number reference")
				}

				// Check comments section if comments exist
				if len(tt.issue.Comments) > 0 {
					if !strings.Contains(got, "## Comments") {
						t.Error("IssueToMarkdown() output missing Comments section")
					}
				}

				// Check reactions if enabled
				if tt.opts.EnableReactions && (tt.issue.Reactions.ThumbsUp > 0 || tt.issue.Reactions.Heart > 0) {
					// Should have reactions in output
					if !strings.Contains(got, "👍") && !strings.Contains(got, "❤️") {
						t.Error("IssueToMarkdown() output missing reactions when enabled")
					}
				}

				// Check user links if enabled
				if tt.opts.EnableUserLinks && strings.Contains(tt.issue.Body, "@") {
					if !strings.Contains(got, "](https://github.com/") {
						t.Error("IssueToMarkdown() output missing user links when enabled")
					}
				}
			}

			// For nil issue, check error
			if tt.wantErr && tt.issue == nil {
				if err == nil {
					t.Error("IssueToMarkdown() expected error for nil issue, got nil")
				}
			}
		})
	}
}

