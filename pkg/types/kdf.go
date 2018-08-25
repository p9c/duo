package types

import (
	"fmt"
	"golang.org/x/crypto/blake2b"
	"hash"
	"time"
)

// KDF takes a password and a random 16 byte initialisation vector and hashes it using Blake2b-384, returning a 32 byte ciphertext and 16 byte initialisation vector from the first 32 bytes and last 16 bytes respectively, after hashing the resultant hash iterations-1 more times.
// Blake2b is used because it is faster than SHA256/SHA512.
func KDF(p *Password, iv *Bytes, iterations int) (C *LockedBuffer, IV *Bytes, err error) {
	if err != nil {
		return
	}
	buf := NewLockedBuffer().WithSize(p.Len() + iv.Len())
	defer buf.Delete()
	Buf := *buf.Buffer()
	P := *p.Buffer()
	for i := range P {
		Buf[i] = P[i]
	}
	for i := 0; i < iv.Len(); i++ {
		Buf[i+p.Len()] = (*iv.Buffer())[i]
	}
	var blake hash.Hash
	blake, err = blake2b.New384(Buf)
	last := blake.Sum(nil)
	for i := 1; i < iterations; i++ {
		blake.Write(last)
		last = blake.Sum(Buf)
	}
	C = NewLockedBuffer().WithSize(32)
	c := *C.Buffer()
	for i := range c {
		c[i] = last[i]
	}
	IV = NewBytes().WithSize(16)
	ivb := *IV.Buffer()
	for i := range ivb {
		ivb[i] = last[i+C.Len()]
	}
	return
}

// KDFBench returns the number of iterations performed in a given time
func KDFBench(t time.Duration) (iter int) {
	timer := time.NewTimer(t)
	fmt.Println(t)
	return
}
