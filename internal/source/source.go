package source

import (
	"context"

	"github.com/PapaDanielVi/secret-shift/internal/provider"
)

type Secret = provider.Secret

type Source interface {
	Read(ctx context.Context) ([]Secret, error)
}
