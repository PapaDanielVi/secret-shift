package provider_test

import (
	"context"
	"testing"

	"github.com/PapaDanielVi/secret-shift/internal/provider"
	_ "github.com/PapaDanielVi/secret-shift/internal/provider/etcd"
	_ "github.com/PapaDanielVi/secret-shift/internal/provider/file"
	_ "github.com/PapaDanielVi/secret-shift/internal/provider/github"
	_ "github.com/PapaDanielVi/secret-shift/internal/provider/gitlab"
	_ "github.com/PapaDanielVi/secret-shift/internal/provider/kubernetes"
	_ "github.com/PapaDanielVi/secret-shift/internal/provider/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry_AllProvidersRegistered(t *testing.T) {
	expected := []provider.Type{provider.GitHub, provider.GitLab, provider.Vault, provider.Etcd, provider.Kubernetes, provider.File}
	for _, typ := range expected {
		reg, err := provider.Get(typ)
		require.NoError(t, err, "provider %s should be registered", typ)
		assert.Equal(t, typ, reg.Name)
	}
}

func TestRegistry_GetUnknown(t *testing.T) {
	_, err := provider.Get(provider.Type("nonexistent"))
	assert.Error(t, err)
}

func TestRegistry_CreateSource_Unknown(t *testing.T) {
	_, err := provider.CreateSource(context.Background(), provider.Type("nonexistent"), nil)
	assert.Error(t, err)
}

func TestRegistry_CreateDestination_Unknown(t *testing.T) {
	_, err := provider.CreateDestination(context.Background(), provider.Type("nonexistent"), nil)
	assert.Error(t, err)
}

func TestRegistry_RegisteredTypes(t *testing.T) {
	types := provider.RegisteredTypes()
	assert.GreaterOrEqual(t, len(types), 6)
	typeSet := make(map[provider.Type]bool)
	for _, t := range types {
		typeSet[t] = true
	}
	assert.True(t, typeSet[provider.GitHub])
	assert.True(t, typeSet[provider.GitLab])
	assert.True(t, typeSet[provider.Vault])
	assert.True(t, typeSet[provider.Etcd])
	assert.True(t, typeSet[provider.Kubernetes])
	assert.True(t, typeSet[provider.File])
}
