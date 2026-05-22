package file

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/PapaDanielVi/secret-shift/internal/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileRead_JSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "secrets.json")

	err := os.WriteFile(path, []byte(`{"KEY1":"val1","KEY2":"val2"}`), 0600)
	require.NoError(t, err)

	p := New(path, "json", false, "")
	secrets, err := p.Read(context.Background())
	require.NoError(t, err)
	assert.Len(t, secrets, 2)

	names := make(map[string]string)
	for _, s := range secrets {
		names[s.Name] = s.Value
	}
	assert.Equal(t, "val1", names["KEY1"])
	assert.Equal(t, "val2", names["KEY2"])
}

func TestFileRead_YAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "secrets.yaml")

	err := os.WriteFile(path, []byte("KEY1: val1\nKEY2: val2\n"), 0600)
	require.NoError(t, err)

	p := New(path, "yaml", false, "")
	secrets, err := p.Read(context.Background())
	require.NoError(t, err)
	assert.Len(t, secrets, 2)
}

func TestFileRead_Encrypted(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "secrets.json")

	p := New(path, "json", true, "test-key-12345")
	data := map[string]string{"SECRET": "hello"}
	raw, err := json.Marshal(data)
	require.NoError(t, err)
	enc, err := p.encryptData(raw)
	require.NoError(t, err)
	err = os.WriteFile(path, enc, 0600)
	require.NoError(t, err)

	secrets, err := p.Read(context.Background())
	require.NoError(t, err)
	names := make(map[string]string)
	for _, s := range secrets {
		names[s.Name] = s.Value
	}
	assert.Equal(t, "hello", names["SECRET"])
}

func TestFileRead_MissingFile(t *testing.T) {
	p := New("/nonexistent/path.json", "json", false, "")
	_, err := p.Read(context.Background())
	assert.Error(t, err)
}

func TestFileWrite_JSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output", "secrets.json")

	p := New(path, "json", false, "")
	err := p.Write(context.Background(), []provider.Secret{
		{Name: "KEY1", Value: "val1"},
		{Name: "KEY2", Value: "val2"},
	})
	require.NoError(t, err)

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(content), `"KEY1": "val1"`)
}

func TestFileWrite_YAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output.yaml")

	p := New(path, "yaml", false, "")
	err := p.Write(context.Background(), []provider.Secret{
		{Name: "KEY1", Value: "val1"},
	})
	require.NoError(t, err)

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(content), "KEY1: val1")
}

func TestFileWrite_Encrypted(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "encrypted.json")

	p := New(path, "json", true, "my-secret-key")
	err := p.Write(context.Background(), []provider.Secret{
		{Name: "SECRET", Value: "sensitive"},
	})
	require.NoError(t, err)

	// Verify it's not plaintext
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.NotContains(t, string(content), "sensitive")

	// Verify we can read it back
	secrets, err := p.Read(context.Background())
	require.NoError(t, err)
	names := make(map[string]string)
	for _, s := range secrets {
		names[s.Name] = s.Value
	}
	assert.Equal(t, "sensitive", names["SECRET"])
}

func TestFileWrite_CreatesDirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a", "b", "c", "secrets.json")

	p := New(path, "json", false, "")
	err := p.Write(context.Background(), []provider.Secret{
		{Name: "K", Value: "v"},
	})
	require.NoError(t, err)
	_, err = os.Stat(path)
	assert.NoError(t, err)
}

func TestFileProvider_Interface(_ *testing.T) {
	var _ provider.Source = (*Provider)(nil)
	var _ provider.Destination = (*Provider)(nil)
}
