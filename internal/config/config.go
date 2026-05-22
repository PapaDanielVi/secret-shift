package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/PapaDanielVi/secret-shift/internal/provider"
	"github.com/spf13/viper"
)

const (
	strategyReplace = "replace"
	strategySkip    = "skip"
	strategyReport  = "report"
)

// VaultConfig holds Vault-specific configuration.
type VaultConfig struct {
	VaultAddress string `json:"vault_address"`
	VaultPath    string `json:"vault_path"`
	VaultMount   string `json:"vault_mount"`
}

// EtcdConfig holds etcd-specific configuration.
type EtcdConfig struct {
	EtcdEndpoints []string `json:"etcd_endpoints"`
	EtcdPrefix    string   `json:"etcd_prefix"`
	EtcdUsername  string   `json:"etcd_username"`
	EtcdPassword  string   `json:"etcd_password"`
}

// KubernetesConfig holds Kubernetes-specific configuration.
type KubernetesConfig struct {
	KubeNamespace  string `json:"kube_namespace"`
	KubeSecretName string `json:"kube_secret_name"`
	KubeConfig     string `json:"kube_config"`
	KubeLabel      string `json:"kube_label"`
}

// GitConfig holds GitHub/GitLab-specific configuration.
type GitConfig struct {
	Repo     string `json:"repo"`
	TokenEnv string `json:"token_env"`
	Token    string `json:"token"`
	URL      string `json:"url"`
}

// FileConfig holds file-specific configuration.
type FileConfig struct {
	Path       string `json:"path"`
	Format     string `json:"format"`
	Encrypt    bool   `json:"encrypt"`
	EncryptKey string `json:"encrypt_key"`
}

// Config is the top-level configuration.
type Config struct {
	Source      SourceConfig      `json:"source"`
	Process     ProcessConfig     `json:"process"`
	Destination DestinationConfig `json:"destination"`
	DryRun      bool              `json:"dry_run"`
}

// SourceConfig holds the source provider configuration.
type SourceConfig struct {
	Type      provider.Type `json:"type"`
	ProjectID string        `json:"project_id"`

	GitConfig
	VaultConfig
	EtcdConfig
	KubernetesConfig
	FileConfig
}

// ProcessConfig holds the processing/filtering configuration.
type ProcessConfig struct {
	AddPrefix    string   `json:"add_prefix"`
	AddSuffix    string   `json:"add_suffix"`
	IncludeRegex string   `json:"include_regex"`
	ExcludeRegex string   `json:"exclude_regex"`
	IncludeTypes []string `json:"include_types"`
	ExcludeTypes []string `json:"exclude_types"`
}

// DestinationConfig holds the destination provider configuration.
type DestinationConfig struct {
	Type             provider.Type `json:"type"`
	ConflictStrategy string        `json:"conflict_strategy"`
	ProjectID        string        `json:"project_id"`

	GitConfig
	VaultConfig
	EtcdConfig
	KubernetesConfig
	FileConfig
}

// Load reads configuration from file and environment.
// Environment variables use the prefix SECRET_SHIFT_SRC_ for source fields
// and SECRET_SHIFT_DST_ for destination fields, e.g.:
//
//	SECRET_SHIFT_SRC_GITHUB_TOKEN=xxx
//	SECRET_SHIFT_DST_GITHUB_TOKEN=yyy
func Load(v *viper.Viper, configFile string) (*Config, error) {
	if configFile != "" {
		v.SetConfigFile(configFile)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("read config file: %w", err)
		}
	}

	v.SetEnvPrefix("SECRET_SHIFT")
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// Resolve tokens from env vars using the token_env field or
	// the per-step env var convention: SECRET_SHIFT_SRC_<PROVIDER>_TOKEN
	// and SECRET_SHIFT_DST_<PROVIDER>_TOKEN.
	resolveSourceToken(&cfg.Source)
	resolveDestToken(&cfg.Destination)

	return &cfg, nil
}

func resolveSourceToken(s *SourceConfig) {
	if s.Token != "" {
		return
	}
	// Check token_env field first
	if s.TokenEnv != "" {
		s.Token = os.Getenv(s.TokenEnv)
		if s.Token != "" {
			return
		}
	}
	// Check per-step env var: SECRET_SHIFT_SRC_<TYPE>_TOKEN
	envKey := fmt.Sprintf("SECRET_SHIFT_SRC_%s_TOKEN", providerEnvSuffix(s.Type))
	s.Token = os.Getenv(envKey)
}

func resolveDestToken(d *DestinationConfig) {
	if d.Token != "" {
		return
	}
	if d.TokenEnv != "" {
		d.Token = os.Getenv(d.TokenEnv)
		if d.Token != "" {
			return
		}
	}
	envKey := fmt.Sprintf("SECRET_SHIFT_DST_%s_TOKEN", providerEnvSuffix(d.Type))
	d.Token = os.Getenv(envKey)
}

