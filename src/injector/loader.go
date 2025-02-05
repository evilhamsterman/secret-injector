package injector

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/afero"
)

func GetSecrets(fs afero.Fs, path string) (*SecretsList, error) {
	// make sure the file exists
	exists, err := afero.Exists(fs, path)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("file %s does not exist", path)
	}

	// read the file
	file, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secrets file: %w", err)
	}

	// decode the file
	secrets := &SecretsList{}
	if err = json.Unmarshal([]byte(file), secrets); err != nil {
		return nil, fmt.Errorf("failed to load secrets file: %w", err)
	}

	return secrets, nil
}
