package parser

import (
	"testing"
)

// TestParseURL is a table-driven test for the ParseURL function.
// It covers various URL formats and error cases as specified in spec.md section 2.1.1.
func TestParseURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    *ParsedURL
		wantErr bool
	}{
		// Valid Issue URLs
		{
			name: "valid standard issue URL",
			url:  "https://github.com/owner/repo/issues/123",
			want: &ParsedURL{
				Owner:  "owner",
				Repo:   "repo",
				Number: 123,
				Type:   TypeIssue,
			},
			wantErr: false,
		},
		{
			name: "valid issue URL with .git suffix",
			url:  "https://github.com/owner/repo.git/issues/123",
			want: &ParsedURL{
				Owner:  "owner",
				Repo:   "repo",
				Number: 123,
				Type:   TypeIssue,
			},
			wantErr: false,
		},
		{
			name: "valid issue URL with www subdomain",
			url:  "https://www.github.com/owner/repo/issues/456",
			want: &ParsedURL{
				Owner:  "owner",
				Repo:   "repo",
				Number: 456,
				Type:   TypeIssue,
			},
			wantErr: false,
		},
		{
			name: "valid issue URL with http protocol",
			url:  "http://github.com/owner/repo/issues/789",
			want: &ParsedURL{
				Owner:  "owner",
				Repo:   "repo",
				Number: 789,
				Type:   TypeIssue,
			},
			wantErr: false,
		},
		{
			name: "valid issue URL with all variations combined",
			url:  "http://www.github.com/owner/repo.git/issues/999",
			want: &ParsedURL{
				Owner:  "owner",
				Repo:   "repo",
				Number: 999,
				Type:   TypeIssue,
			},
			wantErr: false,
		},

		// Valid Pull Request URLs
		{
			name: "valid standard PR URL",
			url:  "https://github.com/owner/repo/pull/123",
			want: &ParsedURL{
				Owner:  "owner",
				Repo:   "repo",
				Number: 123,
				Type:   TypePullRequest,
			},
			wantErr: false,
		},
		{
			name: "valid PR URL with .git suffix",
			url:  "https://github.com/owner/repo.git/pull/456",
			want: &ParsedURL{
				Owner:  "owner",
				Repo:   "repo",
				Number: 456,
				Type:   TypePullRequest,
			},
			wantErr: false,
		},
		{
			name: "valid PR URL with www subdomain",
			url:  "https://www.github.com/owner/repo/pull/789",
			want: &ParsedURL{
				Owner:  "owner",
				Repo:   "repo",
				Number: 789,
				Type:   TypePullRequest,
			},
			wantErr: false,
		},
		{
			name: "valid PR URL with http protocol",
			url:  "http://github.com/owner/repo/pull/101",
			want: &ParsedURL{
				Owner:  "owner",
				Repo:   "repo",
				Number: 101,
				Type:   TypePullRequest,
			},
			wantErr: false,
		},

		// Valid Discussion URLs
		{
			name: "valid standard discussion URL",
			url:  "https://github.com/owner/repo/discussions/123",
			want: &ParsedURL{
				Owner:  "owner",
				Repo:   "repo",
				Number: 123,
				Type:   TypeDiscussion,
			},
			wantErr: false,
		},
		{
			name: "valid discussion URL with .git suffix",
			url:  "https://github.com/owner/repo.git/discussions/456",
			want: &ParsedURL{
				Owner:  "owner",
				Repo:   "repo",
				Number: 456,
				Type:   TypeDiscussion,
			},
			wantErr: false,
		},
		{
			name: "valid discussion URL with http protocol",
			url:  "http://www.github.com/owner/repo/discussions/789",
			want: &ParsedURL{
				Owner:  "owner",
				Repo:   "repo",
				Number: 789,
				Type:   TypeDiscussion,
			},
			wantErr: false,
		},

		// Invalid URLs - format errors
		{
			name:    "invalid URL - not a URL at all",
			url:     "not-a-url",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid URL - empty string",
			url:     "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid URL - missing scheme",
			url:     "github.com/owner/repo/issues/123",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid URL - wrong scheme",
			url:     "ftp://github.com/owner/repo/issues/123",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid URL - missing path",
			url:     "https://github.com",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid URL - only owner, no repo",
			url:     "https://github.com/owner",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid URL - owner/repo but no resource path",
			url:     "https://github.com/owner/repo",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid URL - missing resource number",
			url:     "https://github.com/owner/repo/issues",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid URL - resource number is not an integer",
			url:     "https://github.com/owner/repo/issues/abc",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid URL - negative number",
			url:     "https://github.com/owner/repo/issues/-1",
			want:    nil,
			wantErr: true,
		},

		// Unsupported URL types
		{
			name:    "unsupported URL - repository home page",
			url:     "https://github.com/owner/repo",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "unsupported URL - actions path",
			url:     "https://github.com/owner/repo/actions",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "unsupported URL - wiki path",
			url:     "https://github.com/owner/repo/wiki",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "unsupported URL - releases path",
			url:     "https://github.com/owner/repo/releases",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "unsupported URL - security path",
			url:     "https://github.com/owner/repo/security",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "unsupported URL - settings path",
			url:     "https://github.com/owner/repo/settings",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "unsupported URL - commits path",
			url:     "https://github.com/owner/repo/commits/main",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "unsupported URL - tree path",
			url:     "https://github.com/owner/repo/tree/main",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "unsupported URL - blob path",
			url:     "https://github.com/owner/repo/blob/main/README.md",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "unsupported URL - different domain",
			url:     "https://gitlab.com/owner/repo/issues/123",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Errorf("ParseURL() returned nil, expected non-nil result")
					return
				}
				if got.Owner != tt.want.Owner {
					t.Errorf("ParseURL().Owner = %v, want %v", got.Owner, tt.want.Owner)
				}
				if got.Repo != tt.want.Repo {
					t.Errorf("ParseURL().Repo = %v, want %v", got.Repo, tt.want.Repo)
				}
				if got.Number != tt.want.Number {
					t.Errorf("ParseURL().Number = %v, want %v", got.Number, tt.want.Number)
				}
				if got.Type != tt.want.Type {
					t.Errorf("ParseURL().Type = %v, want %v", got.Type, tt.want.Type)
				}
			}
		})
	}
}
