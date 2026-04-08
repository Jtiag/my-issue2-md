// Package github provides GitHub API interaction functionality.
//
// It defines data structures for GitHub resources and clients
// for fetching issues, pull requests, and discussions.
package github

import (
	"context"
	"fmt"

	gh "github.com/google/go-github/v56/github"
)

// Client is a GitHub API client.
type Client struct {
	token        string
	githubClient *gh.Client
}

// NewClient creates a new GitHub API client.
// The token is used for authentication with the GitHub API.
func NewClient(token string) *Client {
	return &Client{
		token:        token,
		githubClient: gh.NewClient(nil).WithAuthToken(token),
	}
}

// FetchIssue fetches an issue from GitHub.
// It retrieves the issue information along with labels, milestone, reactions, and all comments.
func (c *Client) FetchIssue(ctx context.Context, owner, repo string, number int) (*Issue, error) {
	// Fetch the issue
	ghIssue, _, err := c.githubClient.Issues.Get(ctx, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issue %s/%s#%d: %w", owner, repo, number, err)
	}

	// Convert to internal Issue
	issue := &Issue{
		Owner:     owner,
		Repo:      repo,
		Number:    ghIssue.GetNumber(),
		Title:     ghIssue.GetTitle(),
		Body:      ghIssue.GetBody(),
		State:     ghIssue.GetState(),
		Author:    ghIssue.GetUser().GetLogin(),
		CreatedAt: ghIssue.GetCreatedAt().Time,
		UpdatedAt: ghIssue.GetUpdatedAt().Time,
		Labels:    make([]Label, 0),
		Comments:  make([]Comment, 0),
	}

	// Convert labels
	for _, ghLabel := range ghIssue.Labels {
		issue.Labels = append(issue.Labels, Label{
			Name:  ghLabel.GetName(),
			Color: ghLabel.GetColor(),
		})
	}

	// Convert milestone
	if ghIssue.Milestone != nil {
		issue.Milestone = &Milestone{
			Title: ghIssue.Milestone.GetTitle(),
			State: ghIssue.Milestone.GetState(),
		}
	}

	// Convert reactions
	if ghIssue.Reactions != nil {
		issue.Reactions = convertReactions(ghIssue.Reactions)
	}

	// Fetch comments
	issue.Comments, err = c.fetchComments(ctx, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments for issue %s/%s#%d: %w", owner, repo, number, err)
	}

	return issue, nil
}

// FetchPullRequest fetches a pull request from GitHub.
// It retrieves the PR information along with files, patch, and comments.
func (c *Client) FetchPullRequest(ctx context.Context, owner, repo string, number int) (*PullRequest, error) {
	// Fetch the PR
	ghPR, _, err := c.githubClient.PullRequests.Get(ctx, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PR %s/%s#%d: %w", owner, repo, number, err)
	}

	// Convert to internal PullRequest (which embeds Issue)
	pr := &PullRequest{
		Issue: Issue{
			Owner:     owner,
			Repo:      repo,
			Number:    ghPR.GetNumber(),
			Title:     ghPR.GetTitle(),
			Body:      ghPR.GetBody(),
			State:     ghPR.GetState(),
			Author:    ghPR.GetUser().GetLogin(),
			CreatedAt: ghPR.GetCreatedAt().Time,
			UpdatedAt: ghPR.GetUpdatedAt().Time,
			Labels:    make([]Label, 0),
			Comments:  make([]Comment, 0),
		},
		Mergeable:    ghPR.GetMergeable(),
		Merged:       ghPR.GetMerged(),
		Additions:    ghPR.GetAdditions(),
		Deletions:    ghPR.GetDeletions(),
		ChangedFiles: ghPR.GetChangedFiles(),
		HeadBranch:   ghPR.GetHead().GetRef(),
		BaseBranch:   ghPR.GetBase().GetRef(),
		Files:        make([]File, 0),
	}

	// Set MergedAt if merged
	if pr.Merged && ghPR.MergedAt != nil {
		mergedAt := ghPR.MergedAt.Time
		pr.MergedAt = &mergedAt
	}

	// Convert labels
	for _, ghLabel := range ghPR.Labels {
		pr.Labels = append(pr.Labels, Label{
			Name:  ghLabel.GetName(),
			Color: ghLabel.GetColor(),
		})
	}

	// Convert milestone
	if ghPR.Milestone != nil {
		pr.Milestone = &Milestone{
			Title: ghPR.Milestone.GetTitle(),
			State: ghPR.Milestone.GetState(),
		}
	}

	// Fetch files
	listOpts := &gh.ListOptions{PerPage: 100}
	for {
		ghFiles, resp, err := c.githubClient.PullRequests.ListFiles(ctx, owner, repo, number, listOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch PR files: %w", err)
		}

		for _, ghFile := range ghFiles {
			pr.Files = append(pr.Files, File{
				Path:      ghFile.GetFilename(),
				Additions: ghFile.GetAdditions(),
				Deletions: ghFile.GetDeletions(),
				Patch:     ghFile.GetPatch(),
				BlobURL:   ghFile.GetBlobURL(),
			})
		}

		if resp.NextPage == 0 {
			break
		}
		listOpts.Page = resp.NextPage
	}

	// Fetch comments
	pr.Comments, err = c.fetchComments(ctx, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments for PR %s/%s#%d: %w", owner, repo, number, err)
	}

	return pr, nil
}

// convertReactions converts a go-github Reactions object to our internal Reactions type.
func convertReactions(r *gh.Reactions) Reactions {
	return Reactions{
		ThumbsUp:   r.GetPlusOne(),
		ThumbsDown: r.GetMinusOne(),
		Laugh:      r.GetLaugh(),
		Hooray:     r.GetHooray(),
		Confused:   r.GetConfused(),
		Heart:      r.GetHeart(),
		Rocket:     r.GetRocket(),
		Eyes:       r.GetEyes(),
	}
}

// fetchComments fetches all comments for an issue or PR with pagination.
func (c *Client) fetchComments(ctx context.Context, owner, repo string, number int) ([]Comment, error) {
	var comments []Comment

	opts := &gh.IssueListCommentsOptions{
		ListOptions: gh.ListOptions{PerPage: 100},
	}

	for {
		ghComments, resp, err := c.githubClient.Issues.ListComments(ctx, owner, repo, number, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch comments: %w", err)
		}

		for _, ghComment := range ghComments {
			comment := Comment{
				ID:        ghComment.GetID(),
				Author:    ghComment.GetUser().GetLogin(),
				Body:      ghComment.GetBody(),
				CreatedAt: ghComment.GetCreatedAt().Time,
			}

			if ghComment.Reactions != nil {
				comment.Reactions = convertReactions(ghComment.Reactions)
			}

			comments = append(comments, comment)
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return comments, nil
}
