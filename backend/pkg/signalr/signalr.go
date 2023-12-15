package signalr

import (
	"crypto/rand"
	"encoding/base64"
)

// Use a good random number generator
func GenerateRandomString(length int) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}
