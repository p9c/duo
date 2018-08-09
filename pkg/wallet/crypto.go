package wallet

import (
	"crypto/aes"
	"crypto/sha512"
)
// Encrypts plaintext using the masterkey and a password
func (m *MKey) Encrypt(p string, b ...[]byte) (r [][]byte, err error) {
	return m.encDec(true, p, b...)
}

// Decrypts a ciphertext using the masterkey and password
func (m *MKey) Decrypt(p string, b ...[]byte) (r [][]byte, err error) {
	return m.encDec(false, p, b...)
}
func (m *MKey) encDec(enc bool, p string, b ...[]byte) (r [][]byte, err error) {
	r = make([][]byte, len(b))
	for i := range r {
		r[i] = make([]byte, 16)
	}
	s := append([]byte(p), m.Salt...)
	sLen := len(s)
	var source [64]byte
	for i := range source {
		if i < sLen {
			source[i] = s[i]
		} else {
			source[i] = byte(64 - sLen)
		}
	}
	for i := 0; i < int(m.Iterations); i++ {
		source = sha512.Sum512(source[:])
	}
	S := make([]byte, 32)
	for i := range S {
		S[i] = source[i]
	}
	if block, err := aes.NewCipher(S); err != nil {
		return nil, err
	} else {
		for i := range b {
			if enc {
				block.Encrypt(r[i], b[i])
			} else {
				block.Decrypt(r[i], b[i])
			}
		}
		for i := range source { source[i] = 0 }
		for i := range S { S[i] = 0 }
		for i := range b { 
			for j := range b[i] {
				b[i][j] = 0
			}
		}
		}	
	return
}
