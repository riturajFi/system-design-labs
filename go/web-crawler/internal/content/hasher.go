package content

import (
	"crypto/sha256"
	"encoding/hex"
)

// Hash returns a stable hex-encoded SHA-256 digest for content bytes.
func Hash(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
