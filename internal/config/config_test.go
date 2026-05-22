package config

import (
	"testing"

	"github.com/PapaDanielVi/secret-shift/internal/provider"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestConfig(t *testing.T, cfg *Config) *Config {
	t.Helper()
	v := viper.New()
	_ = v
	return cfg
}

func TestValidate_MissingSourceType(t *testing.T) {
	cfg := newTestConfig(t, &Config{
		Source: SourceConfig{GitConfig: GitConfig{Token: "tok"}},
		Destination: DestinationConfig{
			Type: provider.File,
			Path: "/tmp/out.json",
		},
	})
	err := cfg.Validate()
	if err == nil || err.Error() != "source.type is required" {
		t.Errorf("expected 'source.type is required', got %v", err)
	}
}

func TestValidate_InvalidSourceType(t *testing.T) {
	cfg := newTestConfig(t, &Config{
		Source: SourceConfig{
			Type:      "unknown",
			GitConfig: GitConfig{Token: "tok", Repo: "o/r"},
		},
		Destination: DestinationConfig{
			Type: provider.File,
			Path: "/tmp/out.json",
		},
	})
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for invalid source type")
	}
}

func TestValidate_MissingRepo(t *testing.T) {
	cfg := newTestConfig(t, &Config{
		Source: SourceConfig{
			Type:      provider.GitHub,
			GitConfig: GitConfig{Token: "tok"},
		},
		Destination: DestinationConfig{
			Type: provider.File,
			Path: "/tmp/out.json",
		},
	})
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for missing repo")
	}
}

func TestValidate_MissingToken(t *testing.T) {
	cfg := newTestConfig(t, &Config{
		Source: SourceConfig{
			Type:      provider.GitHub,
			GitConfig: GitConfig{Repo: "o/r"},
		},
		Destination: DestinationConfig{
			Type: provider.File,
			Path: "/tmp/out.json",
		},
	})
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for missing token")
	}
}

func TestValidate_InvalidDestType(t *testing.T) {
	cfg := newTestConfig(t, &Config{
		Source: SourceConfig{
			Type:      provider.GitHub,
			GitConfig: GitConfig{Token: "tok", Repo: "o/r"},
		},
		Destination: DestinationConfig{
			Type: "unknown",
		},
	})
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for invalid dest type")
	}
}

func TestValidate_FileDestMissingPath(t *testing.T) {
	cfg := newTestConfig(t, &Config{
		Source: SourceConfig{
			Type:      provider.GitHub,
			GitConfig: GitConfig{Token: "tok", Repo: "o/r"},
		},
		Destination: DestinationConfig{
			Type: provider.File,
		},
	})
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for missing file path")
	}
}

func TestValidate_InvalidStrategy(t *testing.T) {
	cfg := newTestConfig(t, &Config{
		Source: SourceConfig{
			Type:      provider.GitHub,
			GitConfig: GitConfig{Token: "tok", Repo: "o/r"},
		},
		Destination: DestinationConfig{
			Type:             provider.File,
			Path:             "/tmp/out.json",
			ConflictStrategy: "invalid",
		},
	})
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for invalid strategy")
	}
}

func TestValidate_DefaultStrategy(t *testing.T) {
	cfg := newTestConfig(t, &Config{
		Source: SourceConfig{
			Type:      provider.GitHub,
			GitConfig: GitConfig{Token: "tok", Repo: "o/r"},
		},
		Destination: DestinationConfig{
			Type: provider.File,
			Path: "/tmp/out.json",
		},
	})
	err := cfg.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if cfg.Destination.ConflictStrategy != "replace" {
		t.Errorf("expected default strategy 'replace', got %s", cfg.Destination.ConflictStrategy)
	}
}

