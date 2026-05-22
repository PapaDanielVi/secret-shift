package etcd

import (
	"context"
	"fmt"
	"strings"

	"github.com/PapaDanielVi/secret-shift/internal/provider"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func init() {
	provider.Register(provider.Registration{
		Name: provider.Etcd,
		SourceFactory: func(_ context.Context, opts map[string]any) (provider.Source, error) {
			endpoints := getStringSlice(opts, "etcd_endpoints")
			prefix := getString(opts, "etcd_prefix")
			username := getString(opts, "etcd_username")
			password := getString(opts, "etcd_password")
			return New(endpoints, prefix, username, password)
		},
		DestFactory: func(_ context.Context, opts map[string]any) (provider.Destination, error) {
			endpoints := getStringSlice(opts, "etcd_endpoints")
			prefix := getString(opts, "etcd_prefix")
			username := getString(opts, "etcd_username")
			password := getString(opts, "etcd_password")
			return New(endpoints, prefix, username, password)
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

func getStringSlice(opts map[string]any, key string) []string {
	if v, ok := opts[key]; ok {
		if arr, ok := v.([]any); ok {
			result := make([]string, 0, len(arr))
			for _, item := range arr {
				if s, ok := item.(string); ok {
					result = append(result, s)
				}
			}
			return result
		}
	}
	return nil
}

type Provider struct {
	client *clientv3.Client
	prefix string
}

func New(endpoints []string, prefix, username, password string) (*Provider, error) {
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

	return &Provider{
		client: client,
		prefix: prefix,
	}, nil
}

func (p *Provider) Read(ctx context.Context) ([]provider.Secret, error) {
	resp, err := p.client.Get(ctx, p.prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("read from etcd with prefix %s: %w", p.prefix, err)
	}

	var result []provider.Secret
	for _, kv := range resp.Kvs {
		name := strings.TrimPrefix(string(kv.Key), p.prefix)
		if name == "" {
			name = string(kv.Key)
		}
		result = append(result, provider.Secret{
			Name:  name,
			Value: string(kv.Value),
			Type:  "env",
		})
	}

	return result, nil
}

func (p *Provider) Write(ctx context.Context, secrets []provider.Secret) error {
	for _, s := range secrets {
		key := p.prefix + s.Name
		if p.prefix != "/" && !strings.HasSuffix(p.prefix, "/") {
			key = p.prefix + "/" + s.Name
		}
		_, err := p.client.Put(ctx, key, s.Value)
		if err != nil {
			return fmt.Errorf("write etcd key %s: %w", key, err)
		}
	}

	return nil
}

func (p *Provider) Close() error {
	return p.client.Close()
}

var _ provider.Source = (*Provider)(nil)
var _ provider.Destination = (*Provider)(nil)
