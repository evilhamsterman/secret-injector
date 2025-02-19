package injector

import (
	"crypto/sha256"

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
	hash  []byte
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

func hashSecretData(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

// GetSecretDataHash returns the sha256 hash of the secret
func (s *SecretData) GetSecretDataHash() []byte {
	if s.hash != nil {
		return s.hash
	}
	s.hash = hashSecretData(s.value)
	return s.hash
}

// CheckSecretFileHash checks if the hash of the file is the same as the secret
func (s *SecretData) CheckSecretFileHash(fs hackpadfs.FS, path string) bool {
	return false
}
