package config

import (
	"testing"
)

func TestValidate_MissingSourceType(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{Token: "tok"},
		Destination: DestinationConfig{
			Type: "file",
			Path: "/tmp/out.json",
		},
	}
	err := cfg.Validate()
	if err == nil || err.Error() != "source.type is required" {
		t.Errorf("expected 'source.type is required', got %v", err)
	}
}

func TestValidate_InvalidSourceType(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{Type: "unknown", Token: "tok", Repo: "o/r"},
		Destination: DestinationConfig{
			Type: "file",
			Path: "/tmp/out.json",
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
			Type: "file",
			Path: "/tmp/out.json",
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
			Type: "file",
			Path: "/tmp/out.json",
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
			Type: "unknown",
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
			Type: "file",
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
			Type: "file",
			Path: "/tmp/out.json",
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
			Type:   "file",
			Path:   "/tmp/out.json",
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
			Type:   "file",
			Path:   "/tmp/out.yaml",
			Format: "yaml",
		},
	}
	err := cfg.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_ValidVaultSource(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{
			Type:         "vault",
			VaultAddress: "https://vault.example.com",
			VaultPath:    "myapp/config",
			Token:        "hvs.xxx",
		},
		Destination: DestinationConfig{
			Type: "file",
			Path: "/tmp/out.json",
		},
	}
	err := cfg.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_VaultSourceMissingAddress(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{
			Type:      "vault",
			VaultPath: "myapp/config",
			Token:     "hvs.xxx",
		},
		Destination: DestinationConfig{
			Type: "file",
			Path: "/tmp/out.json",
		},
	}
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for missing vault address")
	}
}

func TestValidate_VaultSourceMissingPath(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{
			Type:         "vault",
			VaultAddress: "https://vault.example.com",
			Token:        "hvs.xxx",
		},
		Destination: DestinationConfig{
			Type: "file",
			Path: "/tmp/out.json",
		},
	}
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for missing vault path")
	}
}

func TestValidate_ValidEtcdSource(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{
			Type:          "etcd",
			EtcdEndpoints: []string{"http://localhost:2379"},
		},
		Destination: DestinationConfig{
			Type: "file",
			Path: "/tmp/out.json",
		},
	}
	err := cfg.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_EtcdSourceMissingEndpoints(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{
			Type: "etcd",
		},
		Destination: DestinationConfig{
			Type: "file",
			Path: "/tmp/out.json",
		},
	}
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for missing etcd endpoints")
	}
}

func TestValidate_ValidKubernetesSource(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{
			Type:          "kubernetes",
			KubeNamespace: "default",
		},
		Destination: DestinationConfig{
			Type: "file",
			Path: "/tmp/out.json",
		},
	}
	err := cfg.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_KubernetesSourceMissingNamespace(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{
			Type: "kubernetes",
		},
		Destination: DestinationConfig{
			Type: "file",
			Path: "/tmp/out.json",
		},
	}
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for missing kube namespace")
	}
}

func TestValidate_ValidVaultDest(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{Type: "github", Repo: "o/r", Token: "tok"},
		Destination: DestinationConfig{
			Type:         "vault",
			VaultAddress: "https://vault.example.com",
			VaultPath:    "myapp/config",
			Token:        "hvs.xxx",
		},
	}
	err := cfg.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_ValidEtcdDest(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{Type: "github", Repo: "o/r", Token: "tok"},
		Destination: DestinationConfig{
			Type:          "etcd",
			EtcdEndpoints: []string{"http://localhost:2379"},
		},
	}
	err := cfg.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_ValidKubernetesDest(t *testing.T) {
	cfg := &Config{
		Source: SourceConfig{Type: "github", Repo: "o/r", Token: "tok"},
		Destination: DestinationConfig{
			Type:          "kubernetes",
			KubeNamespace: "default",
		},
	}
	err := cfg.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
