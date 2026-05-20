package etcd

import (
	"context"
	"fmt"
	"strings"

	"github.com/PapaDanielVi/secret-shift/internal/destination"
	"github.com/PapaDanielVi/secret-shift/internal/source"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Destination struct {
	client *clientv3.Client
	prefix string
}

func New(endpoints []string, prefix, username, password string) (*Destination, error) {
	cfg := clientv3.Config{
		Endpoints: endpoints,
	}
	if username != "" {
		cfg.Username = username
		cfg.Password = password
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("create etcd client: %w", err)
	}

	if prefix == "" {
		prefix = "/"
	}

	return &Destination{
		client: client,
		prefix: prefix,
	}, nil
}

func (d *Destination) Write(ctx context.Context, secrets []source.Secret) error {
	for _, s := range secrets {
		key := d.prefix + s.Name
		if d.prefix != "/" && !strings.HasSuffix(d.prefix, "/") {
			key = d.prefix + "/" + s.Name
		}
		_, err := d.client.Put(ctx, key, s.Value)
		if err != nil {
			return fmt.Errorf("write etcd key %s: %w", key, err)
		}
	}

	return nil
}

func (d *Destination) Close() error {
	return d.client.Close()
}

var _ destination.Destination = (*Destination)(nil)
