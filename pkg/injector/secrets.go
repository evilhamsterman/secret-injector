package injector

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"path/filepath"

	"github.com/hack-pad/hackpadfs"
	v1 "k8s.io/api/core/v1"
)

type Secret struct {
	Name      string
	Namespace string
	Path      string
	data      []SecretData
}

type SecretData struct {
	Name  string
	value []byte
}

// NewSecretFromKubeSecret creates a new Secret from a kubernetes secret
func NewSecretFromKubeSecret(s *v1.Secret) *Secret {
	secret := &Secret{
		Name:      s.Name,
		Namespace: s.Namespace,
		Path:      s.Annotations["secret-injector/path"],
	}

	for k, v := range s.Data {
		secret.data = append(secret.data, SecretData{
			Name:  k,
			value: v,
		})
	}

	return secret
}

func hashData(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

// GetSecretDataHash returns the sha256 hash of the secret
func (s *SecretData) GetSecretDataHash() []byte {
	return hashData(s.value)
}

// CheckSecretFileHash checks if the hash of the file is the same as the secret
func (s *SecretData) CheckSecretFileHash(fs hackpadfs.FS, path string) bool {
	f, err := hackpadfs.ReadFile(fs, path)
	if err != nil {
		return false
	}
	return bytes.Equal(s.GetSecretDataHash(), hashData(f))
}

// WriteSecretData writes the secret data to the filesystem
func (s *SecretData) WriteSecretData(fs hackpadfs.FS, path string) error {
	err := hackpadfs.MkdirAll(fs, filepath.Dir(path), 0o755)
	if err != nil {
		return fmt.Errorf("Failed to create directory: %w", err)
	}
	f, err := hackpadfs.OpenFile(
		fs,
		path,
		hackpadfs.FlagCreate|hackpadfs.FlagWriteOnly|hackpadfs.FlagTruncate,
		0o644,
	)
	if err != nil {
		return fmt.Errorf("Failed to create file: %w", err)
	}
	defer f.Close()

	_, err = hackpadfs.WriteFile(f, s.value)
	if err != nil {
		return fmt.Errorf("Failed to write file: %w", err)
	}

	return nil
}

// UpdateSecretData updates the secret data if the hash is different
func (s *SecretData) UpdateSecretData(fs hackpadfs.FS, path string) error {
	if s.CheckSecretFileHash(fs, path) {
		return nil
	}
	return s.WriteSecretData(fs, path)
}
