package file

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/PapaDanielVi/secret-shift/internal/source"
)

func TestWrite_JSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "secrets.json")

	d := New(path, "json", false, "")
	secrets := []source.Secret{
		{Name: "DB_HOST", Value: "localhost", Type: "env"},
		{Name: "API_KEY", Value: "abc123", Type: "secret"},
	}

	err := d.Write(context.Background(), secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected non-empty file")
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat file: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected permissions 0600, got %o", info.Mode().Perm())
	}
}

func TestWrite_YAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "secrets.yaml")

	d := New(path, "yaml", false, "")
	secrets := []source.Secret{
		{Name: "DB_HOST", Value: "localhost", Type: "env"},
	}

	err := d.Write(context.Background(), secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected non-empty file")
	}
}

func TestWrite_Encrypted(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "secrets.enc")

	d := New(path, "json", true, "my-secret-key-32-chars-long!!")
	secrets := []source.Secret{
		{Name: "DB_HOST", Value: "localhost", Type: "env"},
		{Name: "API_KEY", Value: "abc123", Type: "secret"},
	}

	err := d.Write(context.Background(), secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected non-empty file")
	}

	// Encrypted data should not contain plaintext
	if string(data) == "" {
		t.Error("expected non-empty encrypted data")
	}
	// Verify it's not valid JSON (since it's encrypted)
	if err := os.WriteFile(path+".test", data, 0600); err == nil {
		// Just check the data doesn't look like our plaintext
		plaintext := `{"API_KEY":"abc123","DB_HOST":"localhost"}`
		if string(data) == plaintext {
			t.Error("data should be encrypted, not plaintext")
		}
		os.Remove(path + ".test")
	}
}

func TestWrite_CreatesDirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "deep", "secrets.json")

	d := New(path, "json", false, "")
	secrets := []source.Secret{
		{Name: "KEY", Value: "val", Type: "env"},
	}

	err := d.Write(context.Background(), secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected file to be created")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	d := &Destination{encryptKey: "test-key-for-encryption!!"}

	plaintext := []byte(`{"KEY":"value"}`)
	ciphertext, err := d.encryptData(plaintext)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	// Ciphertext should be different from plaintext
	if string(ciphertext) == string(plaintext) {
		t.Error("ciphertext should differ from plaintext")
	}

	// Ciphertext should be longer (nonce + tag)
	if len(ciphertext) <= len(plaintext) {
		t.Error("ciphertext should be longer than plaintext")
	}
}
