package vault

import (
	"testing"

	"github.com/PapaDanielVi/secret-shift/internal/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringifyVaultValue_String(t *testing.T) {
	result, err := stringifyVaultValue("hello")
	require.NoError(t, err)
	assert.Equal(t, "hello", result)
}

func TestStringifyVaultValue_Bytes(t *testing.T) {
	result, err := stringifyVaultValue([]byte("world"))
	require.NoError(t, err)
	assert.Equal(t, "world", result)
}

func TestStringifyVaultValue_Nil(t *testing.T) {
	result, err := stringifyVaultValue(nil)
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestStringifyVaultValue_Struct(t *testing.T) {
	result, err := stringifyVaultValue(map[string]any{"nested": "value"})
	require.NoError(t, err)
	assert.Contains(t, result, "nested")
	assert.Contains(t, result, "value")
}

var _ provider.Source = (*Provider)(nil)
var _ provider.Destination = (*Provider)(nil)
