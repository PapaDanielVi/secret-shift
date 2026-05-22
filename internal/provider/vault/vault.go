// Package vault implements the provider interface for HashiCorp Vault.
package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/PapaDanielVi/secret-shift/internal/provider"
	vaultapi "github.com/hashicorp/vault/api"
)

func init() {
	provider.Register(provider.Registration{
		Name: provider.Vault,
		SourceFactory: func(_ context.Context, opts map[string]any) (provider.Source, error) {
			token := getString(opts, "token")
			address := getString(opts, "vault_address")
			path := getString(opts, "vault_path")
			mount := getString(opts, "vault_mount")
			return New(token, address, path, mount)
		},
		DestFactory: func(_ context.Context, opts map[string]any) (provider.Destination, error) {
			token := getString(opts, "token")
			address := getString(opts, "vault_address")
			path := getString(opts, "vault_path")
			mount := getString(opts, "vault_mount")
			return New(token, address, path, mount)
		},
	})
}

func getString(opts map[string]any, key string) string {
	if v, ok := opts[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

type Provider struct {
	client *vaultapi.Client
	path   string
	mount  string
}

func New(token, address, path, mount string) (*Provider, error) {
	config := vaultapi.DefaultConfig()
	if address != "" {
		config.Address = address
	}

	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create vault client: %w", err)
	}

	client.SetToken(token)

	if mount == "" {
		mount = "secret"
	}

	return &Provider{
		client: client,
		path:   path,
		mount:  mount,
	}, nil
}

func (p *Provider) Read(ctx context.Context) ([]provider.Secret, error) {
	secret, err := p.client.KVv2(p.mount).Get(ctx, p.path)
	if err != nil {
		return nil, fmt.Errorf("read vault secret at %s/%s: %w", p.mount, p.path, err)
	}

	var result []provider.Secret
	for k, v := range secret.Data {
		strVal, err := stringifyVaultValue(v)
		if err != nil {
			return nil, fmt.Errorf("stringify value for key %s: %w", k, err)
		}
		result = append(result, provider.Secret{
			Name:  k,
			Value: strVal,
			Type:  "env",
		})
	}

	return result, nil
}

func stringifyVaultValue(v any) (string, error) {
	switch val := v.(type) {
	case string:
		return val, nil
	case []byte:
		return string(val), nil
	case json.Number:
		return val.String(), nil
	case nil:
		return "", nil
	default:
		slog.Warn("unexpected vault value type", "type", fmt.Sprintf("%T", v))
		b, err := json.Marshal(v)
		if err != nil {
			return "", fmt.Errorf("marshal vault value: %w", err)
		}
		return string(b), nil
	}
}

func (p *Provider) Write(ctx context.Context, secrets []provider.Secret) error {
	data := make(map[string]any)
	for _, s := range secrets {
		data[s.Name] = s.Value
	}

	_, err := p.client.KVv2(p.mount).Put(ctx, p.path, data)
	if err != nil {
		return fmt.Errorf("write vault secret at %s/%s: %w", p.mount, p.path, err)
	}

	return nil
}

var _ provider.Source = (*Provider)(nil)
var _ provider.Destination = (*Provider)(nil)
