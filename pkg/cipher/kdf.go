package cipher

import (
	"errors"
	"gitlab.com/parallelcoin/duo/pkg/buf/bytes"
	"gitlab.com/parallelcoin/duo/pkg/buf/pass"
	"gitlab.com/parallelcoin/duo/pkg/buf/sec"
	"golang.org/x/crypto/blake2b"
	"hash"
	"time"
)

// Gen takes a password and a random 12 byte initialisation vector and hashes it using Blake2b-384, returning a 32 byte ciphertext and 12 byte initialisation vector from the first 32 bytes and last 12 bytes respectively, after hashing the resultant hash iterations-1 more times.
//
// Blake2b is used because it is faster than SHA256/SHA512.
func Gen(p *passbuf.Password, iv *buf.Unsafe, iterations int) (C *buf.Fenced, IV *buf.Unsafe, err error) {
	if iterations < 1 {
		return nil, nil, errors.New("iterations less than 1")
	}
	if p == nil {
		return nil, nil, errors.New("nil password")
	}
	if iv == nil {
		return nil, nil, errors.New("nil IV")
	}

	buf := buf.New().New(p.Len() + iv.Len())
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
	C = buf.New().New(32).(*buf.Fenced)
	c := *C.Buf().(*[]byte)
	for i := range c {
		c[i] = last[i]
	}
	IV = buf.New().New(12).(*buf.Unsafe)
	ivb := *IV.Buf().(*[]byte)
	for i := range ivb {
		ivb[i] = last[i+C.Len()]
	}
	return
}

// Bench returns the number of iterations performed in a given time on the current hardware
func Bench(t time.Duration) (iter int) {
	P := passbuf.New().Rand(12)
	p := *P.Buf().(*[]byte)
	iv := buf.New().Rand(12)
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