func providerEnvSuffix(t provider.Type) string {
	switch t {
	case provider.GitHub:
		return "GITHUB"
	case provider.GitLab:
		return "GITLAB"
	case provider.Vault:
		return "VAULT"
	case provider.Etcd:
		return "ETCD"
	case provider.Kubernetes:
		return "KUBERNETES"
	case provider.File:
		return "FILE"
	default:
		return string(t)
	}
}

// Validate checks the configuration for correctness.
func (c *Config) Validate() error {
	if c.Source.Type == "" {
		return errors.New("source.type is required")
	}
	if !validSourceTypes[c.Source.Type] {
		return fmt.Errorf("unsupported source type: %s", c.Source.Type)
	}
	if err := validateSourceFields(&c.Source); err != nil {
		return err
	}

	if c.Destination.Type == "" {
		return errors.New("destination.type is required")
	}
	if !validDestTypes[c.Destination.Type] {
		return fmt.Errorf("unsupported destination type: %s", c.Destination.Type)
	}
	if err := validateDestFields(&c.Destination); err != nil {
		return err
	}

	if c.Destination.ConflictStrategy == "" {
		c.Destination.ConflictStrategy = strategyReplace
	}
	validStrategy := map[string]bool{strategyReplace: true, strategySkip: true, strategyReport: true}
	if !validStrategy[c.Destination.ConflictStrategy] {
		return fmt.Errorf("unsupported conflict strategy: %s", c.Destination.ConflictStrategy)
	}

	return nil
}

func validateSourceFields(s *SourceConfig) error {
	switch s.Type {
	case provider.GitHub:
		if s.Repo == "" {
			return errors.New("source.repo is required for github source")
		}
		if s.Token == "" {
			return errors.New("source.token is required for github source (set token, token_env, or SECRET_SHIFT_SRC_GITHUB_TOKEN)")
		}
	case provider.GitLab:
		if s.ProjectID == "" {
			return errors.New("source.project_id is required for gitlab source")
		}
		if s.Token == "" {
			return errors.New("source.token is required for gitlab source (set token, token_env, or SECRET_SHIFT_SRC_GITLAB_TOKEN)")
		}
	case provider.Vault:
		if s.VaultAddress == "" {
			return errors.New("source.vault_address is required for vault source")
		}
		if s.VaultPath == "" {
			return errors.New("source.vault_path is required for vault source")
		}
	case provider.Etcd:
		if len(s.EtcdEndpoints) == 0 {
			return errors.New("source.etcd_endpoints is required for etcd source")
		}
	case provider.Kubernetes:
		if s.KubeNamespace == "" {
			return errors.New("source.kube_namespace is required for kubernetes source")
		}
	case provider.File:
		if s.Path == "" {
			return errors.New("source.path is required for file source")
		}
	}
	return nil
}

func validateDestFields(d *DestinationConfig) error {
	switch d.Type {
	case provider.File:
		if d.Path == "" {
			return errors.New("destination.path is required for file destination")
		}
	case provider.GitHub:
		if d.Repo == "" {
			return errors.New("destination.repo is required for github destination")
		}
		if d.Token == "" {
			return errors.New("destination.token is required for github destination (set token, token_env, or SECRET_SHIFT_DST_GITHUB_TOKEN)")
		}
	case provider.GitLab:
		if d.ProjectID == "" {
			return errors.New("destination.project_id is required for gitlab destination")
		}
		if d.Token == "" {
			return errors.New("destination.token is required for gitlab destination (set token, token_env, or SECRET_SHIFT_DST_GITLAB_TOKEN)")
		}
	case provider.Vault:
		if d.VaultAddress == "" {
			return errors.New("destination.vault_address is required for vault destination")
		}
		if d.VaultPath == "" {
			return errors.New("destination.vault_path is required for vault destination")
		}
	case provider.Etcd:
		if len(d.EtcdEndpoints) == 0 {
			return errors.New("destination.etcd_endpoints is required for etcd destination")
		}
	case provider.Kubernetes:
		if d.KubeNamespace == "" {
			return errors.New("destination.kube_namespace is required for kubernetes destination")
		}
	}
	return nil
}

//nolint:gochecknoglobals
var validSourceTypes = map[provider.Type]bool{
	provider.GitHub:     true,
	provider.GitLab:     true,
	provider.Vault:      true,
	provider.Etcd:       true,
	provider.Kubernetes: true,
	provider.File:       true,
}

//nolint:gochecknoglobals
var validDestTypes = map[provider.Type]bool{
	provider.File:       true,
	provider.GitHub:     true,
	provider.GitLab:     true,
	provider.Vault:      true,
	provider.Etcd:       true,
	provider.Kubernetes: true,
}
