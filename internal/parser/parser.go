// Package parser provides GitHub URL parsing functionality.
//
// It parses GitHub repository URLs and extracts information about
// issues, pull requests, and discussions.
package parser

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ResourceType represents the type of GitHub resource.
type ResourceType string

const (
	// TypeIssue represents a GitHub Issue.
	TypeIssue ResourceType = "issue"
	// TypePullRequest represents a GitHub Pull Request.
	TypePullRequest ResourceType = "pull_request"
	// TypeDiscussion represents a GitHub Discussion.
	TypeDiscussion ResourceType = "discussion"
)

// ParsedURL represents a parsed GitHub URL.
type ParsedURL struct {
	// Owner is the repository owner (user or organization).
	Owner string
	// Repo is the repository name.
	Repo string
	// Number is the issue/PR/discussion number.
	Number int
	// Type is the resource type.
	Type ResourceType
}

// pathSegments represents validated and split URL path segments.
type pathSegments struct {
	// owner is the repository owner.
	owner string
	// repo is the repository name (without .git suffix).
	repo string
	// resourceType is the type of resource (issues, pull, discussions).
	resourceType string
	// number is the resource number as a string.
	number string
}

// ParseURL parses a GitHub URL and extracts resource information.
//
// Supported formats:
//   - https://github.com/owner/repo/issues/123
//   - https://github.com/owner/repo/pull/456
//   - https://github.com/owner/repo/discussions/789
//   - Variations with .git suffix, www subdomain, http protocol, etc.
//
// Returns ParsedURL or an error if the URL format is invalid or unsupported.
func ParseURL(rawURL string) (*ParsedURL, error) {
	// Validate URL and extract host/path
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// Validate and normalize host
	if err := validateHost(u.Host); err != nil {
		return nil, err
	}

	// Validate scheme
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("unsupported URL scheme: only http and https are supported")
	}

	// Split and validate path segments
	segments, err := splitAndValidatePath(rawURL, u.Path)
	if err != nil {
		return nil, err
	}

	// Determine resource type
	resourceType, err := parseResourceType(segments.resourceType)
	if err != nil {
		return nil, err
	}

	// Parse and validate number
	number, err := parseNumber(segments.number)
	if err != nil {
		return nil, err
	}

	return &ParsedURL{
		Owner:  segments.owner,
		Repo:   segments.repo,
		Number: number,
		Type:   resourceType,
	}, nil
}

// validateHost validates that the host is github.com or www.github.com.
func validateHost(host string) error {
	if host == "" {
		return fmt.Errorf("invalid URL: missing host")
	}

	// Remove port if present
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}

	if host != "github.com" && host != "www.github.com" {
		return fmt.Errorf("unsupported URL: not a GitHub URL")
	}

	return nil
}

// splitAndValidatePath splits the URL path and validates its structure.
func splitAndValidatePath(rawURL, path string) (*pathSegments, error) {
	if path == "" {
		return nil, fmt.Errorf("invalid URL: missing path")
	}

	// Remove leading slash and split
	path = strings.TrimPrefix(path, "/")
	segments := strings.Split(path, "/")

	if len(segments) < 4 {
		return nil, fmt.Errorf("invalid URL: expected owner/repo/resource/number format")
	}

	// Extract and validate owner
	owner := segments[0]
	if owner == "" {
		return nil, fmt.Errorf("invalid URL: missing owner")
	}

	// Extract and validate repo name
	repo := segments[1]
	if repo == "" {
		return nil, fmt.Errorf("invalid URL: missing repo")
	}
	repo = strings.TrimSuffix(repo, ".git")
	if repo == "" {
		return nil, fmt.Errorf("invalid URL: repo name is empty after removing .git suffix")
	}

	// Extract resource type and number
	resourceType := segments[2]
	if resourceType == "" {
		return nil, fmt.Errorf("invalid URL: missing resource type")
	}

	number := segments[3]
	if number == "" {
		return nil, fmt.Errorf("invalid URL: missing number")
	}

	// Ensure no extra segments
	if len(segments) > 4 {
		return nil, fmt.Errorf("invalid URL: too many path segments")
	}

	return &pathSegments{
		owner:        owner,
		repo:         repo,
		resourceType: resourceType,
		number:       number,
	}, nil
}

// parseResourceType converts a resource type string to ResourceType.
func parseResourceType(resourceType string) (ResourceType, error) {
	switch resourceType {
	case "issues":
		return TypeIssue, nil
	case "pull":
		return TypePullRequest, nil
	case "discussions":
		return TypeDiscussion, nil
	default:
		return "", fmt.Errorf("unsupported resource type %q: must be issues, pull, or discussions", resourceType)
	}
}

// parseNumber converts a string to a positive integer.
func parseNumber(numberStr string) (int, error) {
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return 0, fmt.Errorf("invalid number %q: %w", numberStr, err)
	}

	if number <= 0 {
		return 0, fmt.Errorf("number must be positive, got %d", number)
	}

	return number, nil
}
