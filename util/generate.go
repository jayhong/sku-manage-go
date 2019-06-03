package util

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateID(n int) (string, error) {
	r := make([]byte, n)
	if _, err := rand.Read(r); err != nil {
		return "", err
	}

	return hex.EncodeToString(r), nil
}
