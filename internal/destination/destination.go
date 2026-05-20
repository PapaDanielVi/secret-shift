package destination

import (
	"context"

	"github.com/PapaDanielVi/secret-shift/internal/provider"
)

type Destination interface {
	Write(ctx context.Context, secrets []provider.Secret) error
}
