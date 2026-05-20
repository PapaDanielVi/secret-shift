package etcd

import (
	"context"
	"fmt"
	"strings"

	"github.com/PapaDanielVi/secret-shift/internal/source"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Source struct {
	client *clientv3.Client
	prefix string
}

func New(endpoints []string, prefix, username, password string) (*Source, error) {
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

	return &Source{
		client: client,
		prefix: prefix,
	}, nil
}

func (s *Source) Read(ctx context.Context) ([]source.Secret, error) {
	resp, err := s.client.Get(ctx, s.prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("read from etcd with prefix %s: %w", s.prefix, err)
	}

	var result []source.Secret
	for _, kv := range resp.Kvs {
		name := strings.TrimPrefix(string(kv.Key), s.prefix)
		if name == "" {
			name = string(kv.Key)
		}
		result = append(result, source.Secret{
			Name:  name,
			Value: string(kv.Value),
			Type:  "env",
		})
	}

	return result, nil
}

func (s *Source) Close() error {
	return s.client.Close()
}
