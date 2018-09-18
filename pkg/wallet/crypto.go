package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"errors"
	"github.com/awnumar/memguard"
	"unsafe"
)

var (
	ErrPaddingSize = errors.New("padding size error")
	PKCS5          = &pkcs5{}
	// difference with pkcs5 only block must be 8
	PKCS7 = &pkcs5{}
)

// pkcs5Padding is a pkcs5 padding struct.
type pkcs5 struct{}

// Padding implements the Padding interface Padding method.
func (p *pkcs5) Padding(s *memguard.LockedBuffer, blockSize int) (r *memguard.LockedBuffer, err error) {
	srcLen := len(s.Buffer())
	padLen := blockSize - (srcLen % blockSize)
	var S *memguard.LockedBuffer
	S, err = NewBuffer(padLen)
	if err != nil {
		return
	}
	for i := range S.Buffer() {
		S.Buffer()[i] = byte(padLen)
	}
	r, err = memguard.Concatenate(s, S)
	AllLockedBuffers = append(AllLockedBuffers, r)
	AllocatedBufferCount++
	AllocatedBufferTotalSize += r.Size()
	if err != nil {
		return
	}
	DeleteBuffers(S, s)
	return
}

// Unpadding implements the Padding interface Unpadding method.
func (p *pkcs5) Unpadding(s *memguard.LockedBuffer, blockSize int) (r *memguard.LockedBuffer, err error) {
	srcLen := s.Size()
	paddingLen := int(s.Buffer()[srcLen-1])
	if paddingLen >= srcLen || paddingLen > blockSize {
		err = ErrPaddingSize
		return
	}
	r, err = memguard.Trim(s, srcLen, paddingLen)
	DeleteBuffer(s)
	return
}

// Decrypts a ciphertext using the masterkey and password
func (s *Serializable) DeriveCipher(pass *memguard.LockedBuffer) (mk *MasterKey, k *memguard.LockedBuffer, iv []byte, err error) {
	var mode cipher.BlockMode
	for i := range s.masterKey {
		pLen, sLen := len(pass.Buffer()), len((s.masterKey)[i].Salt)
		var Buf *memguard.LockedBuffer
		Buf, err = NewBuffer(pLen + sLen)
		buf := Buf.Buffer()
		if err != nil {
			return
		}
		for j := range pass.Buffer() {
			buf[j] = pass.Buffer()[j]
		}
		for j := range s.masterKey[i].Salt {
			buf[j+pLen] = s.masterKey[i].Salt[j]
		}
		seed, _ := PKCS7.Padding(Buf, 8)
		var l *memguard.LockedBuffer
		l, err = NewBuffer(64)
		if err != nil {
			return
		}
		var source *[64]byte
		source = (*[64]byte)(unsafe.Pointer(&l.Buffer()[0]))
		*source = sha512.Sum512(seed.Buffer())
		for j := 0; j < int(s.masterKey[i].Iterations-1); j++ {
			*source = sha512.Sum512(l.Buffer())
		}
		k, err = NewBuffer(64)
		if err != nil {
			return
		}
		for j := range k.Buffer() {
			k.Buffer()[j] = source[j]
		}
		var ckey, ivb *memguard.LockedBuffer
		ckey, ivb, err = memguard.Split(k, 32)
		if err != nil {
			return
		}
		block, err := aes.NewCipher(ckey.Buffer())
		if err != nil {
			break
		}
		iv = ivb.Buffer()[:block.BlockSize()]
		mode = cipher.NewCBCDecrypter(block, iv)
		mode.CryptBlocks(k.Buffer(), s.masterKey[i].EncryptedKey)
		mk = s.masterKey[i]
		// DeleteBuffers(k, l)
	}
	if mk == nil {
		err = errors.New("Password did not unlock any of the available master keys")
	}
	return
}
