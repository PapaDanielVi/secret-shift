package file

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/PapaDanielVi/secret-shift/internal/provider"
	"gopkg.in/yaml.v3"
)

type Destination struct {
	path       string
	format     string
	encrypt    bool
	encryptKey string
}

func New(path, format string, encrypt bool, encryptKey string) *Destination {
	return &Destination{
		path:       path,
		format:     format,
		encrypt:    encrypt,
		encryptKey: encryptKey,
	}
}

func (d *Destination) Write(ctx context.Context, secrets []provider.Secret) error {
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

	if d.encrypt {
		output, err = d.encryptData(output)
		if err != nil {
			return fmt.Errorf("encrypt output: %w", err)
		}
	}

	if err := os.WriteFile(d.path, output, 0600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func (d *Destination) encryptData(plaintext []byte) ([]byte, error) {
	key := d.deriveKey()

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func (d *Destination) deriveKey() []byte {
	hash := sha256.Sum256([]byte(d.encryptKey))
	return hash[:]
}
