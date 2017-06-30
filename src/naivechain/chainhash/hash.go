package chainhash

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashSize ...
const HashSize = 32

// Hash ...
type Hash [HashSize]byte

// HashH ...
func HashH(b []byte) Hash {
	hash := sha256.Sum256(b)
	return Hash(hash)
}

// DoubleHashH ...
func DoubleHashH(b []byte) Hash {
	first := sha256.Sum256(b)
	second := sha256.Sum256(first[:])
	return Hash(second)
}

// String ...
func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}
