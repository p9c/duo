package types

import (
	"fmt"
)

type Bytes struct {
	value *[]byte
	err   error
}
type bytes interface {
	Value() *[]byte
	Len() int
	WithSize(int) *Bytes
	ToBytes() []byte
	FromBytes([]byte) *Bytes
	ToByteSlice() []byte
	FromByteSlice([]byte) *Bytes
	ToString() string
	FromString(string) *Bytes
}

// ----------------------------------------------------------------------------
// Methods relating to Bytes

// NewBytes creates a new Bytes object
func NewBytes() *Bytes {
	return new(Bytes)
}

// Len returns the length of the data stored in a Bytes
func (b *Bytes) Len() int {
	if b.value == nil {
		return 0
	}
	return len(*b.value)
}

// WithSize allocates a defined size of byte slice in a Bytes structure. Always returns a crispy fresh new buffer, so don't forget you have the old if you call it on an existing Bytes.
func (b *Bytes) WithSize(size int) (B *Bytes) {
	B = new(Bytes)
	bb := make([]byte, size)
	B.FromByteSlice(&bb)
	return
}

// ToBytes returns a copy of a bytes object
func (b *Bytes) ToBytes() (B *Bytes) {
	B = NewBytes()
	bb := make([]byte, b.Len())
	for i := range *b.value {
		bb[i] = (*b.value)[i]
	}
	B.value = &bb
	return B
}

// FromBytes moves the contents of one Bytes to another
func (b *Bytes) FromBytes(B *Bytes) *Bytes {
	if b.value != nil {
		for i := range *b.value {
			(*b.value)[i] = 0
		}
	}
	value := make([]byte, B.Len())
	for i := range *B.value {
		value[i] = (*B.value)[i]
		(*B.value)[i] = 0
	}
	b.value = &value
	return b
}

// ----------------------------------------------------------------------------
// Methods relating to builtin byte slices

// ToByteSlice moves the byte slice inside it out
func (b *Bytes) ToByteSlice() (out *[]byte) {
	out = b.value
	b.value = nil
	return
}

// FromByteSlice moves the content of a byte slice into the Bytes buffer
func (b *Bytes) FromByteSlice(in *[]byte) *Bytes {
	if b.value != nil {
		for i := range *b.value {
			(*b.value)[i] = 0
		}
	}
	I := *in
	B := make([]byte, len(I))
	if in != nil {
		for i := range *in {
			B[i] = (*in)[i]
			I[i] = 0
		}
	}
	b.value = &B
	return b
}

// ToString moves the contents of a string into its buffer.
// WARNING: strings are immutable and you cannot control where or how many times they get copied!
func (b *Bytes) ToString() (s *string) {
	S := string(*b.value)
	return &S
}

// FromString moves the contents of its buffer into a string
// WARNING: strings are immutable and you cannot control where or how many times they get copied!
func (b *Bytes) FromString(s *string) *Bytes {
	S := *s
	B := []byte(S)
	b.value = &B
	return b
}

func init() {
	fmt.Println("pkg/types/bytes.go initialising")
}
