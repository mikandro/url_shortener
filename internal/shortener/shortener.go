package shortener

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateShortCode(longUrl string) string {
	hash := sha256.Sum256([]byte(longUrl))

	hashString := hex.EncodeToString(hash[:])

	shortCode := hashString[:8]
	return shortCode
}
