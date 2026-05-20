package github

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

func encryptSecret(publicKey, secretValue string) (string, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", fmt.Errorf("decode public key: %w", err)
	}

	pubKeyInterface, err := x509.ParsePKIXPublicKey(decodedKey)
	if err != nil {
		return "", fmt.Errorf("parse public key: %w", err)
	}

	pubKey, ok := pubKeyInterface.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("public key is not RSA")
	}

	encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, []byte(secretValue), nil)
	if err != nil {
		return "", fmt.Errorf("encrypt: %w", err)
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}
