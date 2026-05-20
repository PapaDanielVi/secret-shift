package pipeline

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/PapaDanielVi/secret-shift/internal/config"
	"github.com/PapaDanielVi/secret-shift/internal/destination"
	dstfile "github.com/PapaDanielVi/secret-shift/internal/destination/file"
	provideretcd "github.com/PapaDanielVi/secret-shift/internal/provider/etcd"
	providergithub "github.com/PapaDanielVi/secret-shift/internal/provider/github"
	providergitlab "github.com/PapaDanielVi/secret-shift/internal/provider/gitlab"
	providerk8s "github.com/PapaDanielVi/secret-shift/internal/provider/kubernetes"
	providervault "github.com/PapaDanielVi/secret-shift/internal/provider/vault"
	"github.com/PapaDanielVi/secret-shift/internal/source"
)

type Pipeline struct {
	src  source.Source
	dst  destination.Destination
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

func Build(cfg *config.Config) (*Pipeline, error) {
	p := &Pipeline{}

	if err := p.initSource(cfg); err != nil {
		return nil, err
	}
	if err := p.initDestination(cfg); err != nil {
		return nil, err
	}

	p.proc = NewProcessor(cfg.Process)
	return p, nil
}

func (p *Pipeline) initSource(cfg *config.Config) error {
	var err error
	switch cfg.Source.Type {
	case "github":
		p.src, err = providergithub.New(cfg.Source.Token, cfg.Source.Repo, cfg.Source.URL, "replace")
	case "gitlab":
		p.src, err = providergitlab.New(cfg.Source.Token, cfg.Source.ProjectID, cfg.Source.URL, "replace")
	case "vault":
		p.src, err = providervault.New(cfg.Source.Token, cfg.Source.VaultAddress, cfg.Source.VaultPath, cfg.Source.VaultMount)
	case "etcd":
		p.src, err = provideretcd.New(cfg.Source.EtcdEndpoints, cfg.Source.EtcdPrefix, cfg.Source.EtcdUsername, cfg.Source.EtcdPassword)
	case "kubernetes":
		p.src, err = providerk8s.New(cfg.Source.KubeConfig, cfg.Source.KubeNamespace, cfg.Source.KubeSecretName, cfg.Source.KubeLabel)
	default:
		return fmt.Errorf("unsupported source type: %s", cfg.Source.Type)
	}
	if err != nil {
		return fmt.Errorf("create source: %w", err)
	}
	return nil
}

func (p *Pipeline) initDestination(cfg *config.Config) error {
	var err error
	switch cfg.Destination.Type {
	case "file":
		p.dst = dstfile.New(cfg.Destination.Path, cfg.Destination.Format, cfg.Destination.Encrypt, cfg.Destination.EncryptKey)
	case "github":
		p.dst, err = providergithub.New(cfg.Destination.Token, cfg.Destination.Repo, cfg.Destination.URL, cfg.Destination.ConflictStrategy)
	case "gitlab":
		p.dst, err = providergitlab.New(cfg.Destination.Token, cfg.Destination.ProjectID, cfg.Destination.URL, cfg.Destination.ConflictStrategy)
	case "vault":
		p.dst, err = providervault.New(cfg.Destination.Token, cfg.Destination.VaultAddress, cfg.Destination.VaultPath, cfg.Destination.VaultMount)
	case "etcd":
		p.dst, err = provideretcd.New(cfg.Destination.EtcdEndpoints, cfg.Destination.EtcdPrefix, cfg.Destination.EtcdUsername, cfg.Destination.EtcdPassword)
	case "kubernetes":
		p.dst, err = providerk8s.New(cfg.Destination.KubeConfig, cfg.Destination.KubeNamespace, cfg.Destination.KubeSecretName, cfg.Destination.KubeLabel)
	default:
		return fmt.Errorf("unsupported destination type: %s", cfg.Destination.Type)
	}
	if err != nil {
		return fmt.Errorf("create destination: %w", err)
	}
	return nil
}

var _ source.Source = (*providergithub.Provider)(nil)
var _ destination.Destination = (*providergithub.Provider)(nil)
var _ source.Source = (*providergitlab.Provider)(nil)
var _ destination.Destination = (*providergitlab.Provider)(nil)
var _ source.Source = (*providervault.Provider)(nil)
var _ destination.Destination = (*providervault.Provider)(nil)
var _ source.Source = (*provideretcd.Provider)(nil)
var _ destination.Destination = (*provideretcd.Provider)(nil)
var _ source.Source = (*providerk8s.Provider)(nil)
var _ destination.Destination = (*providerk8s.Provider)(nil)
