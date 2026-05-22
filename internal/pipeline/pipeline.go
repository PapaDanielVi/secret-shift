package pipeline

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/PapaDanielVi/secret-shift/internal/config"
	"github.com/PapaDanielVi/secret-shift/internal/provider"
)

type Pipeline struct {
	src    provider.Source
	dst    provider.Destination
	proc   *Processor
	dryRun bool
}

func (p *Pipeline) Run(ctx context.Context) error {
	slog.Info("Starting sync pipeline")

	secrets, err := p.src.Read(ctx)
	if err != nil {
		return fmt.Errorf("read from source: %w", err)
	}
	slog.Info("Read secrets from source", "count", len(secrets))

	processed := p.proc.Process(secrets)
	slog.Info("After processing", "count", len(processed))

	if p.dryRun {
		slog.Info("Dry run — skipping write", "count", len(processed))
		for _, s := range processed {
			slog.Debug("Would write secret", "name", s.Name, "type", s.Type)
		}
		return nil
	}

	if err := p.dst.Write(ctx, processed); err != nil {
		return fmt.Errorf("write to destination: %w", err)
	}

	slog.Info("Sync complete")
	return nil
}

func Build(ctx context.Context, cfg *config.Config) (*Pipeline, error) {
	p := &Pipeline{
		dryRun: cfg.DryRun,
	}

	if err := p.initSource(ctx, cfg); err != nil {
		return nil, err
	}
	if err := p.initDestination(ctx, cfg); err != nil {
		return nil, err
	}

	proc, err := NewProcessor(cfg.Process)
	if err != nil {
		return nil, fmt.Errorf("create processor: %w", err)
	}
	p.proc = proc
	return p, nil
}

func (p *Pipeline) initSource(ctx context.Context, cfg *config.Config) error {
	opts := buildSourceOpts(cfg.Source)

	var err error
	p.src, err = provider.CreateSource(ctx, cfg.Source.Type, opts)
	if err != nil {
		return fmt.Errorf("create source: %w", err)
	}
	return nil
}

func (p *Pipeline) initDestination(ctx context.Context, cfg *config.Config) error {
	opts := buildDestOpts(cfg.Destination)

	var err error
	p.dst, err = provider.CreateDestination(ctx, cfg.Destination.Type, opts)
	if err != nil {
		return fmt.Errorf("create destination: %w", err)
	}
	return nil
}

func buildSourceOpts(src config.SourceConfig) map[string]any {
	opts := map[string]any{
		"token":      src.Token,
		"project_id": src.ProjectID,
	}

	switch src.Type {
	case provider.GitHub:
		opts["repo"] = src.Repo
		opts["url"] = src.URL
	case provider.GitLab:
		opts["url"] = src.URL
	case provider.Vault:
		opts["vault_address"] = src.VaultAddress
		opts["vault_path"] = src.VaultPath
		opts["vault_mount"] = src.VaultMount
	case provider.Etcd:
		opts["etcd_endpoints"] = src.EtcdEndpoints
		opts["etcd_prefix"] = src.EtcdPrefix
		opts["etcd_username"] = src.EtcdUsername
		opts["etcd_password"] = src.EtcdPassword
	case provider.Kubernetes:
		opts["kube_config"] = src.KubeConfig
		opts["kube_namespace"] = src.KubeNamespace
		opts["kube_secret_name"] = src.KubeSecretName
		opts["kube_label"] = src.KubeLabel
	case provider.File:
		opts["path"] = src.Path
		opts["format"] = src.Format
		opts["encrypt"] = src.Encrypt
		opts["encrypt_key"] = src.EncryptKey
	}

	return opts
}

func buildDestOpts(dst config.DestinationConfig) map[string]any {
	opts := map[string]any{
		"token":      dst.Token,
		"project_id": dst.ProjectID,
		"strategy":   dst.ConflictStrategy,
	}

	switch dst.Type {
	case provider.GitHub:
		opts["repo"] = dst.Repo
		opts["url"] = dst.URL
	case provider.GitLab:
		opts["url"] = dst.URL
	case provider.Vault:
		opts["vault_address"] = dst.VaultAddress
		opts["vault_path"] = dst.VaultPath
		opts["vault_mount"] = dst.VaultMount
	case provider.Etcd:
		opts["etcd_endpoints"] = dst.EtcdEndpoints
		opts["etcd_prefix"] = dst.EtcdPrefix
		opts["etcd_username"] = dst.EtcdUsername
		opts["etcd_password"] = dst.EtcdPassword
	case provider.Kubernetes:
		opts["kube_config"] = dst.KubeConfig
		opts["kube_namespace"] = dst.KubeNamespace
		opts["kube_secret_name"] = dst.KubeSecretName
		opts["kube_label"] = dst.KubeLabel
	case provider.File:
		opts["path"] = dst.Path
		opts["format"] = dst.Format
		opts["encrypt"] = dst.Encrypt
		opts["encrypt_key"] = dst.EncryptKey
	}

	return opts
}
