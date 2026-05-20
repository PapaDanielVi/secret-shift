package source

import "context"

type Secret struct {
	Name  string
	Value string
	Type  string // "env" or "secret"
}

type Source interface {
	Read(ctx context.Context) ([]Secret, error)
}
