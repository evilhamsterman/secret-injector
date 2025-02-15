package injector

import (
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

func TestGetSecrets(t *testing.T) {
	appFS := afero.NewMemMapFs()
	path := "/test/secrets.json"

	cases := []struct {
		name        string
		json        string
		secretsList *SecretsList
	}{
		{
			name: "secret with all fields",
			json: `[{
				"name":"test-secret",
				"namespace":"test-namespace",
				"path":"/test/secret",
				"keys":["key1","key2"]
			}]`,
			secretsList: &SecretsList{
				{
					Name:      "test-secret",
					Namespace: "test-namespace",
					Path:      "/test/secret",
					Keys:      []string{"key1", "key2"},
				},
			},
		},
		{
			name: "secret with no namespace",
			json: `[{
				"name":"test-secret",
				"path":"/test/secret",
				"keys":["key1","key2"]
			}]`,
			secretsList: &SecretsList{
				{
					Name: "test-secret",
					Path: "/test/secret",
					Keys: []string{"key1", "key2"},
				},
			},
		},
		{
			name: "secret with no keys",
			json: `[{
				"name":"test-secret",
				"namespace":"test-namespace",
				"path":"/test/secret"
			}]`,
			secretsList: &SecretsList{
				{
					Name:      "test-secret",
					Namespace: "test-namespace",
					Path:      "/test/secret",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// write the json to the file
			err := afero.WriteFile(appFS, path, []byte(c.json), 0644)
			if err != nil {
				t.Fatalf("failed to write file: %v", err)
			}

			// get the secrets
			secrets, err := GetSecrets(appFS, path)
			if err != nil {
				t.Fatalf("failed to get secrets: %v", err)
			}

			// ensure we have the right number of secrets
			if len(*secrets) != 1 {
				t.Fatalf("expected 1 secrets, got %d", len(*secrets))
			}

			// ensure the secret is correct
			if !reflect.DeepEqual(secrets, c.secretsList) {
				t.Fatalf("expected %v, got %v", c.secretsList, secrets)
			}
		})
	}
}
