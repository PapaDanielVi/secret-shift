package vault

import (
	"context"
	"fmt"

	"github.com/PapaDanielVi/secret-shift/internal/source"
	vault "github.com/hashicorp/vault/api"
)

type Source struct {
	client  *vault.Client
	path    string
	mount   string
}

func New(token, address, path, mount string) (*Source, error) {
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

	return &Source{
		client: client,
		path:   path,
		mount:  mount,
	}, nil
}

func (s *Source) Read(ctx context.Context) ([]source.Secret, error) {
	secret, err := s.client.KVv2(s.mount).Get(ctx, s.path)
	if err != nil {
		return nil, fmt.Errorf("read vault secret at %s/%s: %w", s.mount, s.path, err)
	}

	var result []source.Secret
	for k, v := range secret.Data {
		strVal := fmt.Sprintf("%v", v)
		result = append(result, source.Secret{
			Name:  k,
			Value: strVal,
			Type:  "env",
		})
	}

	return result, nil
}
