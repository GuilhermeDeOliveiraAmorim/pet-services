package entities

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func HashToken(token string) string {
	normalized := strings.TrimSpace(token)
	if normalized == "" {
		return ""
	}

	sum := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(sum[:])
}
