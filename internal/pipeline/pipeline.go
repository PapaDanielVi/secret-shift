package pipeline

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/PapaDanielVi/secret-shift/internal/config"
	dstfile "github.com/PapaDanielVi/secret-shift/internal/destination/file"
	"github.com/PapaDanielVi/secret-shift/internal/provider"
	provideretcd "github.com/PapaDanielVi/secret-shift/internal/provider/etcd"
	providergithub "github.com/PapaDanielVi/secret-shift/internal/provider/github"
	providergitlab "github.com/PapaDanielVi/secret-shift/internal/provider/gitlab"
	providerk8s "github.com/PapaDanielVi/secret-shift/internal/provider/kubernetes"
	providervault "github.com/PapaDanielVi/secret-shift/internal/provider/vault"
)

type Pipeline struct {
	src  provider.Source
	dst  provider.Destination
	proc *Processor
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

	if err := p.dst.Write(ctx, processed); err != nil {
		return fmt.Errorf("write to destination: %w", err)
	}

	slog.Info("Sync complete")
	return nil
}

func Build(ctx context.Context, cfg *config.Config) (*Pipeline, error) {
	p := &Pipeline{}

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
	var err error
	src := cfg.Source
	switch src.Type {
	case provider.GitHub:
		p.src, err = providergithub.New(ctx, src.Token, src.Repo, src.URL, "replace")
	case provider.GitLab:
		p.src, err = providergitlab.New(ctx, src.Token, src.ProjectID, src.URL, "replace")
	case provider.Vault:
		p.src, err = providervault.New(src.Token, src.VaultAddress, src.VaultPath, src.VaultMount)
	case provider.Etcd:
		p.src, err = provideretcd.New(src.EtcdEndpoints, src.EtcdPrefix, src.EtcdUsername, src.EtcdPassword)
	case provider.Kubernetes:
		p.src, err = providerk8s.New(src.KubeConfig, src.KubeNamespace, src.KubeSecretName, src.KubeLabel)
	default:
		return fmt.Errorf("unsupported source type: %s", src.Type)
	}
	if err != nil {
		return fmt.Errorf("create source: %w", err)
	}
	return nil
}

func (p *Pipeline) initDestination(ctx context.Context, cfg *config.Config) error {
	var err error
	dst := cfg.Destination
	switch dst.Type {
	case provider.File:
		p.dst = dstfile.New(dst.Path, dst.Format, dst.Encrypt, dst.EncryptKey)
	case provider.GitHub:
		p.dst, err = providergithub.New(ctx, dst.Token, dst.Repo, dst.URL, dst.ConflictStrategy)
	case provider.GitLab:
		p.dst, err = providergitlab.New(ctx, dst.Token, dst.ProjectID, dst.URL, dst.ConflictStrategy)
	case provider.Vault:
		p.dst, err = providervault.New(dst.Token, dst.VaultAddress, dst.VaultPath, dst.VaultMount)
	case provider.Etcd:
		p.dst, err = provideretcd.New(dst.EtcdEndpoints, dst.EtcdPrefix, dst.EtcdUsername, dst.EtcdPassword)
	case provider.Kubernetes:
		p.dst, err = providerk8s.New(dst.KubeConfig, dst.KubeNamespace, dst.KubeSecretName, dst.KubeLabel)
	default:
		return fmt.Errorf("unsupported destination type: %s", dst.Type)
	}
	if err != nil {
		return fmt.Errorf("create destination: %w", err)
	}
	return nil
}
