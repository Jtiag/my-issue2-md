// Package config provides configuration management for the CLI.
package config

import (
	"os"
	"strings"
	"testing"
)

// TestParseFlags is a table-driven test for the ParseFlags function.
// It covers all flag combinations and positional arguments as specified in spec.md.
func TestParseFlags(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantConfig  *Config
		wantErr     bool
		wantErrCode int // For exit code validation
	}{
		{
			name:    "empty arguments",
			args:    []string{},
			wantErr: true,
		},
		{
			name: "help flag -h",
			args: []string{"-h"},
			wantConfig: &Config{
				helpRequested: true,
			},
			wantErr: false,
		},
		{
			name: "help flag -help",
			args: []string{"-help"},
			wantConfig: &Config{
				helpRequested: true,
			},
			wantErr: false,
		},
		{
			name: "version flag -v",
			args: []string{"-v"},
			wantConfig: &Config{
				versionRequested: true,
			},
			wantErr: false,
		},
		{
			name: "version flag -version",
			args: []string{"-version"},
			wantConfig: &Config{
				versionRequested: true,
			},
			wantErr: false,
		},
		{
			name: "enable reactions flag",
			args: []string{"-enable-reactions", "https://github.com/owner/repo/issues/123"},
			wantConfig: &Config{
				EnableReactions: true,
				OutputFile:      "",
			},
			wantErr: false,
		},
		{
			name: "enable user links flag",
			args: []string{"-enable-user-links", "https://github.com/owner/repo/issues/123"},
			wantConfig: &Config{
				EnableUserLinks: true,
				OutputFile:      "",
			},
			wantErr: false,
		},
		{
			name: "both enable flags",
			args: []string{"-enable-reactions", "-enable-user-links", "https://github.com/owner/repo/issues/123"},
			wantConfig: &Config{
				EnableReactions: true,
				EnableUserLinks: true,
				OutputFile:      "",
			},
			wantErr: false,
		},
		{
			name: "with URL only",
			args: []string{"https://github.com/owner/repo/issues/123"},
			wantConfig: &Config{
				EnableReactions: false,
				EnableUserLinks: false,
				OutputFile:      "",
			},
			wantErr: false,
		},
		{
			name: "with URL and output file",
			args: []string{"https://github.com/owner/repo/issues/123", "output.md"},
			wantConfig: &Config{
				EnableReactions: false,
				EnableUserLinks: false,
				OutputFile:      "output.md",
			},
			wantErr: false,
		},
		{
			name: "with all flags and URL and output",
			args: []string{"-enable-reactions", "-enable-user-links", "https://github.com/owner/repo/issues/123", "output.md"},
			wantConfig: &Config{
				EnableReactions: true,
				EnableUserLinks: true,
				OutputFile:      "output.md",
			},
			wantErr: false,
		},
		{
			name: "PR URL",
			args: []string{"https://github.com/owner/repo/pull/456"},
			wantConfig: &Config{
				EnableReactions: false,
				EnableUserLinks: false,
				OutputFile:      "",
			},
			wantErr: false,
		},
		{
			name: "discussion URL",
			args: []string{"https://github.com/owner/repo/discussions/789"},
			wantConfig: &Config{
				EnableReactions: false,
				EnableUserLinks: false,
				OutputFile:      "",
			},
			wantErr: false,
		},
		{
			name:    "missing URL with enable flag",
			args:    []string{"-enable-reactions"},
			wantErr: true,
		},
		{
			name:    "invalid flag",
			args:    []string{"-invalid-flag"},
			wantErr: true,
		},
		{
			name:    "too many positional args",
			args:    []string{"https://github.com/owner/repo/issues/123", "output.md", "extra"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFlags(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Error("ParseFlags() returned nil config, expected non-nil")
					return
				}

				if got.EnableReactions != tt.wantConfig.EnableReactions {
					t.Errorf("ParseFlags().EnableReactions = %v, want %v",
						got.EnableReactions, tt.wantConfig.EnableReactions)
				}

				if got.EnableUserLinks != tt.wantConfig.EnableUserLinks {
					t.Errorf("ParseFlags().EnableUserLinks = %v, want %v",
						got.EnableUserLinks, tt.wantConfig.EnableUserLinks)
				}

				if got.OutputFile != tt.wantConfig.OutputFile {
					t.Errorf("ParseFlags().OutputFile = %v, want %v",
						got.OutputFile, tt.wantConfig.OutputFile)
				}

				if got.helpRequested != tt.wantConfig.helpRequested {
					t.Errorf("ParseFlags().helpRequested = %v, want %v",
						got.helpRequested, tt.wantConfig.helpRequested)
				}

				if got.versionRequested != tt.wantConfig.versionRequested {
					t.Errorf("ParseFlags().versionRequested = %v, want %v",
						got.versionRequested, tt.wantConfig.versionRequested)
				}
			}
		})
	}
}

