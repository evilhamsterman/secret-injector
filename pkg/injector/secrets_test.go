package injector

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/hack-pad/hackpadfs"
	memfs "github.com/hack-pad/hackpadfs/mem"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	testSecretData = struct {
		hash  []byte
		value *SecretData
	}{
		hash: []byte("cd42404d52ad55ccfa9aca4adc828aa5800ad9d385a0671fbcbf724118320619"),
		value: &SecretData{
			Name:  "key",
			value: []byte("value"),
		},
	}
)

func TestNewSecretFromKubeSecret(t *testing.T) {
	tests := []struct {
		name string
		s    *v1.Secret
		want *Secret
	}{
		{
			name: "Test secret creation",
			s: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					Annotations: map[string]string{
						"secret-injector/path": "/tmp",
					},
				},
				Data: map[string][]byte{
					"key": []byte("value"),
				},
			},
			want: &Secret{
				Name:      "test",
				Namespace: "default",
				Path:      "/tmp",
				data:      []SecretData{*testSecretData.value},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSecretFromKubeSecret(tt.s)
			if got.Name != tt.want.Name {
				t.Errorf("NewSecretFromKubeSecret() = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Namespace != tt.want.Namespace {
				t.Errorf("NewSecretFromKubeSecret() = %v, want %v", got.Namespace, tt.want.Namespace)
			}
			if got.Path != tt.want.Path {
				t.Errorf("NewSecretFromKubeSecret() = %v, want %v", got.Path, tt.want.Path)
			}
			if got.data[0].Name != tt.want.data[0].Name {
				t.Errorf("NewSecretFromKubeSecret() = %v, want %v", got.data[0].Name, tt.want.data[0].Name)
			}
			if !bytes.Equal(got.data[0].value, tt.want.data[0].value) {
				t.Errorf("NewSecretFromKubeSecret() = %v, want %v", got.data[0].value, tt.want.data[0].value)
			}
		})
	}
}

func TestSecretData_GetSecretDataHash(t *testing.T) {
	testString := []byte("value")
	testHash, _ := hex.DecodeString("cd42404d52ad55ccfa9aca4adc828aa5800ad9d385a0671fbcbf724118320619")
	tests := []struct {
		name string
		s    *SecretData
		want []byte
	}{
		{
			name: "Test secret data hash",
			s: &SecretData{
				Name:  "key",
				value: testString,
			},
			want: testHash,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.GetSecretDataHash()
			if !bytes.Equal(got, tt.want) {
				t.Errorf("GetSecretDataHash() = %x, want %x", got, tt.want)
			}
		})
	}
}

func TestSecret_CheckSecretFileHash(t *testing.T) {
	path := "test"
	fs, _ := memfs.NewFS()
	f, _ := hackpadfs.OpenFile(fs, path, hackpadfs.FlagCreate|hackpadfs.FlagReadWrite, 0644)
	hackpadfs.WriteFile(f, []byte("value"))
	f.Close()

	tests := []struct {
		name string
		s    *SecretData
		want bool
	}{
		{
			name: "Test secret file hash is the same",
			s:    testSecretData.value,
			want: true,
		},
		{
			name: "Test secret file hash is not the same",
			s: &SecretData{
				Name:  "key",
				value: []byte("not the same"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.CheckSecretFileHash(fs, path)
			if got != tt.want {
				t.Errorf("CheckSecretFileHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_WriteSecretData(t *testing.T) {
	fs, _ := memfs.NewFS()
	path := "test"

	tests := []struct {
		name     string
		s        *SecretData
		fileData []byte
	}{
		{
			name:     "Test write secret data",
			s:        testSecretData.value,
			fileData: testSecretData.value.value,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.s.WriteSecretData(fs, path)
			if err != nil {
				t.Fatalf("WriteSecretData() error = %v", err)
			}
			f, _ := hackpadfs.ReadFile(fs, path)
			if !bytes.Equal(f, tt.fileData) {
				t.Errorf("WriteSecretData() = %s, want %s", f, tt.fileData)
			}
		})
	}
}
