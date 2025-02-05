package injector

type SecretPathErrorNotExists struct {
	Path string
}

func (e *SecretPathErrorNotExists) Error() string {
	return "file " + e.Path + " does not exist"
}

type SecretHashError struct {
	ExpectedHash string
	ActualHash   string
}

func (e *SecretHashError) Error() string {
	return "hash mismatch: expected " + e.ExpectedHash + " but got " + e.ActualHash
}
