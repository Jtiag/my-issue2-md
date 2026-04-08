// Package converter provides GitHub data to Markdown conversion functionality.
package converter

import (
	"strings"
	"testing"
	"time"

	"github.com/bigwhite/my-issue2md/internal/config"
	"github.com/bigwhite/my-issue2md/internal/github"
)

// TestDiscussionToMarkdown is a table-driven test for the DiscussionToMarkdown function.
// It covers various discussion conversion scenarios as specified in spec.md section 2.4.
func TestDiscussionToMarkdown(t *testing.T) {
	fixedTime := time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name    string
		disc    *github.Discussion
		opts    config.Options
		wantErr bool
	}{
		{
			name: "basic discussion conversion",
			disc: &github.Discussion{
				Number:    123,
				Title:     "Discussion Title",
				Body:      "This is a discussion",
				Author:    "discuser",
				CreatedAt: fixedTime,
				Category:  "General",
				Upvotes:   5,
				Replies:   []github.DiscussionReply{},
			},
			opts:    config.Options{},
			wantErr: false,
		},
		{
			name: "discussion with single reply",
			disc: &github.Discussion{
				Number:    456,
				Title:     "Question",
				Body:      "How do I do this?",
				Author:    "questionuser",
				CreatedAt: fixedTime,
				Category:  "Q&A",
				Upvotes:   2,
				Replies: []github.DiscussionReply{
					{
						ID:        "reply_1",
						Author:    "answeruser",
						Body:      "Here's the answer",
						CreatedAt: fixedTime.Add(time.Hour),
						Replies:   []github.DiscussionReply{},
					},
				},
			},
			opts:    config.Options{},
			wantErr: false,
		},
		{
			name: "discussion with nested replies",
			disc: &github.Discussion{
				Number:    789,
				Title:     "Nested Discussion",
				Body:      "Main topic",
				Author:    "mainuser",
				CreatedAt: fixedTime,
				Category:  "General",
				Upvotes:   10,
				Replies: []github.DiscussionReply{
					{
						ID:        "reply_1",
						Author:    "replyuser1",
						Body:      "First reply",
						CreatedAt: fixedTime.Add(time.Hour),
						Replies: []github.DiscussionReply{
							{
								ID:        "reply_2",
								Author:    "replyuser2",
								Body:      "Nested reply",
								CreatedAt: fixedTime.Add(2 * time.Hour),
								Replies:   []github.DiscussionReply{},
							},
						},
					},
				},
			},
			opts:    config.Options{},
			wantErr: false,
		},
		{
			name: "discussion with multiple top-level replies",
			disc: &github.Discussion{
				Number:    101,
				Title:     "Multi-reply Discussion",
				Body:      "Topic with many replies",
				Author:    "topicuser",
				CreatedAt: fixedTime,
				Category:  "Ideas",
				Upvotes:   15,
				Replies: []github.DiscussionReply{
					{
						ID:        "reply_1",
						Author:    "user1",
						Body:      "First response",
						CreatedAt: fixedTime.Add(time.Hour),
						Replies:   []github.DiscussionReply{},
					},
					{
						ID:        "reply_2",
						Author:    "user2",
						Body:      "Second response",
						CreatedAt: fixedTime.Add(2 * time.Hour),
						Replies:   []github.DiscussionReply{},
					},
					{
						ID:        "reply_3",
						Author:    "user3",
						Body:      "Third response",
						CreatedAt: fixedTime.Add(3 * time.Hour),
						Replies:   []github.DiscussionReply{},
					},
				},
			},
			opts:    config.Options{},
			wantErr: false,
		},
		{
			name: "discussion with deeply nested replies",
			disc: &github.Discussion{
				Number:    202,
				Title:     "Deep Thread",
				Body:      "Starting a deep conversation",
				Author:    "starter",
				CreatedAt: fixedTime,
				Category:  "General",
				Upvotes:   3,
				Replies: []github.DiscussionReply{
					{
						ID:        "r1",
						Author:    "u1",
						Body:      "Level 1",
						CreatedAt: fixedTime.Add(time.Hour),
						Replies: []github.DiscussionReply{
							{
								ID:        "r2",
								Author:    "u2",
								Body:      "Level 2",
								CreatedAt: fixedTime.Add(2 * time.Hour),
								Replies: []github.DiscussionReply{
									{
										ID:        "r3",
										Author:    "u3",
										Body:      "Level 3",
										CreatedAt: fixedTime.Add(3 * time.Hour),
										Replies:   []github.DiscussionReply{},
									},
								},
							},
						},
					},
				},
			},
			opts:    config.Options{},
			wantErr: false,
		},
		{
			name: "discussion with user links enabled",
			disc: &github.Discussion{
				Number:    303,
				Title:     "Discussion with Mentions",
				Body:      "@team please review",
				Author:    "mentionuser",
				CreatedAt: fixedTime,
				Category:  "Announcements",
				Upvotes:   8,
				Replies: []github.DiscussionReply{
					{
						ID:        "reply_1",
						Author:    "responder",
						Body:      "@mentionuser I'll review",
						CreatedAt: fixedTime.Add(time.Hour),
						Replies:   []github.DiscussionReply{},
					},
				},
			},
			opts:    config.Options{EnableUserLinks: true},
			wantErr: false,
		},
		{
			name:    "nil discussion",
			disc:    nil,
			opts:    config.Options{},
			wantErr: true,
		},
		{
			name: "discussion with empty body",
			disc: &github.Discussion{
				Number:    404,
				Title:     "Empty Body Discussion",
				Body:      "",
				Author:    "emptyuser",
				CreatedAt: fixedTime,
				Category:  "General",
				Upvotes:   0,
				Replies:   []github.DiscussionReply{},
			},
			opts:    config.Options{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DiscussionToMarkdown(tt.disc, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("DiscussionToMarkdown() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.disc != nil {
				// Check YAML front matter
				if !strings.Contains(got, "---") {
					t.Error("DiscussionToMarkdown() output missing YAML front matter")
				}

				// Check title
				if !strings.Contains(got, tt.disc.Title) {
					t.Errorf("DiscussionToMarkdown() output missing title %q", tt.disc.Title)
				}

				// Check author
				if !strings.Contains(got, tt.disc.Author) {
					t.Errorf("DiscussionToMarkdown() output missing author %q", tt.disc.Author)
				}

				// Check category
				if !strings.Contains(got, tt.disc.Category) {
					t.Errorf("DiscussionToMarkdown() output missing category %q", tt.disc.Category)
				}

				// Check upvotes
				if !strings.Contains(got, "upvote") {
					t.Error("DiscussionToMarkdown() output missing upvotes")
				}

				// Check replies section if replies exist
				if len(tt.disc.Replies) > 0 {
					if !strings.Contains(got, "Replies") {
						t.Error("DiscussionToMarkdown() output missing Replies section")
					}
				}

				// Check nested reply indentation
				if hasNestedReplies(tt.disc.Replies) {
					// Should have blockquotes for nested replies
					if !strings.Contains(got, ">") {
						t.Error("DiscussionToMarkdown() output missing blockquote indentation for nested replies")
					}
				}

				// Check user links if enabled
				if tt.opts.EnableUserLinks && (strings.Contains(tt.disc.Body, "@") || hasMentionsInReplies(tt.disc.Replies)) {
					if !strings.Contains(got, "](https://github.com/") {
						t.Error("DiscussionToMarkdown() output missing user links when enabled")
					}
				}
			}
		})
	}
}

// hasNestedReplies checks if there are any nested replies.
func hasNestedReplies(replies []github.DiscussionReply) bool {
	for _, r := range replies {
		if len(r.Replies) > 0 {
			return true
		}
	}
	return false
}

// hasMentionsInReplies checks if any reply contains @mentions.
func hasMentionsInReplies(replies []github.DiscussionReply) bool {
	for _, r := range replies {
		if strings.Contains(r.Body, "@") {
			return true
		}
		if hasMentionsInReplies(r.Replies) {
			return true
		}
	}
	return false
}
