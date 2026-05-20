package pipeline

import (
	"context"
	"fmt"

	"github.com/PapaDanielVi/secret-shift/internal/config"
	"github.com/PapaDanielVi/secret-shift/internal/destination"
	"github.com/PapaDanielVi/secret-shift/internal/source"
	dstfile "github.com/PapaDanielVi/secret-shift/internal/destination/file"
	provideretcd "github.com/PapaDanielVi/secret-shift/internal/provider/etcd"
	providergithub "github.com/PapaDanielVi/secret-shift/internal/provider/github"
	providergitlab "github.com/PapaDanielVi/secret-shift/internal/provider/gitlab"
	providerk8s "github.com/PapaDanielVi/secret-shift/internal/provider/kubernetes"
	providervault "github.com/PapaDanielVi/secret-shift/internal/provider/vault"
)

type Pipeline struct {
	src  source.Source
	dst  destination.Destination
	proc *Processor
}

func (p *Pipeline) Run(ctx context.Context) error {
	fmt.Println("Starting sync pipeline...")

	secrets, err := p.src.Read(ctx)
	if err != nil {
		return fmt.Errorf("read from source: %w", err)
	}
	fmt.Printf("Read %d secrets from source\n", len(secrets))

	processed := p.proc.Process(secrets)
	fmt.Printf("After processing: %d secrets\n", len(processed))

	if err := p.dst.Write(ctx, processed); err != nil {
		return fmt.Errorf("write to destination: %w", err)
	}

	fmt.Println("Sync complete")
	return nil
}

func Build(cfg *config.Config) (*Pipeline, error) {
	p := &Pipeline{}

	switch cfg.Source.Type {
	case "github":
		src, err := providergithub.New(cfg.Source.Token, cfg.Source.Repo, cfg.Source.URL, "replace")
		if err != nil {
			return nil, fmt.Errorf("create github source: %w", err)
		}
		p.src = src
	case "gitlab":
		src, err := providergitlab.New(cfg.Source.Token, cfg.Source.ProjectID, cfg.Source.URL, "replace")
		if err != nil {
			return nil, fmt.Errorf("create gitlab source: %w", err)
		}
		p.src = src
	case "vault":
		src, err := providervault.New(cfg.Source.Token, cfg.Source.VaultAddress, cfg.Source.VaultPath, cfg.Source.VaultMount)
		if err != nil {
			return nil, fmt.Errorf("create vault source: %w", err)
		}
		p.src = src
	case "etcd":
		src, err := provideretcd.New(cfg.Source.EtcdEndpoints, cfg.Source.EtcdPrefix, cfg.Source.EtcdUsername, cfg.Source.EtcdPassword)
		if err != nil {
			return nil, fmt.Errorf("create etcd source: %w", err)
		}
		p.src = src
	case "kubernetes":
		src, err := providerk8s.New(cfg.Source.KubeConfig, cfg.Source.KubeNamespace, cfg.Source.KubeSecretName, cfg.Source.KubeLabel)
		if err != nil {
			return nil, fmt.Errorf("create kubernetes source: %w", err)
		}
		p.src = src
	default:
		return nil, fmt.Errorf("unsupported source type: %s", cfg.Source.Type)
	}

	switch cfg.Destination.Type {
	case "file":
		p.dst = dstfile.New(cfg.Destination.Path, cfg.Destination.Format, cfg.Destination.Encrypt, cfg.Destination.EncryptKey)
	case "github":
		dst, err := providergithub.New(cfg.Destination.Token, cfg.Destination.Repo, cfg.Destination.URL, cfg.Destination.ConflictStrategy)
		if err != nil {
			return nil, fmt.Errorf("create github destination: %w", err)
		}
		p.dst = dst
	case "gitlab":
		dst, err := providergitlab.New(cfg.Destination.Token, cfg.Destination.ProjectID, cfg.Destination.URL, cfg.Destination.ConflictStrategy)
		if err != nil {
			return nil, fmt.Errorf("create gitlab destination: %w", err)
		}
		p.dst = dst
	case "vault":
		dst, err := providervault.New(cfg.Destination.Token, cfg.Destination.VaultAddress, cfg.Destination.VaultPath, cfg.Destination.VaultMount)
		if err != nil {
			return nil, fmt.Errorf("create vault destination: %w", err)
		}
		p.dst = dst
	case "etcd":
		dst, err := provideretcd.New(cfg.Destination.EtcdEndpoints, cfg.Destination.EtcdPrefix, cfg.Destination.EtcdUsername, cfg.Destination.EtcdPassword)
		if err != nil {
			return nil, fmt.Errorf("create etcd destination: %w", err)
		}
		p.dst = dst
	case "kubernetes":
		dst, err := providerk8s.New(cfg.Destination.KubeConfig, cfg.Destination.KubeNamespace, cfg.Destination.KubeSecretName, cfg.Destination.KubeLabel)
		if err != nil {
			return nil, fmt.Errorf("create kubernetes destination: %w", err)
		}
		p.dst = dst
	default:
		return nil, fmt.Errorf("unsupported destination type: %s", cfg.Destination.Type)
	}

	p.proc = NewProcessor(cfg.Process)

	return p, nil
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
