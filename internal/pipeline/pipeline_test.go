package pipeline

import (
	"context"
	"errors"
	"testing"

	"github.com/PapaDanielVi/secret-shift/internal/config"
	"github.com/PapaDanielVi/secret-shift/internal/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockSource implements provider.Source for testing.
type mockSource struct {
	secrets []provider.Secret
	err     error
}

func (m *mockSource) Read(_ context.Context) ([]provider.Secret, error) {
	return m.secrets, m.err
}

// mockDestination implements provider.Destination for testing.
type mockDestination struct {
	written []provider.Secret
	err     error
}

func (m *mockDestination) Write(_ context.Context, secrets []provider.Secret) error {
	if m.err != nil {
		return m.err
	}
	m.written = append(m.written, secrets...)
	return nil
}

func TestPipeline_Run_Success(t *testing.T) {
	src := &mockSource{
		secrets: []provider.Secret{
			{Name: "DB_HOST", Value: "localhost", Type: "env"},
			{Name: "DB_PORT", Value: "5432", Type: "env"},
			{Name: "API_KEY", Value: "abc123", Type: "secret"},
		},
	}
	dst := &mockDestination{}
	proc, err := NewProcessor(config.ProcessConfig{})
	require.NoError(t, err)

	p := &Pipeline{src: src, dst: dst, proc: proc}
	err = p.Run(context.Background())
	require.NoError(t, err)
	assert.Len(t, dst.written, 3)
}

func TestPipeline_Run_SourceError(t *testing.T) {
	src := &mockSource{err: errors.New("source failure")}
	dst := &mockDestination{}
	proc, err := NewProcessor(config.ProcessConfig{})
	require.NoError(t, err)

	p := &Pipeline{src: src, dst: dst, proc: proc}
	err = p.Run(context.Background())
	require.ErrorContains(t, err, "read from source")
	assert.ErrorContains(t, err, "source failure")
}

func TestPipeline_Run_DestError(t *testing.T) {
	src := &mockSource{
		secrets: []provider.Secret{
			{Name: "KEY", Value: "val", Type: "env"},
		},
	}
	dst := &mockDestination{err: errors.New("dest failure")}
	proc, err := NewProcessor(config.ProcessConfig{})
	require.NoError(t, err)

	p := &Pipeline{src: src, dst: dst, proc: proc}
	err = p.Run(context.Background())
	require.ErrorContains(t, err, "write to destination")
	assert.ErrorContains(t, err, "dest failure")
}

func TestPipeline_Process(t *testing.T) {
	src := &mockSource{
		secrets: []provider.Secret{
			{Name: "DB_HOST", Value: "localhost", Type: "env"},
			{Name: "DB_PORT", Value: "5432", Type: "env"},
			{Name: "API_KEY", Value: "abc", Type: "secret"},
			{Name: "CACHE_HOST", Value: "redis", Type: "env"},
		},
	}
	dst := &mockDestination{}
	proc, err := NewProcessor(config.ProcessConfig{
		AddPrefix:    "PROD_",
		IncludeRegex: "^DB_",
		IncludeTypes: []string{"env"},
	})
	require.NoError(t, err)

	p := &Pipeline{src: src, dst: dst, proc: proc}
	err = p.Run(context.Background())
	require.NoError(t, err)
	require.Len(t, dst.written, 2)
	assert.Equal(t, "PROD_DB_HOST", dst.written[0].Name)
	assert.Equal(t, "PROD_DB_PORT", dst.written[1].Name)
}
