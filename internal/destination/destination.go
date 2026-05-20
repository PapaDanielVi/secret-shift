package destination

import (
	"context"

	"github.com/PapaDanielVi/secret-shift/internal/source"
)

type Destination interface {
	Write(ctx context.Context, secrets []source.Secret) error
}
