package base62

import (
	"urlshortener/internal/domain"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Encode converts a numeric ID to a Base62 short URL.
// This function is pure and deterministic.
func Encode(id domain.ID) domain.ShortURL {
	if id == 0 {
		return domain.ShortURL(string(alphabet[0]))
	}

	n := int64(id)
	var encoded []byte

	for n > 0 {
		remainder := n % 62
		encoded = append(encoded, alphabet[remainder])
		n = n / 62
	}

	// reverse result
	for i, j := 0, len(encoded)-1; i < j; i, j = i+1, j-1 {
		encoded[i], encoded[j] = encoded[j], encoded[i]
	}

	return domain.ShortURL(encoded)
}
