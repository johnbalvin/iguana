package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

//GenerateChecksum256 generates a 256 checksum
func GenerateChecksum256(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
