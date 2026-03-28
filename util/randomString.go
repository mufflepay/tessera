package util

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomString(length int) (string, error) {
	key := make([]byte, length) // 256-bit key
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}
