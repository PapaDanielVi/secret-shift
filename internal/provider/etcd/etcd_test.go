package etcd

import (
	"context"
	"testing"

	"github.com/PapaDanielVi/secret-shift/internal/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

func TestRegisteredFactoryAcceptsStringEndpoints(t *testing.T) {
	t.Parallel()

	registration, err := provider.Get(provider.Etcd)
	require.NoError(t, err)

	source, err := registration.SourceFactory(context.Background(), map[string]any{
		"etcd_endpoints": []string{"http://127.0.0.1:2379"},
	})
	require.NoError(t, err)

	etcdProvider, ok := source.(*Provider)
	require.True(t, ok)
	assert.Equal(t, []string{"http://127.0.0.1:2379"}, etcdProvider.client.Endpoints())
	require.NoError(t, etcdProvider.Close())
}

func TestSecretsFromKVsUsesStoredValues(t *testing.T) {
	t.Parallel()

	secrets := secretsFromKVs([]*mvccpb.KeyValue{
		{Key: []byte("/app/database-url"), Value: []byte("postgres://database")},
	}, "/app/")

	require.Len(t, secrets, 1)
	assert.Equal(t, "database-url", secrets[0].Name)
	assert.Equal(t, "postgres://database", secrets[0].Value)
	assert.Equal(t, "env", secrets[0].Type)
}
