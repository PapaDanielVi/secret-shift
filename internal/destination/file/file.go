package file

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/PapaDanielVi/secret-shift/internal/source"
	"gopkg.in/yaml.v3"
)

type Destination struct {
	path   string
	format string
}

func New(path, format string) *Destination {
	return &Destination{
		path:   path,
		format: format,
	}
}

func (d *Destination) Write(ctx context.Context, secrets []source.Secret) error {
	if err := os.MkdirAll(filepath.Dir(d.path), 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	data := make(map[string]string)
	for _, s := range secrets {
		data[s.Name] = s.Value
	}

	var output []byte
	var err error

	switch d.format {
	case "yaml", "yml":
		output, err = yaml.Marshal(data)
	default:
		output, err = json.MarshalIndent(data, "", "  ")
	}

	if err != nil {
		return fmt.Errorf("marshal output: %w", err)
	}

	if err := os.WriteFile(d.path, output, 0600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

var _ interface {
	Write(ctx context.Context, secrets []source.Secret) error
} = (*Destination)(nil)
