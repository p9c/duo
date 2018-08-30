package testprototype

import (
	. "gitlab.com/parallelcoin/duo/pkg/bytes"
)

type TestPrototype struct {
	val  *[]byte
	set  bool
	utf8 bool
	err  error
}

type testPrototype interface {
	Buf() []byte
	Copy(*Bytes) *Bytes
	Delete()
	Error() string
	IsSet() bool
	IsUTF8() bool
	Len() int
	Link(*Bytes) *Bytes
	Load(*[]byte) *Bytes
	MarshalJSON() ([]byte, error)
	Move(*Bytes) *Bytes
	New(int) *Bytes
	Null() *Bytes
	Rand(int) *Bytes
	SetBin() *Bytes
	SetError(string) *Bytes
	SetUTF8() *Bytes
	String() string
}
