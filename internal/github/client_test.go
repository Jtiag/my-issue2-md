// Package github provides GitHub API interaction functionality.
package github

import (
	"context"
	"os"
	"testing"
)

// TestFetchIssue_Integration is an integration test for FetchIssue.
// It uses the real GitHub API and requires GITHUB_TOKEN to be set.
func TestFetchIssue_Integration(t *testing.T) {
	// Skip in short mode (unit tests only)
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Check for GITHUB_TOKEN
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("GITHUB_TOKEN not set, skipping integration test")
	}

	// Create client
	client := NewClient(token)

	t.Run("fetch public issue", func(t *testing.T) {
		ctx := context.Background()

		// Use a known issue from golang/go repo
		// Issue 40655 is a real, simple issue
		issue, err := client.FetchIssue(ctx, "golang", "go", 40655)
		if err != nil {
			t.Fatalf("FetchIssue() failed: %v", err)
		}

		// Validate basic fields
		if issue == nil {
			t.Fatal("FetchIssue() returned nil issue")
		}

		if issue.Number != 40655 {
			t.Errorf("Issue.Number = %d, want 40655", issue.Number)
		}

		if issue.Title == "" {
			t.Error("Issue.Title is empty")
		}

		if issue.Author == "" {
			t.Error("Issue.Author is empty")
		}

		if issue.State != "open" && issue.State != "closed" {
			t.Errorf("Issue.State = %s, want 'open' or 'closed'", issue.State)
		}

		if issue.CreatedAt.IsZero() {
			t.Error("Issue.CreatedAt is zero")
		}

		// Labels should be present (may be empty)
		if issue.Labels == nil {
			t.Error("Issue.Labels is nil, should be initialized")
		}

		// Comments should be present (may be empty)
		if issue.Comments == nil {
			t.Error("Issue.Comments is nil, should be initialized")
		}
	})

	t.Run("fetch issue with comments", func(t *testing.T) {
		ctx := context.Background()

		// Use an issue known to have comments
		// Issue 41266 from golang/go typically has multiple comments
		issue, err := client.FetchIssue(ctx, "golang", "go", 41266)
		if err != nil {
			t.Fatalf("FetchIssue() failed: %v", err)
		}

		if issue == nil {
			t.Fatal("FetchIssue() returned nil issue")
		}

		// This issue should have comments
		if len(issue.Comments) == 0 {
			t.Error("Expected comments on issue 41266, got none")
		}

		// Verify comment fields
		for i, comment := range issue.Comments {
			if comment.ID == 0 {
				t.Errorf("Comment %d has ID=0", i)
			}
			if comment.Author == "" {
				t.Errorf("Comment %d has empty Author", i)
			}
			if comment.CreatedAt.IsZero() {
				t.Errorf("Comment %d has zero CreatedAt", i)
			}
		}
	})

	t.Run("fetch issue with labels", func(t *testing.T) {
		ctx := context.Background()

		// Use an issue known to have labels
		issue, err := client.FetchIssue(ctx, "golang", "go", 40655)
		if err != nil {
			t.Fatalf("FetchIssue() failed: %v", err)
		}

		if issue == nil {
			t.Fatal("FetchIssue() returned nil issue")
		}

		// Verify label fields
		for i, label := range issue.Labels {
			if label.Name == "" {
				t.Errorf("Label %d has empty Name", i)
			}
		}
	})
}

// TestFetchPullRequest_Integration is an integration test for FetchPullRequest.
func TestFetchPullRequest_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("GITHUB_TOKEN not set, skipping integration test")
	}

	client := NewClient(token)

	t.Run("fetch merged PR with files", func(t *testing.T) {
		ctx := context.Background()

		// Use a known PR from rust-lang/rust repo
		// PR 109331 is a real PR
		pr, err := client.FetchPullRequest(ctx, "rust-lang", "rust", 109331)
		if err != nil {
			t.Fatalf("FetchPullRequest() failed: %v", err)
		}

		if pr == nil {
			t.Fatal("FetchPullRequest() returned nil PR")
		}

		// Validate basic fields
		if pr.Title == "" {
			t.Error("PR.Title is empty")
		}

		// Validate PR-specific fields
		if pr.HeadBranch == "" {
			t.Error("PR.HeadBranch is empty")
		}

		if pr.BaseBranch == "" {
			t.Error("PR.BaseBranch is empty")
		}

		// Files should be present
		if pr.Files == nil {
			t.Error("PR.Files is nil, should be initialized")
		}

		// Verify file fields
		for i, file := range pr.Files {
			if file.Path == "" {
				t.Errorf("File %d has empty Path", i)
			}
		}
	})
}

// TestNewClient creates a client with token.
func TestNewClient(t *testing.T) {
	client := NewClient("test-token")

	if client == nil {
		t.Fatal("NewClient() returned nil")
	}

	if client.token != "test-token" {
		t.Errorf("Client.token = %q, want %q", client.token, "test-token")
	}
}
