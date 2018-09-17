package blockcrypt

import (
	"crypto/rand"
	"errors"
	"github.com/parallelcointeam/duo/pkg/buf"
	"golang.org/x/crypto/blake2b"
	"hash"
	"time"
)

// Gen takes a password and a random 12 byte initialisation vector and hashes it using Blake2b-384, returning a 32 byte ciphertext.
func Gen(p *buf.Secure, iv *buf.Bytes, iterations int) (C *buf.Secure, IV *buf.Bytes, err error) {
	if iterations < 1 {
		return nil, nil, errors.New("iterations less than 1")
	}
	if p == nil {
		return nil, nil, errors.New("nil password")
	}
	if iv == nil {
		return nil, nil, errors.New("nil IV")
	}
	b := make([]byte, p.Len()+iv.Len())
	b1 := buf.NewSecure().Copy(&b)
	defer b1.Free()
	bb := b1.Bytes()
	B := *bb
	pp := p.Bytes()
	P := *pp
	for i := range *pp {
		B[i] = P[i]
	}
	for i := 0; i < iv.Len(); i++ {
		B[i+p.Len()] = (*iv.Bytes())[i]
	}
	var blake hash.Hash
	blake, err = blake2b.New384(B)
	last := blake.Sum(nil)
	for i := 1; i < iterations; i++ {
		blake.Write(last)
		last = blake.Sum(B)
	}
	b = make([]byte, 32)
	C = buf.NewSecure().Copy(&b).(*buf.Secure)
	c := *C.Bytes()
	for i := range c {
		c[i] = last[i]
		last[i] = 0
	}
	ivv := c[32:44]
	IV = buf.NewBytes().Copy(&ivv).(*buf.Bytes)
	return
}

// Bench returns the number of iterations performed in a given time on the current hardware
func Bench(t time.Duration) (iter int) {
	pp := make([]byte, 16)
	rand.Read(pp)
	P := buf.NewSecure().Copy(&pp).(*buf.Secure)
	p := *P.Bytes()
	ivv := make([]byte, 12)
	iv := buf.NewBytes().Copy(&ivv).(*buf.Bytes)
	Buf := make([]byte, P.Len()+iv.Len())
	for i := range p {
		Buf[i] = p[i]
	}
	for i := 0; i < iv.Len(); i++ {
		Buf[i+P.Len()] = (*iv.Bytes())[i]
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
