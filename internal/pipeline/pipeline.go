package pipeline

import (
	"context"
	"fmt"

	"github.com/PapaDanielVi/secret-shift/internal/config"
	"github.com/PapaDanielVi/secret-shift/internal/destination"
	dstetcd "github.com/PapaDanielVi/secret-shift/internal/destination/etcd"
	dstfile "github.com/PapaDanielVi/secret-shift/internal/destination/file"
	dstgithub "github.com/PapaDanielVi/secret-shift/internal/destination/github"
	dstgitlab "github.com/PapaDanielVi/secret-shift/internal/destination/gitlab"
	dstk8s "github.com/PapaDanielVi/secret-shift/internal/destination/kubernetes"
	dstvault "github.com/PapaDanielVi/secret-shift/internal/destination/vault"
	"github.com/PapaDanielVi/secret-shift/internal/source"
	srcetcd "github.com/PapaDanielVi/secret-shift/internal/source/etcd"
	srcgithub "github.com/PapaDanielVi/secret-shift/internal/source/github"
	srcgitlab "github.com/PapaDanielVi/secret-shift/internal/source/gitlab"
	srck8s "github.com/PapaDanielVi/secret-shift/internal/source/kubernetes"
	srcvault "github.com/PapaDanielVi/secret-shift/internal/source/vault"
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
		src, err := srcgithub.New(cfg.Source.Token, cfg.Source.Repo, cfg.Source.URL)
		if err != nil {
			return nil, fmt.Errorf("create github source: %w", err)
		}
		p.src = src
	case "gitlab":
		src, err := srcgitlab.New(cfg.Source.Token, cfg.Source.ProjectID, cfg.Source.URL)
		if err != nil {
			return nil, fmt.Errorf("create gitlab source: %w", err)
		}
		p.src = src
	case "vault":
		src, err := srcvault.New(cfg.Source.Token, cfg.Source.VaultAddress, cfg.Source.VaultPath, cfg.Source.VaultMount)
		if err != nil {
			return nil, fmt.Errorf("create vault source: %w", err)
		}
		p.src = src
	case "etcd":
		src, err := srcetcd.New(cfg.Source.EtcdEndpoints, cfg.Source.EtcdPrefix, cfg.Source.EtcdUsername, cfg.Source.EtcdPassword)
		if err != nil {
			return nil, fmt.Errorf("create etcd source: %w", err)
		}
		p.src = src
	case "kubernetes":
		src, err := srck8s.New(cfg.Source.KubeConfig, cfg.Source.KubeNamespace, cfg.Source.KubeSecretName, cfg.Source.KubeLabel)
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
		dst, err := dstgithub.New(cfg.Destination.Token, cfg.Destination.Repo, cfg.Destination.URL, cfg.Destination.ConflictStrategy)
		if err != nil {
			return nil, fmt.Errorf("create github destination: %w", err)
		}
		p.dst = dst
	case "gitlab":
		dst, err := dstgitlab.New(cfg.Destination.Token, cfg.Destination.ProjectID, cfg.Destination.URL, cfg.Destination.ConflictStrategy)
		if err != nil {
			return nil, fmt.Errorf("create gitlab destination: %w", err)
		}
		p.dst = dst
	case "vault":
		dst, err := dstvault.New(cfg.Destination.Token, cfg.Destination.VaultAddress, cfg.Destination.VaultPath, cfg.Destination.VaultMount)
		if err != nil {
			return nil, fmt.Errorf("create vault destination: %w", err)
		}
		p.dst = dst
	case "etcd":
		dst, err := dstetcd.New(cfg.Destination.EtcdEndpoints, cfg.Destination.EtcdPrefix, cfg.Destination.EtcdUsername, cfg.Destination.EtcdPassword)
		if err != nil {
			return nil, fmt.Errorf("create etcd destination: %w", err)
		}
		p.dst = dst
	case "kubernetes":
		dst, err := dstk8s.New(cfg.Destination.KubeConfig, cfg.Destination.KubeNamespace, cfg.Destination.KubeSecretName)
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
