// Package github provides GitHub API interaction functionality.
//
// It defines data structures for GitHub resources and clients
// for fetching issues, pull requests, and discussions.
package github

import "time"

// Issue represents a GitHub Issue.
type Issue struct {
	// Owner is the repository owner (user or organization).
	Owner string
	// Repo is the repository name.
	Repo string
	// Number is the issue number.
	Number int
	// Title is the issue title.
	Title string
	// Body is the issue description content.
	Body string
	// State is the issue state ("open" or "closed").
	State string
	// Author is the username of the issue author.
	Author string
	// CreatedAt is the timestamp when the issue was created.
	CreatedAt time.Time
	// UpdatedAt is the timestamp when the issue was last updated.
	UpdatedAt time.Time
	// Labels are the labels associated with the issue.
	Labels []Label
	// Milestone is the milestone associated with the issue, if any.
	Milestone *Milestone
	// Reactions are the reactions on the issue.
	Reactions Reactions
	// Comments are the issue comments in chronological order.
	Comments []Comment
}

// PullRequest represents a GitHub Pull Request.
// It embeds Issue and adds PR-specific fields.
type PullRequest struct {
	Issue
	// Mergeable indicates whether the PR can be merged.
	Mergeable bool
	// Merged indicates whether the PR has been merged.
	Merged bool
	// MergedAt is the timestamp when the PR was merged, if applicable.
	MergedAt *time.Time
	// Additions is the number of lines added in the PR.
	Additions int
	// Deletions is the number of lines deleted in the PR.
	Deletions int
	// ChangedFiles is the number of files changed in the PR.
	ChangedFiles int
	// HeadBranch is the name of the source branch.
	HeadBranch string
	// BaseBranch is the name of the target branch.
	BaseBranch string
	// Files are the files changed in the PR.
	Files []File
}

// Discussion represents a GitHub Discussion.
type Discussion struct {
	// Number is the discussion number.
	Number int
	// Title is the discussion title.
	Title string
	// Body is the discussion content.
	Body string
	// Author is the username of the discussion author.
	Author string
	// CreatedAt is the timestamp when the discussion was created.
	CreatedAt time.Time
	// Category is the discussion category name.
	Category string
	// Upvotes is the number of upvotes the discussion received.
	Upvotes int
	// Replies are the discussion replies (possibly nested).
	Replies []DiscussionReply
}

// Label represents a GitHub label.
type Label struct {
	// Name is the label name.
	Name string
	// Color is the label color (hex code without #).
	Color string
}

// Milestone represents a GitHub milestone.
type Milestone struct {
	// Title is the milestone title.
	Title string
	// State is the milestone state ("open" or "closed").
	State string
}

// Reactions represents GitHub reactions on a resource.
type Reactions struct {
	// ThumbsUp is the count of +1 reactions.
	ThumbsUp int
	// ThumbsDown is the count of -1 reactions.
	ThumbsDown int
	// Laugh is the count of laugh reactions.
	Laugh int
	// Hooray is the count of hooray reactions.
	Hooray int
	// Confused is the count of confused reactions.
	Confused int
	// Heart is the count of heart reactions.
	Heart int
	// Rocket is the count of rocket reactions.
	Rocket int
	// Eyes is the count of eyes reactions.
	Eyes int
}

// Total returns the total count of all reactions.
func (r Reactions) Total() int {
	return r.ThumbsUp + r.ThumbsDown + r.Laugh + r.Hooray +
		r.Confused + r.Heart + r.Rocket + r.Eyes
}

// Comment represents a comment on an issue or pull request.
type Comment struct {
	// ID is the comment ID.
	ID int64
	// Author is the username of the comment author.
	Author string
	// Body is the comment content.
	Body string
	// CreatedAt is the timestamp when the comment was created.
	CreatedAt time.Time
	// Reactions are the reactions on the comment.
	Reactions Reactions
}

// File represents a file changed in a pull request.
type File struct {
	// Path is the file path.
	Path string
	// Additions is the number of lines added.
	Additions int
	// Deletions is the number of lines deleted.
	Deletions int
	// Patch is the unified diff for the file.
	Patch string
	// BlobURL is the GitHub URL for the file.
	BlobURL string
}

// DiscussionReply represents a reply in a GitHub Discussion.
type DiscussionReply struct {
	// ID is the GraphQL ID of the reply.
	ID string
	// Author is the username of the reply author.
	Author string
	// Body is the reply content.
	Body string
	// CreatedAt is the timestamp when the reply was created.
	CreatedAt time.Time
	// Replies are nested replies to this reply.
	Replies []DiscussionReply
}
