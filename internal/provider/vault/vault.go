package vault

import (
	"context"
	"fmt"

	"github.com/PapaDanielVi/secret-shift/internal/provider"
	vaultapi "github.com/hashicorp/vault/api"
)

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
		strVal := fmt.Sprintf("%v", v)
		result = append(result, provider.Secret{
			Name:  k,
			Value: strVal,
			Type:  "env",
		})
	}

	return result, nil
}

func (p *Provider) Write(ctx context.Context, secrets []provider.Secret) error {
	data := make(map[string]interface{})
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
