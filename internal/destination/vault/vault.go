package vault

import (
	"context"
	"fmt"

	"github.com/PapaDanielVi/secret-shift/internal/destination"
	"github.com/PapaDanielVi/secret-shift/internal/source"
	vault "github.com/hashicorp/vault/api"
)

type Destination struct {
	client  *vault.Client
	path    string
	mount   string
}

func New(token, address, path, mount string) (*Destination, error) {
	config := vault.DefaultConfig()
	if address != "" {
		config.Address = address
	}

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create vault client: %w", err)
	}

	client.SetToken(token)

	if mount == "" {
		mount = "secret"
	}

	return &Destination{
		client: client,
		path:   path,
		mount:  mount,
	}, nil
}

func (d *Destination) Write(ctx context.Context, secrets []source.Secret) error {
	data := make(map[string]interface{})
	for _, s := range secrets {
		data[s.Name] = s.Value
	}

	_, err := d.client.KVv2(d.mount).Put(ctx, d.path, data)
	if err != nil {
		return fmt.Errorf("write vault secret at %s/%s: %w", d.mount, d.path, err)
	}

	return nil
}

var _ destination.Destination = (*Destination)(nil)
