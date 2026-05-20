package config

import (
	"testing"
)

func TestValidate_MissingSourceType(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{Token: "tok"},
		Destination: DestinationConfig{
			Type:  "file",
			Path:  "/tmp/out.json",
			Token: "tok",
		},
	}
	err := cfg.Validate()
	if err == nil || err.Error() != "source.type is required" {
		t.Errorf("expected 'source.type is required', got %v", err)
	}
}

func TestValidate_InvalidSourceType(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{Type: "vault", Token: "tok", Repo: "o/r"},
		Destination: DestinationConfig{
			Type:  "file",
			Path:  "/tmp/out.json",
			Token: "tok",
		},
	}
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for invalid source type")
	}
}

func TestValidate_MissingRepo(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{Type: "github", Token: "tok"},
		Destination: DestinationConfig{
			Type:  "file",
			Path:  "/tmp/out.json",
			Token: "tok",
		},
	}
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for missing repo")
	}
}

func TestValidate_MissingToken(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{Type: "github", Repo: "o/r"},
		Destination: DestinationConfig{
			Type:  "file",
			Path:  "/tmp/out.json",
			Token: "tok",
		},
	}
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for missing token")
	}
}

func TestValidate_InvalidDestType(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{Type: "github", Repo: "o/r", Token: "tok"},
		Destination: DestinationConfig{
			Type:  "vault",
			Token: "tok",
		},
	}
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for invalid dest type")
	}
}

func TestValidate_FileDestMissingPath(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{Type: "github", Repo: "o/r", Token: "tok"},
		Destination: DestinationConfig{
			Type:  "file",
			Token: "tok",
		},
	}
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for missing file path")
	}
}

func TestValidate_InvalidStrategy(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{Type: "github", Repo: "o/r", Token: "tok"},
		Destination: DestinationConfig{
			Type:             "file",
			Path:             "/tmp/out.json",
			Token:            "tok",
			ConflictStrategy: "invalid",
		},
	}
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for invalid strategy")
	}
}

func TestValidate_DefaultStrategy(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{Type: "github", Repo: "o/r", Token: "tok"},
		Destination: DestinationConfig{
			Type:  "file",
			Path:  "/tmp/out.json",
			Token: "tok",
		},
	}
	err := cfg.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if cfg.Destination.ConflictStrategy != "replace" {
		t.Errorf("expected default strategy 'replace', got %s", cfg.Destination.ConflictStrategy)
	}
}

func TestValidate_ValidGithubToFile(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{Type: "github", Repo: "owner/repo", Token: "ghp_xxx"},
		Destination: DestinationConfig{
			Type:  "file",
			Path:  "/tmp/out.json",
			Format: "json",
		},
	}
	err := cfg.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_ValidGitlabToFile(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{Type: "gitlab", ProjectID: "123", Token: "glpat-xxx"},
		Destination: DestinationConfig{
			Type:  "file",
			Path:  "/tmp/out.yaml",
			Format: "yaml",
		},
	}
	err := cfg.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
