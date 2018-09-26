// Package hash160 is a library for turning bytes into a 'hash160' hash, which is an sha256 hash followed by ripemd160, 20 bytes long, used by bitcoin and other addresses derived from the public keys, and used by scripts, which are a special type of address
package hash160

import (
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)

// Sum returns the result of sha256 and then ripemd160 on a message
func Sum(message *[]byte) *[]byte {
	h := sha256.Sum256(*message)
	o := ripemd160.New()
	o.Write(h[:])
	O := o.Sum(nil)
	return &O
}
