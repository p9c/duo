package key

import (
	"crypto/aes"
	"crypto/cipher"
)

// Check if a key is valid
func Check(b []byte) bool {
	Max := [32]byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE,
		0xBA, 0xAE, 0xDC, 0xE6, 0xAF, 0x48, 0xA0, 0x3B,
		0xBF, 0xD2, 0x5E, 0x8C, 0xD0, 0x36, 0x41, 0x40,
	}
	IsZero := true
	for i := range b {
		if b[i] != 0 {
			IsZero = false
		}
	}
	if IsZero {
		return false
	}
	result := true
	for i := range b {
		if b[i] > Max[i] {
			result = false
		}
	}
	return result
}

// NewCipher creates a new aes-cbc-256 block cipher
func NewCipher(passphrase []byte) (block cipher.Block, err error) {
	block, err = aes.NewCipher(append(passphrase, make([]byte, len(passphrase)-32)...))
	return
}
