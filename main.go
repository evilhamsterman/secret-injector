package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/spf13/afero"
)

func main() {
	data := "hello there"

	hash := sha256.New()
	hash.Write([]byte(data))
	hashValue := hash.Sum(nil)

	fmt.Printf("Data: %s\n Hash: %s\n", data, hex.EncodeToString(hashValue))

	fs := afero.NewOsFs()

	f, _ := afero.ReadFile(fs, "playground/test.txt")

	hash.Reset()
	hash.Write(f)
	hashValue = hash.Sum(nil)

	fmt.Printf("Data: %s\n Hash: %s\n", f, hex.EncodeToString(hashValue))

}
