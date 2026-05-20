package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
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

	// Explicit env bindings for direct token passing
	for _, key := range []string{
		"source.type", "source.repo", "source.token", "source.url", "source.project_id",
		"destination.type", "destination.path", "destination.format", "destination.repo",
		"destination.token", "destination.url", "destination.project_id", "destination.conflict_strategy",
		"process.add_prefix", "process.add_suffix", "process.include_regex", "process.exclude_regex",
	} {
		_ = viper.BindEnv(key)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// Resolve tokens from env if not set directly
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
		return fmt.Errorf("source.type is required")
	}
	if c.Source.Type != "github" && c.Source.Type != "gitlab" {
		return fmt.Errorf("unsupported source type: %s", c.Source.Type)
	}
	if c.Source.Repo == "" && c.Source.ProjectID == "" {
		return fmt.Errorf("source.repo or source.project_id is required")
	}
	if c.Source.Token == "" {
		return fmt.Errorf("source token is required (set via config or %s)", c.Source.TokenEnv)
	}

	if c.Destination.Type == "" {
		return fmt.Errorf("destination.type is required")
	}
	validDest := map[string]bool{"file": true, "github": true, "gitlab": true}
	if !validDest[c.Destination.Type] {
		return fmt.Errorf("unsupported destination type: %s", c.Destination.Type)
	}
	if c.Destination.Type == "file" && c.Destination.Path == "" {
		return fmt.Errorf("destination.path is required for file destination")
	}
	if c.Destination.Type != "file" && c.Destination.Token == "" {
		return fmt.Errorf("destination token is required for %s destination", c.Destination.Type)
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
