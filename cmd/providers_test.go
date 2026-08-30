package cmd

import (
	"testing"

	"github.com/PapaDanielVi/secret-shift/internal/provider"
	"github.com/stretchr/testify/require"
)

func TestProvidersAreRegistered(t *testing.T) {
	t.Parallel()

	providerTypes := []provider.Type{
		provider.GitHub,
		provider.GitLab,
		provider.Vault,
		provider.Etcd,
		provider.Kubernetes,
		provider.File,
	}
	for _, providerType := range providerTypes {
		_, err := provider.Get(providerType)
		require.NoError(t, err, "provider %q is not registered", providerType)
	}
}
