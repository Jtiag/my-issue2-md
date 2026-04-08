// Package cli provides the command-line interface for issue2md.
package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TestExecute is a table-driven test for the Execute function.
// It covers various input scenarios as specified in spec.md.
func TestExecute(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantExit     int
		wantInStdout string
		wantInStderr string
		setupEnv     func() func() // Setup environment, returns cleanup function
	}{
		{
			name:         "help flag outputs help and returns 0",
			args:         []string{"-h"},
			wantExit:     0,
			wantInStdout: "issue2md",
			wantInStderr: "",
		},
		{
			name:         "help flag -help outputs help and returns 0",
			args:         []string{"-help"},
			wantExit:     0,
			wantInStdout: "issue2md",
		},
		{
			name:         "version flag outputs version and returns 0",
			args:         []string{"-v"},
			wantExit:     0,
			wantInStdout: "version",
		},
		{
			name:         "version flag -version outputs version and returns 0",
			args:         []string{"-version"},
			wantExit:     0,
			wantInStdout: "version",
		},
		{
			name:         "no arguments returns error code 1",
			args:         []string{},
			wantExit:     1,
			wantInStderr: "required",
		},
		{
			name:         "invalid URL returns error code 1",
			args:         []string{"not-a-url"},
			wantExit:     1,
			wantInStderr: "invalid",
		},
		{
			name:     "missing GITHUB_TOKEN returns error code 4",
			args:     []string{"https://github.com/owner/repo/issues/123"},
			wantExit: 4,
			setupEnv: func() func() {
				// Save and unset token
				original := os.Getenv("GITHUB_TOKEN")
				os.Unsetenv("GITHUB_TOKEN")
				return func() {
					if original != "" {
						os.Setenv("GITHUB_TOKEN", original)
					}
				}
			},
			wantInStderr: "GITHUB_TOKEN",
		},
		{
			name:     "successful conversion returns 0",
			args:     []string{"https://github.com/golang/go/issues/40655"},
			wantExit: 0,
			setupEnv: func() func() {
				// This test will be skipped if no token is set
				return func() {}
			},
			wantInStdout: "#",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment if needed
			if tt.setupEnv != nil {
				cleanup := tt.setupEnv()
				defer cleanup()
			}

			// Skip token-dependent tests if token is not set
			if strings.Contains(tt.name, "GITHUB_TOKEN") && os.Getenv("GITHUB_TOKEN") == "" {
				t.Skip("GITHUB_TOKEN not set, skipping token-dependent test")
			}

			// Capture stdout and stderr
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			stdin := &bytes.Buffer{}

			// Execute
			exitCode := Execute(stdin, stdout, stderr, tt.args)

			// Check exit code
			if exitCode != tt.wantExit {
				t.Errorf("Execute() exitCode = %d, want %d", exitCode, tt.wantExit)
			}

			// Check stdout content
			if tt.wantInStdout != "" {
				stdoutStr := stdout.String()
				if !strings.Contains(strings.ToLower(stdoutStr), strings.ToLower(tt.wantInStdout)) {
					t.Errorf("Execute() stdout should contain %q, got %q", tt.wantInStdout, stdoutStr)
				}
			}

			// Check stderr content
			if tt.wantInStderr != "" {
				stderrStr := stderr.String()
				if !strings.Contains(strings.ToLower(stderrStr), strings.ToLower(tt.wantInStderr)) {
					t.Errorf("Execute() stderr should contain %q, got %q", tt.wantInStderr, stderrStr)
				}
			}
		})
	}
}

// TestRun tests the Run function.
func TestRun(t *testing.T) {
	// This is a simple test to ensure Run calls Execute with correct args
	// Full integration testing would require subprocess testing

	// Test with help flag (should exit 0)
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"issue2md", "-h"}
	// We can't actually test Run() since it calls os.Exit()
	// But we can test that Execute is called correctly
}
