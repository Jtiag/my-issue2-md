// Package converter provides GitHub data to Markdown conversion functionality.
package converter

import (
	"strings"
	"testing"
	"time"

	"github.com/bigwhite/my-issue2md/internal/config"
	"github.com/bigwhite/my-issue2md/internal/github"
)

// TestPullRequestToMarkdown is a table-driven test for the PullRequestToMarkdown function.
// It covers various PR conversion scenarios as specified in spec.md section 2.3.
func TestPullRequestToMarkdown(t *testing.T) {
	fixedTime := time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name    string
		pr      *github.PullRequest
		opts    config.Options
		wantErr bool
	}{
		{
			name: "basic PR conversion",
			pr: &github.PullRequest{
				Issue: github.Issue{
					Number:    123,
					Title:     "Test PR",
					Body:      "This is a test PR",
					State:     "open",
					Author:    "testuser",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					Labels:    []github.Label{},
					Comments:  []github.Comment{},
				},
				Mergeable:    true,
				Merged:       false,
				Additions:    100,
				Deletions:    50,
				ChangedFiles: 3,
				HeadBranch:   "feature-branch",
				BaseBranch:   "main",
				Files: []github.File{
					{Path: "README.md", Additions: 10, Deletions: 5},
					{Path: "main.go", Additions: 90, Deletions: 45},
				},
			},
			opts:    config.Options{},
			wantErr: false,
		},
		{
			name: "merged PR",
			pr: &github.PullRequest{
				Issue: github.Issue{
					Number:    456,
					Title:     "Merged PR",
					Body:      "This PR was merged",
					State:     "closed",
					Author:    "mergeuser",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime.Add(24 * time.Hour),
					Labels:    []github.Label{{Name: "merged", Color: "6f42c1"}},
					Comments:  []github.Comment{},
				},
				Mergeable:    true,
				Merged:       true,
				Additions:    200,
				Deletions:    100,
				ChangedFiles: 5,
				HeadBranch:   "feature-2",
				BaseBranch:   "main",
				Files:        []github.File{},
			},
			opts:    config.Options{},
			wantErr: false,
		},
		{
			name: "PR with files and patch",
			pr: &github.PullRequest{
				Issue: github.Issue{
					Number:    789,
					Title:     "PR with Patch",
					Body:      "Changes included",
					State:     "open",
					Author:    "patchuser",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					Labels:    []github.Label{},
					Comments:  []github.Comment{},
				},
				Mergeable:    true,
				Merged:       false,
				Additions:    50,
				Deletions:    25,
				ChangedFiles: 2,
				HeadBranch:   "fix-branch",
				BaseBranch:   "main",
				Files: []github.File{
					{
						Path:      "bugfix.go",
						Additions: 30,
						Deletions: 15,
						Patch:     "@@ -1,5 +1,10 @@\n-old line\n+new line",
					},
				},
			},
			opts:    config.Options{},
			wantErr: false,
		},
		{
			name: "large PR with many files",
			pr: &github.PullRequest{
				Issue: github.Issue{
					Number:    101,
					Title:     "Large PR",
					Body:      "Many changes",
					State:     "open",
					Author:    "largeuser",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					Labels:    []github.Label{},
					Comments:  []github.Comment{},
				},
				Mergeable:    true,
				Merged:       false,
				Additions:    5000,
				Deletions:    2000,
				ChangedFiles: 50,
				HeadBranch:   "large-feature",
				BaseBranch:   "main",
				Files:        generateTestFiles(50),
			},
			opts:    config.Options{},
			wantErr: false,
		},
		{
			name: "PR with comments",
			pr: &github.PullRequest{
				Issue: github.Issue{
					Number:    202,
					Title:     "PR with Comments",
					Body:      "Please review",
					State:     "open",
					Author:    "commentuser",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					Labels:    []github.Label{},
					Comments: []github.Comment{
						{ID: 1, Author: "reviewer1", Body: "LGTM", CreatedAt: fixedTime.Add(time.Hour)},
					},
				},
				Mergeable:    true,
				Merged:       false,
				Additions:    100,
				Deletions:    50,
				ChangedFiles: 3,
				HeadBranch:   "review-branch",
				BaseBranch:   "main",
				Files:        []github.File{},
			},
			opts:    config.Options{},
			wantErr: false,
		},
		{
			name:    "nil PR",
			pr:      nil,
			opts:    config.Options{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PullRequestToMarkdown(tt.pr, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("PullRequestToMarkdown() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.pr != nil {
				// Check YAML front matter
				if !strings.Contains(got, "---") {
					t.Error("PullRequestToMarkdown() output missing YAML front matter")
				}

				// Check title
				if !strings.Contains(got, tt.pr.Title) {
					t.Errorf("PullRequestToMarkdown() output missing title %q", tt.pr.Title)
				}

				// Check PR-specific fields
				if !strings.Contains(got, "Pull Request") {
					t.Error("PullRequestToMarkdown() output should indicate it's a PR")
				}

				// Check branch info
				if !strings.Contains(got, tt.pr.HeadBranch) {
					t.Errorf("PullRequestToMarkdown() output missing head branch %q", tt.pr.HeadBranch)
				}
				if !strings.Contains(got, tt.pr.BaseBranch) {
					t.Errorf("PullRequestToMarkdown() output missing base branch %q", tt.pr.BaseBranch)
				}

				// Check stats
				if !strings.Contains(got, "additions") || !strings.Contains(got, "deletions") {
					t.Error("PullRequestToMarkdown() output missing PR stats")
				}

				// Check files section if files exist
				if len(tt.pr.Files) > 0 {
					if !strings.Contains(got, "Changed Files") {
						t.Error("PullRequestToMarkdown() output missing Files section")
					}
				}

				// Check merged status
				if tt.pr.Merged && !strings.Contains(got, "Merged") {
					t.Error("PullRequestToMarkdown() output should show merged status")
				}
			}
		})
	}
}

// generateTestFiles creates test file data for large PR testing.
func generateTestFiles(count int) []github.File {
	files := make([]github.File, count)
	for i := 0; i < count; i++ {
		files[i] = github.File{
			Path:      generateTestPath(i),
			Additions: 100,
			Deletions: 40,
		}
	}
	return files
}

// generateTestPath generates a test file path.
func generateTestPath(i int) string {
	return "path/to/file" + string(rune('a'+i%26)) + ".go"
}
