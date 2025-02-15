package injector

type Secret struct {
	Name      string   `json:"name"`
	Namespace string   `json:"namespace,omitempty"`
	Path      string   `json:"path"`
	Keys      []string `json:"keys,omitempty"`
}

type SecretsList []Secret