func TestValidate_ValidGithubToFile(t *testing.T) {
	cfg := newTestConfig(t, &Config{
		Source: SourceConfig{
			Type:      provider.GitHub,
			GitConfig: GitConfig{Repo: "owner/repo", Token: "ghp_xxx"},
		},
		Destination: DestinationConfig{
			Type:   provider.File,
			Path:   "/tmp/out.json",
			Format: "json",
		},
	})
	err := cfg.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_ValidGitlabToFile(t *testing.T) {
	cfg := newTestConfig(t, &Config{
		Source: SourceConfig{
			Type:      provider.GitLab,
			ProjectID: "123",
			GitConfig: GitConfig{Token: "glpat-xxx"},
		},
		Destination: DestinationConfig{
			Type:   provider.File,
			Path:   "/tmp/out.yaml",
			Format: "yaml",
		},
	})
	err := cfg.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_ValidVaultSource(t *testing.T) {
	cfg := newTestConfig(t, &Config{
		Source: SourceConfig{
			Type:        provider.Vault,
			VaultConfig: VaultConfig{VaultAddress: "https://vault.example.com", VaultPath: "myapp/config"},
			GitConfig:   GitConfig{Token: "hvs.xxx"},
		},
		Destination: DestinationConfig{
			Type: provider.File,
			Path: "/tmp/out.json",
		},
	})
	err := cfg.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_VaultSourceMissingAddress(t *testing.T) {
	cfg := newTestConfig(t, &Config{
		Source: SourceConfig{
			Type:        provider.Vault,
			VaultConfig: VaultConfig{VaultPath: "myapp/config"},
			GitConfig:   GitConfig{Token: "hvs.xxx"},
		},
		Destination: DestinationConfig{
			Type: provider.File,
			Path: "/tmp/out.json",
		},
	})
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for missing vault address")
	}
}

func TestValidate_VaultSourceMissingPath(t *testing.T) {
	cfg := newTestConfig(t, &Config{
		Source: SourceConfig{
			Type:        provider.Vault,
			VaultConfig: VaultConfig{VaultAddress: "https://vault.example.com"},
			GitConfig:   GitConfig{Token: "hvs.xxx"},
		},
		Destination: DestinationConfig{
			Type: provider.File,
			Path: "/tmp/out.json",
		},
	})
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for missing vault path")
	}
}

func TestValidate_ValidEtcdSource(t *testing.T) {
	cfg := newTestConfig(t, &Config{
		Source: SourceConfig{
			Type:       provider.Etcd,
			EtcdConfig: EtcdConfig{EtcdEndpoints: []string{"http://localhost:2379"}},
		},
		Destination: DestinationConfig{
			Type: provider.File,
			Path: "/tmp/out.json",
		},
	})
	err := cfg.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_EtcdSourceMissingEndpoints(t *testing.T) {
	cfg := newTestConfig(t, &Config{
		Source: SourceConfig{
			Type: provider.Etcd,
		},
		Destination: DestinationConfig{
			Type: provider.File,
			Path: "/tmp/out.json",
		},
	})
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for missing etcd endpoints")
	}
}

func TestValidate_ValidKubernetesSource(t *testing.T) {
	cfg := newTestConfig(t, &Config{
		Source: SourceConfig{
			Type:             provider.Kubernetes,
			KubernetesConfig: KubernetesConfig{KubeNamespace: "default"},
		},
		Destination: DestinationConfig{
			Type: provider.File,
			Path: "/tmp/out.json",
		},
	})
	err := cfg.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_KubernetesSourceMissingNamespace(t *testing.T) {
	cfg := newTestConfig(t, &Config{
		Source: SourceConfig{
			Type: provider.Kubernetes,
		},
		Destination: DestinationConfig{
			Type: provider.File,
			Path: "/tmp/out.json",
		},
	})
	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for missing kube namespace")
	}
}

func TestValidate_ValidVaultDest(t *testing.T) {
	cfg := newTestConfig(t, &Config{
		Source: SourceConfig{
			Type:      provider.GitHub,
			GitConfig: GitConfig{Repo: "o/r", Token: "tok"},
		},
		Destination: DestinationConfig{
			Type:        provider.Vault,
			VaultConfig: VaultConfig{VaultAddress: "https://vault.example.com", VaultPath: "myapp/config"},
			GitConfig:   GitConfig{Token: "hvs.xxx"},
		},
	})
	err := cfg.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_ValidEtcdDest(t *testing.T) {
	cfg := newTestConfig(t, &Config{
		Source: SourceConfig{
			Type:      provider.GitHub,
			GitConfig: GitConfig{Repo: "o/r", Token: "tok"},
		},
		Destination: DestinationConfig{
			Type:       provider.Etcd,
			EtcdConfig: EtcdConfig{EtcdEndpoints: []string{"http://localhost:2379"}},
		},
	})
	err := cfg.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_ValidKubernetesDest(t *testing.T) {
	cfg := newTestConfig(t, &Config{
		Source: SourceConfig{
			Type:      provider.GitHub,
			GitConfig: GitConfig{Repo: "o/r", Token: "tok"},
		},
		Destination: DestinationConfig{
			Type:             provider.Kubernetes,
			KubernetesConfig: KubernetesConfig{KubeNamespace: "default"},
		},
	})
	err := cfg.Validate()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *Config
		wantErr   bool
		errSubstr string
	}{
		{
			name: "valid github to file",
			cfg: &Config{
				Source: SourceConfig{
					Type:      provider.GitHub,
					GitConfig: GitConfig{Repo: "o/r", Token: "tok"},
				},
				Destination: DestinationConfig{
					Type: provider.File,
					Path: "/tmp/out.json",
				},
			},
			wantErr: false,
		},
		{
			name: "valid gitlab to vault",
			cfg: &Config{
				Source: SourceConfig{
					Type:      provider.GitLab,
					ProjectID: "123",
					GitConfig: GitConfig{Token: "tok"},
				},
				Destination: DestinationConfig{
					Type:        provider.Vault,
					VaultConfig: VaultConfig{VaultAddress: "https://vault.example.com", VaultPath: "app/config"},
					GitConfig:   GitConfig{Token: "hvs.xxx"},
				},
			},
			wantErr: false,
		},
		{
			name: "valid kubernetes to etcd",
			cfg: &Config{
				Source: SourceConfig{
					Type:             provider.Kubernetes,
					KubernetesConfig: KubernetesConfig{KubeNamespace: "prod"},
				},
				Destination: DestinationConfig{
					Type:       provider.Etcd,
					EtcdConfig: EtcdConfig{EtcdEndpoints: []string{"http://etcd:2379"}},
				},
			},
			wantErr: false,
		},
		{
			name: "empty source type",
			cfg: &Config{
				Source: SourceConfig{},
				Destination: DestinationConfig{
					Type: provider.File,
					Path: "/tmp/out.json",
				},
			},
			wantErr:   true,
			errSubstr: "source.type is required",
		},
		{
			name: "empty dest type",
			cfg: &Config{
				Source: SourceConfig{
					Type:      provider.GitHub,
					GitConfig: GitConfig{Repo: "o/r", Token: "tok"},
				},
				Destination: DestinationConfig{},
			},
			wantErr:   true,
			errSubstr: "destination.type is required",
		},
		{
			name: "invalid source type",
			cfg: &Config{
				Source: SourceConfig{Type: "invalid"},
				Destination: DestinationConfig{
					Type: provider.File,
					Path: "/tmp/out.json",
				},
			},
			wantErr:   true,
			errSubstr: "unsupported source type",
		},
		{
			name: "invalid dest type",
			cfg: &Config{
				Source: SourceConfig{
					Type:      provider.GitHub,
					GitConfig: GitConfig{Repo: "o/r", Token: "tok"},
				},
				Destination: DestinationConfig{Type: "invalid"},
			},
			wantErr:   true,
			errSubstr: "unsupported destination type",
		},
		{
			name: "github source missing repo",
			cfg: &Config{
				Source: SourceConfig{
					Type:      provider.GitHub,
					GitConfig: GitConfig{Token: "tok"},
				},
				Destination: DestinationConfig{
					Type: provider.File,
					Path: "/tmp/out.json",
				},
			},
			wantErr:   true,
			errSubstr: "source.repo is required",
		},
		{
			name: "github source missing token",
			cfg: &Config{
				Source: SourceConfig{
					Type:      provider.GitHub,
					GitConfig: GitConfig{Repo: "o/r"},
				},
				Destination: DestinationConfig{
					Type: provider.File,
					Path: "/tmp/out.json",
				},
			},
			wantErr:   true,
			errSubstr: "source.token is required",
		},
		{
			name: "file dest missing path",
			cfg: &Config{
				Source: SourceConfig{
					Type:      provider.GitHub,
					GitConfig: GitConfig{Repo: "o/r", Token: "tok"},
				},
				Destination: DestinationConfig{
					Type: provider.File,
				},
			},
			wantErr:   true,
			errSubstr: "destination.path is required",
		},
		{
			name: "invalid conflict strategy",
			cfg: &Config{
				Source: SourceConfig{
					Type:      provider.GitHub,
					GitConfig: GitConfig{Repo: "o/r", Token: "tok"},
				},
				Destination: DestinationConfig{
					Type:             provider.File,
					Path:             "/tmp/out.json",
					ConflictStrategy: "invalid",
				},
			},
			wantErr:   true,
			errSubstr: "unsupported conflict strategy",
		},
		{
			name: "default conflict strategy is replace",
			cfg: &Config{
				Source: SourceConfig{
					Type:      provider.GitHub,
					GitConfig: GitConfig{Repo: "o/r", Token: "tok"},
				},
				Destination: DestinationConfig{
					Type: provider.File,
					Path: "/tmp/out.json",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newTestConfig(t, tt.cfg)
			err := cfg.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errSubstr)
				return
			}
			require.NoError(t, err)
		})
	}
}
