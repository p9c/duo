package hash160

import (
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)

// Sum returns the result of sha256 and then ripemd160 on a message
func Sum(message *[]byte) *[]byte {
	h := sha256.Sum256(*message)
	o := ripemd160.New().Sum(h[:])
	return &o
}
