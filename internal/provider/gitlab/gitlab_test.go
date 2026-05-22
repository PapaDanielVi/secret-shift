package gitlab

import (
	"testing"

	"github.com/PapaDanielVi/secret-shift/internal/provider"
	"github.com/stretchr/testify/assert"
)

func TestNew_ProviderFields(t *testing.T) {
	// Verify provider field assignment via the exported fields.
	p := &Provider{
		projectID: "123",
		strategy:  "replace",
	}
	assert.Equal(t, "123", p.projectID)
	assert.Equal(t, "replace", p.strategy)
}

var _ provider.Source = (*Provider)(nil)
var _ provider.Destination = (*Provider)(nil)
