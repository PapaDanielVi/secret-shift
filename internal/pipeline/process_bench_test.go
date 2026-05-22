package pipeline

import (
	"fmt"
	"testing"

	"github.com/PapaDanielVi/secret-shift/internal/config"
	"github.com/PapaDanielVi/secret-shift/internal/provider"
)

func BenchmarkProcessor_Process_NoFilters(b *testing.B) {
	secrets := make([]provider.Secret, 100)
	for i := range secrets {
		secrets[i] = provider.Secret{
			Name:  fmt.Sprintf("SECRET_%d", i),
			Value: "value",
			Type:  "env",
		}
	}
	proc, err := NewProcessor(config.ProcessConfig{})
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		_ = proc.Process(secrets)
	}
}

func BenchmarkProcessor_Process_RegexFilter(b *testing.B) {
	secrets := make([]provider.Secret, 100)
	for i := range secrets {
		secrets[i] = provider.Secret{
			Name:  fmt.Sprintf("SECRET_%d", i),
			Value: "value",
			Type:  "env",
		}
	}
	proc, err := NewProcessor(config.ProcessConfig{
		IncludeRegex: "^SECRET_[0-4]",
	})
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		_ = proc.Process(secrets)
	}
}
