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

type VaultConfig struct {
	VaultAddress string `json:"vault_address"`
	VaultPath    string `json:"vault_path"`
	VaultMount   string `json:"vault_mount"`
}

type EtcdConfig struct {
	EtcdEndpoints []string `json:"etcd_endpoints"`
	EtcdPrefix    string   `json:"etcd_prefix"`
	EtcdUsername  string   `json:"etcd_username"`
	EtcdPassword  string   `json:"etcd_password"`
}

type KubernetesConfig struct {
	KubeNamespace  string `json:"kube_namespace"`
	KubeSecretName string `json:"kube_secret_name"`
	KubeConfig     string `json:"kube_config"`
	KubeLabel      string `json:"kube_label"`
}

type GitConfig struct {
	Repo     string `json:"repo"`
	TokenEnv string `json:"token_env"`
	Token    string `json:"token"`
	URL      string `json:"url"`
}

type Config struct {
	Source      SourceConfig      `json:"source"`
	Process     ProcessConfig     `json:"process"`
	Destination DestinationConfig `json:"destination"`
}

type SourceConfig struct {
	GitConfig
	VaultConfig
	EtcdConfig
	KubernetesConfig

	Type      provider.Type `json:"type"`
	ProjectID string        `json:"project_id"`
}

type ProcessConfig struct {
	AddPrefix    string   `json:"add_prefix"`
	AddSuffix    string   `json:"add_suffix"`
	IncludeRegex string   `json:"include_regex"`
	ExcludeRegex string   `json:"exclude_regex"`
	IncludeTypes []string `json:"include_types"`
	ExcludeTypes []string `json:"exclude_types"`
}

type DestinationConfig struct {
	GitConfig
	VaultConfig
	EtcdConfig
	KubernetesConfig

	Type             provider.Type `json:"type"`
	Path             string        `json:"path"`
	Format           string        `json:"format"`
	ConflictStrategy string        `json:"conflict_strategy"`
	Encrypt          bool          `json:"encrypt"`
	EncryptKey       string        `json:"encrypt_key"`
	ProjectID        string        `json:"project_id"`
}

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

	if cfg.Source.Token == "" && cfg.Source.TokenEnv != "" {
		cfg.Source.Token = os.Getenv(cfg.Source.TokenEnv)
	}
	if cfg.Destination.Token == "" && cfg.Destination.TokenEnv != "" {
		cfg.Destination.Token = os.Getenv(cfg.Destination.TokenEnv)
	}

	return &cfg, nil
}

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
			return errors.New("source.token is required for github source")
		}
	case provider.GitLab:
		if s.ProjectID == "" {
			return errors.New("source.project_id is required for gitlab source")
		}
		if s.Token == "" {
			return errors.New("source.token is required for gitlab source")
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
	default:
		return nil
	}
	return nil
}

//exhaustive:ignore
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
			return errors.New("destination.token is required for github destination")
		}
	case provider.GitLab:
		if d.ProjectID == "" {
			return errors.New("destination.project_id is required for gitlab destination")
		}
		if d.Token == "" {
			return errors.New("destination.token is required for gitlab destination")
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
