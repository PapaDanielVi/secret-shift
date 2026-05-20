package provider

import "context"

type Secret struct {
	Name            string
	Value           string
	Type            string // "env" or "secret" or file
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
