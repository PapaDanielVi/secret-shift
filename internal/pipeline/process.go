package pipeline

import (
	"fmt"
	"regexp"
	"slices"

	"github.com/PapaDanielVi/secret-shift/internal/config"
	"github.com/PapaDanielVi/secret-shift/internal/provider"
)

type Processor struct {
	cfg         config.ProcessConfig
	excludeRe   *regexp.Regexp
	includeRe   *regexp.Regexp
}

func NewProcessor(cfg config.ProcessConfig) (*Processor, error) {
	p := &Processor{cfg: cfg}

	if cfg.ExcludeRegex != "" {
		re, err := regexp.Compile(cfg.ExcludeRegex)
		if err != nil {
			return nil, fmt.Errorf("compile exclude regex: %w", err)
		}
		p.excludeRe = re
	}

	if cfg.IncludeRegex != "" {
		re, err := regexp.Compile(cfg.IncludeRegex)
		if err != nil {
			return nil, fmt.Errorf("compile include regex: %w", err)
		}
		p.includeRe = re
	}

	return p, nil
}

func (p *Processor) Process(secrets []provider.Secret) []provider.Secret {
	var result []provider.Secret

	for _, s := range secrets {
		if p.excludeType(s.Type) {
			continue
		}
		if !p.includeType(s.Type) {
			continue
		}

		name := s.Name

		if p.matchExclude(name) {
			continue
		}
		if !p.matchInclude(name) {
			continue
		}

		name = p.cfg.AddPrefix + name + p.cfg.AddSuffix

		result = append(result, provider.Secret{
			Name:  name,
			Value: s.Value,
			Type:  s.Type,
		})
	}

	return result
}

func (p *Processor) excludeType(t string) bool {
	return slices.Contains(p.cfg.ExcludeTypes, t)
}

func (p *Processor) includeType(t string) bool {
	if len(p.cfg.IncludeTypes) == 0 {
		return true
	}
	return slices.Contains(p.cfg.IncludeTypes, t)
}

func (p *Processor) matchExclude(name string) bool {
	if p.excludeRe == nil {
		return false
	}
	return p.excludeRe.MatchString(name)
}

func (p *Processor) matchInclude(name string) bool {
	if p.includeRe == nil {
		return true
	}
	return p.includeRe.MatchString(name)
}
