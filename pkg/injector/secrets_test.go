package injector

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
				data: []SecretData{
					{
						Name:  "key",
						value: []byte("value"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSecretFromKubeSecret(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSecretFromKubeSecret() = %v, want %v", got, tt.want)
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

// func TestSecret_CheckSecretFileHash(t *testing.T) {
// 	fs := memfs.NewFS()

// 	tests := []struct {
// 		name string
// 		s    *SecretData
// 		path string
// 		want bool
// 	}{
// 		{
// 			name: "Test secret file hash is the same",
// 			s: &SecretData{
// 				Name:  "key",
// 				value: []byte("value"),
// 			},
