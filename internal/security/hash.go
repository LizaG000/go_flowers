package security

import (
	"crypto/sha256"
	"encoding/hex"
)

func CalculatePayloadHash(payload []byte) string {
	hash := sha256.Sum256(payload)

	return hex.EncodeToString(hash[:])
}
