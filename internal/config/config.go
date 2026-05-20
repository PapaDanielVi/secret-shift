package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	typeGithub      = "github"
	typeGitlab      = "gitlab"
	typeVault       = "vault"
	typeEtcd        = "etcd"
	typeKubernetes  = "kubernetes"
	typeFile        = "file"
)

type Config struct {
	Source      SourceConfig      `mapstructure:"source"`
	Process     ProcessConfig     `mapstructure:"process"`
	Destination DestinationConfig `mapstructure:"destination"`
}

type SourceConfig struct {
	Type      string `mapstructure:"type"`
	Repo      string `mapstructure:"repo"`
	TokenEnv  string `mapstructure:"token_env"`
	Token     string `mapstructure:"token"`
	URL       string `mapstructure:"url"`
	ProjectID string `mapstructure:"project_id"`

	// Vault-specific
	VaultAddress string `mapstructure:"vault_address"`
	VaultPath    string `mapstructure:"vault_path"`
	VaultMount   string `mapstructure:"vault_mount"`

	// etcd-specific
	EtcdEndpoints []string `mapstructure:"etcd_endpoints"`
	EtcdPrefix   string   `mapstructure:"etcd_prefix"`
	EtcdUsername string   `mapstructure:"etcd_username"`
	EtcdPassword string   `mapstructure:"etcd_password"`

	// Kubernetes-specific
	KubeNamespace  string `mapstructure:"kube_namespace"`
	KubeSecretName string `mapstructure:"kube_secret_name"`
	KubeConfig     string `mapstructure:"kube_config"`
	KubeLabel      string `mapstructure:"kube_label"`
}

type ProcessConfig struct {
	AddPrefix    string   `mapstructure:"add_prefix"`
	AddSuffix    string   `mapstructure:"add_suffix"`
	IncludeRegex string   `mapstructure:"include_regex"`
	ExcludeRegex string   `mapstructure:"exclude_regex"`
	IncludeTypes []string `mapstructure:"include_types"`
	ExcludeTypes []string `mapstructure:"exclude_types"`
}

type DestinationConfig struct {
	Type             string `mapstructure:"type"`
	Path             string `mapstructure:"path"`
	Format           string `mapstructure:"format"`
	Repo             string `mapstructure:"repo"`
	TokenEnv         string `mapstructure:"token_env"`
	Token            string `mapstructure:"token"`
	URL              string `mapstructure:"url"`
	ProjectID        string `mapstructure:"project_id"`
	ConflictStrategy string `mapstructure:"conflict_strategy"`

	// File encryption
	Encrypt    bool   `mapstructure:"encrypt"`
	EncryptKey string `mapstructure:"encrypt_key"`

	// Vault-specific
	VaultAddress string `mapstructure:"vault_address"`
	VaultPath    string `mapstructure:"vault_path"`
	VaultMount   string `mapstructure:"vault_mount"`

	// etcd-specific
	EtcdEndpoints []string `mapstructure:"etcd_endpoints"`
	EtcdPrefix   string   `mapstructure:"etcd_prefix"`
	EtcdUsername string   `mapstructure:"etcd_username"`
	EtcdPassword string   `mapstructure:"etcd_password"`

	// Kubernetes-specific
	KubeNamespace  string `mapstructure:"kube_namespace"`
	KubeSecretName string `mapstructure:"kube_secret_name"`
	KubeConfig     string `mapstructure:"kube_config"`
	KubeLabel      string `mapstructure:"kube_label"`
}

var validSourceTypes = map[string]bool{ //nolint:gochecknoglobals // global config state
	typeGithub: true, typeGitlab: true, typeVault: true, typeEtcd: true, typeKubernetes: true,
}

var validDestTypes = map[string]bool{ //nolint:gochecknoglobals // global config state
	typeFile: true, typeGithub: true, typeGitlab: true, typeVault: true, typeEtcd: true, typeKubernetes: true,
}

func Load(configFile string) (*Config, error) {
	if configFile != "" {
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("read config file: %w", err)
		}
	}

	viper.SetEnvPrefix("SECRET_SHIFT")
	viper.AutomaticEnv()

	for _, key := range []string{
		"source.type", "source.repo", "source.token", "source.url", "source.project_id",
		"source.vault_address", "source.vault_path", "source.vault_mount",
		"source.etcd_endpoints", "source.etcd_prefix", "source.etcd_username", "source.etcd_password",
		"source.kube_namespace", "source.kube_secret_name", "source.kube_config", "source.kube_label",
		"destination.type", "destination.path", "destination.format", "destination.repo",
		"destination.token", "destination.url", "destination.project_id", "destination.conflict_strategy",
		"destination.encrypt", "destination.encrypt_key",
		"destination.vault_address", "destination.vault_path", "destination.vault_mount",
		"destination.etcd_endpoints", "destination.etcd_prefix", "destination.etcd_username", "destination.etcd_password",
		"destination.kube_namespace", "destination.kube_secret_name", "destination.kube_config", "destination.kube_label",
		"process.add_prefix", "process.add_suffix", "process.include_regex", "process.exclude_regex",
	} {
		_ = viper.BindEnv(key)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
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
		c.Destination.ConflictStrategy = "replace"
	}
	validStrategy := map[string]bool{"replace": true, "skip": true, "report": true}
	if !validStrategy[c.Destination.ConflictStrategy] {
		return fmt.Errorf("unsupported conflict strategy: %s", c.Destination.ConflictStrategy)
	}

	return nil
}

func validateSourceFields(s *SourceConfig) error {
	switch s.Type {
	case typeGithub:
		if s.Repo == "" {
			return errors.New("source.repo is required for github source")
		}
		if s.Token == "" {
			return errors.New("source token is required for github source")
		}
	case typeGitlab:
		if s.ProjectID == "" {
			return errors.New("source.project_id is required for gitlab source")
		}
		if s.Token == "" {
			return errors.New("source token is required for gitlab source")
		}
	case typeVault:
		if s.VaultAddress == "" {
			return errors.New("source.vault_address is required for vault source")
		}
		if s.VaultPath == "" {
			return errors.New("source.vault_path is required for vault source")
		}
	case typeEtcd:
		if len(s.EtcdEndpoints) == 0 {
			return errors.New("source.etcd_endpoints is required for etcd source")
		}
	case typeKubernetes:
		if s.KubeNamespace == "" {
			return errors.New("source.kube_namespace is required for kubernetes source")
		}
	}
	return nil
}

func validateDestFields(d *DestinationConfig) error {
	switch d.Type {
	case typeFile:
		if d.Path == "" {
			return errors.New("destination.path is required for file destination")
		}
	case typeGithub:
		if d.Repo == "" {
			return errors.New("destination.repo is required for github destination")
		}
		if d.Token == "" {
			return errors.New("destination token is required for github destination")
		}
	case typeGitlab:
		if d.ProjectID == "" {
			return errors.New("destination.project_id is required for gitlab destination")
		}
		if d.Token == "" {
			return errors.New("destination token is required for gitlab destination")
		}
	case typeVault:
		if d.VaultAddress == "" {
			return errors.New("destination.vault_address is required for vault destination")
		}
		if d.VaultPath == "" {
			return errors.New("destination.vault_path is required for vault destination")
		}
	case typeEtcd:
		if len(d.EtcdEndpoints) == 0 {
			return errors.New("destination.etcd_endpoints is required for etcd destination")
		}
	case typeKubernetes:
		if d.KubeNamespace == "" {
			return errors.New("destination.kube_namespace is required for kubernetes destination")
		}
	}
	return nil
}
