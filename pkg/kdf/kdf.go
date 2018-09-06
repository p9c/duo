// Package kdf implements a blake2b based key derivation function for stretching passwords
package kdf

import (
	"errors"
	. "gitlab.com/parallelcoin/duo/pkg/bytes"
	. "gitlab.com/parallelcoin/duo/pkg/lockedbuffer"
	. "gitlab.com/parallelcoin/duo/pkg/password"
	"golang.org/x/crypto/blake2b"
	"hash"
	"time"
)

// Gen takes a password and a random 12 byte initialisation vector and hashes it using Blake2b-384, returning a 32 byte ciphertext and 12 byte initialisation vector from the first 32 bytes and last 12 bytes respectively, after hashing the resultant hash iterations-1 more times.
//
// Blake2b is used because it is faster than SHA256/SHA512.
func Gen(p *Password, iv *Bytes, iterations int) (C *LockedBuffer, IV *Bytes, err error) {
	if p == nil {
		return nil, nil, errors.New("nil password")
	}
	if iv == nil {
		return nil, nil, errors.New("nil IV")
	}
	if iterations < 1 {
		return nil, nil, errors.New("iterations less than 1")
	}
	buf := NewLockedBuffer().New(p.Len() + iv.Len())
	if buf.Error() != "" {
		return nil, nil, errors.New(buf.Error())
	}
	defer buf.Free()
	bb := buf.Buf().(*[]byte)
	B := *bb
	pp := p.Buf().(*[]byte)
	P := *pp
	for i := range *pp {
		B[i] = P[i]
	}
	for i := 0; i < iv.Len(); i++ {
		B[i+p.Len()] = (*iv.Buf().(*[]byte))[i]
	}
	var blake hash.Hash
	blake, err = blake2b.New384(B)
	last := blake.Sum(nil)
	for i := 1; i < iterations; i++ {
		blake.Write(last)
		last = blake.Sum(B)
	}
	C = NewLockedBuffer().New(32).(*LockedBuffer)
	c := *C.Buf().(*[]byte)
	for i := range c {
		c[i] = last[i]
	}
	IV = NewBytes().New(12).(*Bytes)
	ivb := *IV.Buf().(*[]byte)
	for i := range ivb {
		ivb[i] = last[i+C.Len()]
	}
	return
}

// Bench returns the number of iterations performed in a given time on the current hardware
func Bench(t time.Duration) (iter int) {
	P := NewPassword().Rand(12)
	p := *P.Buf().(*[]byte)
	iv := NewBytes().Rand(12)
	Buf := make([]byte, P.Len()+iv.Len())
	for i := range p {
		Buf[i] = p[i]
	}
	for i := 0; i < iv.Len(); i++ {
		Buf[i+P.Len()] = (*iv.Buf().(*[]byte))[i]
	}
	var blake hash.Hash
	blake, _ = blake2b.New384(Buf)
	timerChan := time.NewTimer(t).C
	last := blake.Sum(nil)
	iter = 1
	for {
		blake.Write(last)
		last = blake.Sum(Buf)
		iter++
		select {
		case <-timerChan:
			P.Free()
			return
		default:
		}
	}
}
