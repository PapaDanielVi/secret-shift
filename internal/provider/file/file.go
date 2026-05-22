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

func init() {
	provider.Register(provider.Registration{
		Name: provider.File,
		SourceFactory: func(ctx context.Context, opts map[string]any) (provider.Source, error) {
			return New(
				getString(opts, "path"),
				getString(opts, "format"),
				getBool(opts, "encrypt"),
				getString(opts, "encrypt_key"),
			), nil
		},
		DestFactory: func(ctx context.Context, opts map[string]any) (provider.Destination, error) {
			return New(
				getString(opts, "path"),
				getString(opts, "format"),
				getBool(opts, "encrypt"),
				getString(opts, "encrypt_key"),
			), nil
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

func getBool(opts map[string]any, key string) bool {
	if v, ok := opts[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

type Provider struct {
	path       string
	format     string
	encrypt    bool
	encryptKey string
}

func New(path, format string, encrypt bool, encryptKey string) *Provider {
	if format == "" {
		format = "json"
	}
	return &Provider{
		path:       path,
		format:     format,
		encrypt:    encrypt,
		encryptKey: encryptKey,
	}
}

func (p *Provider) Read(ctx context.Context) ([]provider.Secret, error) {
	data, err := os.ReadFile(p.path)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", p.path, err)
	}

	if p.encrypt {
		data, err = p.decryptData(data)
		if err != nil {
			return nil, fmt.Errorf("decrypt file %s: %w", p.path, err)
		}
	}

	raw := make(map[string]string)
	switch p.format {
	case "yaml", "yml":
		if err := yaml.Unmarshal(data, &raw); err != nil {
			return nil, fmt.Errorf("unmarshal yaml: %w", err)
		}
	default:
		if err := json.Unmarshal(data, &raw); err != nil {
			return nil, fmt.Errorf("unmarshal json: %w", err)
		}
	}

	var result []provider.Secret
	for k, v := range raw {
		result = append(result, provider.Secret{
			Name:  k,
			Value: v,
			Type:  "env",
		})
	}
	return result, nil
}

func (p *Provider) Write(ctx context.Context, secrets []provider.Secret) error {
	if err := os.MkdirAll(filepath.Dir(p.path), 0750); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	data := make(map[string]string)
	for _, s := range secrets {
		data[s.Name] = s.Value
	}

	var output []byte
	var err error

	switch p.format {
	case "yaml", "yml":
		output, err = yaml.Marshal(data)
	default:
		output, err = json.MarshalIndent(data, "", "  ")
	}

	if err != nil {
		return fmt.Errorf("marshal output: %w", err)
	}

	if p.encrypt {
		output, err = p.encryptData(output)
		if err != nil {
			return fmt.Errorf("encrypt output: %w", err)
		}
	}

	if err := os.WriteFile(p.path, output, 0600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func (p *Provider) encryptData(plaintext []byte) ([]byte, error) {
	key := p.deriveKey()

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

func (p *Provider) decryptData(ciphertext []byte) ([]byte, error) {
	key := p.deriveKey()

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func (p *Provider) deriveKey() []byte {
	hash := sha256.Sum256([]byte(p.encryptKey))
	return hash[:]
}

var _ provider.Source = (*Provider)(nil)
var _ provider.Destination = (*Provider)(nil)
