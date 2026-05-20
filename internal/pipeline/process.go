package pipeline

import (
	"regexp"
	"slices"

	"github.com/PapaDanielVi/secret-shift/internal/config"
	"github.com/PapaDanielVi/secret-shift/internal/source"
)

type Processor struct {
	cfg config.ProcessConfig
}

func NewProcessor(cfg config.ProcessConfig) *Processor {
	return &Processor{cfg: cfg}
}

func (p *Processor) Process(secrets []source.Secret) []source.Secret {
	var result []source.Secret

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

		result = append(result, source.Secret{
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
	if p.cfg.ExcludeRegex == "" {
		return false
	}
	re, err := regexp.Compile(p.cfg.ExcludeRegex)
	if err != nil {
		return false
	}
	return re.MatchString(name)
}

func (p *Processor) matchInclude(name string) bool {
	if p.cfg.IncludeRegex == "" {
		return true
	}
	re, err := regexp.Compile(p.cfg.IncludeRegex)
	if err != nil {
		return true
	}
	return re.MatchString(name)
}
