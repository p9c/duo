package kdf

import (
	"crypto/rand"
	"errors"
	"github.com/parallelcointeam/duo/pkg/buf"
	"golang.org/x/crypto/blake2b"
	"hash"
	"time"
)

// Gen takes a password and a random 12 byte initialisation vector and hashes it using Blake2b-384, returning a 32 byte ciphertext that is used to encrypt and decrypt the ciphertext from the crypt
func Gen(p *buf.Secure, iv *buf.Byte, iterations int64) (C *buf.Secure, err error) {
	C = buf.NewSecure()
	switch {
	case p == nil:
		err = errors.New("nil password")
	case iv == nil:
		err = errors.New("nil IV")
	case iterations < 1:
		err = errors.New("iterations less than 1")
	default:
		b := make([]byte, p.Len()+iv.Len())
		b1 := buf.NewSecure()
		b1.Copy(&b)
		if !b1.OK() {
			err = errors.New(b1.Error())
		} else {
			defer b1.Free()
			var blake hash.Hash
			blake, err = blake2b.New384(nil)
			if err != nil {
				return nil, err
			}
			blake.Write(*p.Bytes())
			blake.Write(*b1.Bytes())
			last := blake.Sum(nil)
			for i := int64(1); i < iterations; i++ {
				N := len(last)
				n, err := blake.Write(last)
				if err != nil {
					return nil, err
				}
				if N != n {
					return nil, errors.New("did not get all bytes from hash")
				}
				last = blake.Sum(nil)
			}
			b = make([]byte, 32)
			C = buf.NewSecure()
			C.Copy(&b)
			c := *C.Bytes()
			for i := range c {
				c[i] = last[i]
				last[i] = 0
			}
		}
	}
	return
}

// Bench returns the number of iterations performed in a given time on the current hardware
func Bench(t time.Duration) (iter int64) {
	p := make([]byte, 16)
	rand.Read(p)
	iv := make([]byte, 12)
	rand.Read(iv)
	var blake hash.Hash
	blake, _ = blake2b.New384(nil)
	blake.Write(p)
	blake.Write(iv)
	timerChan := time.NewTimer(t).C
	last := blake.Sum(nil)
	iter = 1
	for {
		last = blake.Sum(last)
		iter++
		select {
		case <-timerChan:
			return
		default:
		}
	}
}