// TestValidate tests the Validate method.
func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "help requested - no error",
			config: &Config{
				helpRequested: true,
			},
			wantErr: false,
		},
		{
			name: "version requested - no error",
			config: &Config{
				versionRequested: true,
			},
			wantErr: false,
		},
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config == nil {
				// Can't call Validate on nil
				t.Log("nil config test - skipping method call")
				return
			}
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestOutputOptions tests the OutputOptions method.
func TestOutputOptions(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   Options
	}{
		{
			name: "both options enabled",
			config: &Config{
				EnableReactions: true,
				EnableUserLinks: true,
			},
			want: Options{
				EnableReactions: true,
				EnableUserLinks: true,
			},
		},
		{
			name: "only reactions enabled",
			config: &Config{
				EnableReactions: true,
				EnableUserLinks: false,
			},
			want: Options{
				EnableReactions: true,
				EnableUserLinks: false,
			},
		},
		{
			name: "only user links enabled",
			config: &Config{
				EnableReactions: false,
				EnableUserLinks: true,
			},
			want: Options{
				EnableReactions: false,
				EnableUserLinks: true,
			},
		},
		{
			name: "no options enabled",
			config: &Config{
				EnableReactions: false,
				EnableUserLinks: false,
			},
			want: Options{
				EnableReactions: false,
				EnableUserLinks: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.OutputOptions()
			if got.EnableReactions != tt.want.EnableReactions {
				t.Errorf("Config.OutputOptions().EnableReactions = %v, want %v",
					got.EnableReactions, tt.want.EnableReactions)
			}
			if got.EnableUserLinks != tt.want.EnableUserLinks {
				t.Errorf("Config.OutputOptions().EnableUserLinks = %v, want %v",
					got.EnableUserLinks, tt.want.EnableUserLinks)
			}
		})
	}
}

// TestVersionInfo tests the VersionInfo function.
func TestVersionInfo(t *testing.T) {
	version := VersionInfo()
	if version == "" {
		t.Error("VersionInfo() returned empty string")
	}
	if !strings.Contains(version, "issue2md") {
		t.Error("VersionInfo() should contain 'issue2md'")
	}
}

// TestHelpOutput tests that help output is correctly formatted.
func TestHelpOutput(t *testing.T) {
	help := HelpText()
	if help == "" {
		t.Error("HelpText() returned empty string")
	}

	// Check for expected sections
	expectedContent := []string{
		"issue2md",
		"usage",
		"flags",
		"-enable-reactions",
		"-enable-user-links",
		"-h",
		"-help",
		"-v",
		"-version",
	}

	for _, content := range expectedContent {
		if !strings.Contains(strings.ToLower(help), content) {
			t.Errorf("HelpText() should contain %q", content)
		}
	}
}

// TestEnvToken tests that GITHUB_TOKEN can be read from environment.
func TestEnvToken(t *testing.T) {
	// Save original value
	original := os.Getenv("GITHUB_TOKEN")
	defer os.Setenv("GITHUB_TOKEN", original)

	// Test with token set
	os.Setenv("GITHUB_TOKEN", "test-token")
	token := os.Getenv("GITHUB_TOKEN")
	if token != "test-token" {
		t.Errorf("Expected GITHUB_TOKEN to be 'test-token', got %q", token)
	}

	// Test with token unset
	os.Unsetenv("GITHUB_TOKEN")
	token = os.Getenv("GITHUB_TOKEN")
	if token != "" {
		t.Errorf("Expected GITHUB_TOKEN to be empty, got %q", token)
	}
}
