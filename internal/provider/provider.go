// Package provider defines the interfaces and types for secret providers.
package provider

import "context"

// Type identifies a secret provider.
type Type string

const (
	GitHub     Type = "github"
	GitLab     Type = "gitlab"
	Vault      Type = "vault"
	Etcd       Type = "etcd"
	Kubernetes Type = "kubernetes"
	File       Type = "file"
)

type Secret struct {
	Name            string
	Value           string
	Type            string // "env" or "secret" or "file"
	OtherAttributes map[string]string
}

type Source interface {
	Read(ctx context.Context) ([]Secret, error)
}

type Destination interface {
	Write(ctx context.Context, secrets []Secret) error
}

type Provider interface {
	Source
	Destination
}
