package types

import (
	"fmt"
)

// Bytes is a simple wrapper around the builtin []byte type that is written to make explicit when data is being copied by an application. This is to promote security hygiene and prevent the unintended duplication of sensitive data in the memory of the process, which obviously improves the odds of an attack uncovering them.
type Bytes struct {
	value *[]byte
	err   error
}
type bytes interface {
	Copy() *[]byte
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

// Copy returns a copy of the contents of a Bytes (consumers of this library may only alter this value via functions as it is not exported). Note that this is in contrast to the To and From methods which *move* the data. This also does not use the builtin slice functions because they can cause side effects by referring to the same underlying buffer (this is why the To and From functions *move* the data also, but the consuming function should be aware of the difference because it is specified clearly).
func (b *Bytes) Copy() *Bytes {
	bb := NewBytes().WithSize(b.Len())
	bv := *b.value
	bbv := *bb.value
	for i := range bv {
		bbv[i] = bv[i]
	}
	return bb
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

// ToBytes moves the contents to a newly allocated buffer.
// WARNING: The slice inside the receiver struct will be zeroed!
func (b *Bytes) ToBytes() (B *Bytes) {
	B = NewBytes()
	bb := make([]byte, b.Len())
	bv := *b.value
	for i := range bv {
		bb[i] = bv[i]
		bv[i] = 0
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
