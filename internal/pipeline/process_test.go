package pipeline

import (
	"testing"

	"github.com/PapaDanielVi/secret-shift/internal/config"
	"github.com/PapaDanielVi/secret-shift/internal/source"
)

func TestProcess_Passthrough(t *testing.T) {
	p := NewProcessor(config.ProcessConfig{})
	input := []source.Secret{
		{Name: "DB_HOST", Value: "localhost", Type: "env"},
		{Name: "API_KEY", Value: "abc123", Type: "secret"},
	}
	result := p.Process(input)
	if len(result) != 2 {
		t.Fatalf("expected 2 secrets, got %d", len(result))
	}
	if result[0].Name != "DB_HOST" || result[1].Name != "API_KEY" {
		t.Errorf("unexpected names: %s, %s", result[0].Name, result[1].Name)
	}
}

func TestProcess_AddPrefix(t *testing.T) {
	p := NewProcessor(config.ProcessConfig{AddPrefix: "PROD_"})
	input := []source.Secret{
		{Name: "DB_HOST", Value: "localhost", Type: "env"},
	}
	result := p.Process(input)
	if len(result) != 1 {
		t.Fatalf("expected 1 secret, got %d", len(result))
	}
	if result[0].Name != "PROD_DB_HOST" {
		t.Errorf("expected PROD_DB_HOST, got %s", result[0].Name)
	}
}

func TestProcess_AddSuffix(t *testing.T) {
	p := NewProcessor(config.ProcessConfig{AddSuffix: "_V2"})
	input := []source.Secret{
		{Name: "DB_HOST", Value: "localhost", Type: "env"},
	}
	result := p.Process(input)
	if result[0].Name != "DB_HOST_V2" {
		t.Errorf("expected DB_HOST_V2, got %s", result[0].Name)
	}
}

func TestProcess_IncludeRegex(t *testing.T) {
	p := NewProcessor(config.ProcessConfig{IncludeRegex: "^DB_"})
	input := []source.Secret{
		{Name: "DB_HOST", Value: "localhost", Type: "env"},
		{Name: "API_KEY", Value: "abc", Type: "secret"},
		{Name: "DB_PORT", Value: "5432", Type: "env"},
	}
	result := p.Process(input)
	if len(result) != 2 {
		t.Fatalf("expected 2 secrets, got %d", len(result))
	}
	for _, s := range result {
		if s.Name != "DB_HOST" && s.Name != "DB_PORT" {
			t.Errorf("unexpected secret: %s", s.Name)
		}
	}
}

func TestProcess_ExcludeRegex(t *testing.T) {
	p := NewProcessor(config.ProcessConfig{ExcludeRegex: "^DEBUG_"})
	input := []source.Secret{
		{Name: "DB_HOST", Value: "localhost", Type: "env"},
		{Name: "DEBUG_MODE", Value: "true", Type: "env"},
	}
	result := p.Process(input)
	if len(result) != 1 {
		t.Fatalf("expected 1 secret, got %d", len(result))
	}
	if result[0].Name != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", result[0].Name)
	}
}

func TestProcess_IncludeTypes(t *testing.T) {
	p := NewProcessor(config.ProcessConfig{IncludeTypes: []string{"secret"}})
	input := []source.Secret{
		{Name: "DB_HOST", Value: "localhost", Type: "env"},
		{Name: "API_KEY", Value: "abc", Type: "secret"},
	}
	result := p.Process(input)
	if len(result) != 1 {
		t.Fatalf("expected 1 secret, got %d", len(result))
	}
	if result[0].Name != "API_KEY" {
		t.Errorf("expected API_KEY, got %s", result[0].Name)
	}
}

func TestProcess_ExcludeTypes(t *testing.T) {
	p := NewProcessor(config.ProcessConfig{ExcludeTypes: []string{"env"}})
	input := []source.Secret{
		{Name: "DB_HOST", Value: "localhost", Type: "env"},
		{Name: "API_KEY", Value: "abc", Type: "secret"},
	}
	result := p.Process(input)
	if len(result) != 1 {
		t.Fatalf("expected 1 secret, got %d", len(result))
	}
	if result[0].Name != "API_KEY" {
		t.Errorf("expected API_KEY, got %s", result[0].Name)
	}
}

func TestProcess_Combined(t *testing.T) {
	p := NewProcessor(config.ProcessConfig{
		AddPrefix:    "PROD_",
		IncludeRegex: "^DB_",
		IncludeTypes: []string{"env"},
	})
	input := []source.Secret{
		{Name: "DB_HOST", Value: "localhost", Type: "env"},
		{Name: "DB_PORT", Value: "5432", Type: "env"},
		{Name: "API_KEY", Value: "abc", Type: "secret"},
		{Name: "CACHE_HOST", Value: "redis", Type: "env"},
	}
	result := p.Process(input)
	if len(result) != 2 {
		t.Fatalf("expected 2 secrets, got %d", len(result))
	}
	if result[0].Name != "PROD_DB_HOST" {
		t.Errorf("expected PROD_DB_HOST, got %s", result[0].Name)
	}
	if result[1].Name != "PROD_DB_PORT" {
		t.Errorf("expected PROD_DB_PORT, got %s", result[1].Name)
	}
}
