package pipeline

import (
	"context"
	"fmt"

	"github.com/PapaDanielVi/secret-shift/internal/config"
	"github.com/PapaDanielVi/secret-shift/internal/destination"
	dstfile "github.com/PapaDanielVi/secret-shift/internal/destination/file"
	dstgithub "github.com/PapaDanielVi/secret-shift/internal/destination/github"
	dstgitlab "github.com/PapaDanielVi/secret-shift/internal/destination/gitlab"
	"github.com/PapaDanielVi/secret-shift/internal/source"
	srcgithub "github.com/PapaDanielVi/secret-shift/internal/source/github"
	srcgitlab "github.com/PapaDanielVi/secret-shift/internal/source/gitlab"
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
	default:
		return nil, fmt.Errorf("unsupported source type: %s", cfg.Source.Type)
	}

	switch cfg.Destination.Type {
	case "file":
		p.dst = dstfile.New(cfg.Destination.Path, cfg.Destination.Format)
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
	default:
		return nil, fmt.Errorf("unsupported destination type: %s", cfg.Destination.Type)
	}

	p.proc = NewProcessor(cfg.Process)

	return p, nil
}
