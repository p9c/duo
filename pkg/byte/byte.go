package bits

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	. "gitlab.com/parallelcoin/duo/pkg/interfaces"
)

// Byte is an implementation of the Buffer, Toggle and Status interface
//
// Note that except for its use in the Bytes type as its' array element most of this is just empty in order to enable the abstraction of the single byte as the element in the Bytes Array implementation.
type Byte struct {
	byte
	isset  bool
	coding int
	error
}

/////////////////////////////////////////
// Buffer implementation
/////////////////////////////////////////

// NewByte creates a new byte
func NewByte() (R Buffer) {
	return new(Byte)
}

func (r *Byte) ifnil() (R *Byte) {
	if r == nil {
		return NewByte().SetError("nil receiver").(*Byte)
	}
	return r
}

// Buf returns itself inside a pointer to a slice
func (r *Byte) Buf() (R *[]byte) {
	r = r.ifnil()
	return &[]byte{r.byte}
}

// Copy copies the (first) byte from another buffer
func (r *Byte) Copy(b Buffer) (R Buffer) {
	r = r.ifnil()
	if b.Size() < 1 {
		return r
	}
	r.byte = (*b.Buf())[0]
	return r
}

// ForEach calls a function that is called with an index and allows iteration neatly with a closure
func (r *Byte) ForEach(f func(int)) (R Buffer) {
	r = r.ifnil()
	f(0)
	return r
}

// Free for byte just creates a fresh new struct and returns it
func (r *Byte) Free() (R Buffer) {
	r = r.ifnil()
	return r
}

// Link points the Byte to the byte of another
func (r *Byte) Link(b Buffer) (R Buffer) {
	r = r.ifnil()
	r.byte = (*b.Buf())[0]
	return r
}

// Load copies the first byte of a buffer
func (r *Byte) Load(b *[]byte) (R Buffer) {
	r = r.ifnil()
	if len(*b) >= 1 {
		r.byte = byte((*b)[0])
	}
	if len(*b) <= 0 {
		r.byte = 0
	}
	return r
}

// Move copies the first byte of a buffer and zeroes the original
func (r *Byte) Move(b Buffer) (R Buffer) {
	r = r.ifnil()
	r.byte = (*b.Buf())[0]
	b.Null()
	return r
}

// New creates a new byte (really just zeroes it, size is ignored, 0 input would be usual)
func (r *Byte) New(size int) (R Buffer) {
	r = r.ifnil()
	r.Null()
	return r
}

// Null makes the Byte zero
func (r *Byte) Null() (R Buffer) {
	r = r.ifnil()
	r.byte = 0
	return r
}

// Rand loads a random byte into the buffer
func (r *Byte) Rand(int) (R Buffer) {
	r = r.ifnil()
	rand.Read([]byte{r.byte})
	return r
}

// Size returns 1 or if unset, 0
func (r *Byte) Size() int {
	r = r.ifnil()
	if r.IsSet() {
		return 1
	}
	return 0
}

// IsSet returns the boolean indicating if the variable has been initialised/loaded
func (r *Byte) IsSet() bool {
	r = r.ifnil()
	return r.isset
}

// Set marks the Byte to set
func (r *Byte) Set() Buffer {
	r = r.ifnil()
	r.isset = true
	return r
}

// Unset marks the Byte to unset
func (r *Byte) Unset() Buffer {
	r = r.ifnil()
	r.isset = false
	return r
}

// Status implementation

// SetError sets the status string in the error field
func (r *Byte) SetError(s string) (R Buffer) {
	r = r.ifnil()
	r.error = errors.New(s)
	return r
}

// UnsetError nils the error in the error field
func (r *Byte) UnsetError() (R Buffer) {
	r = r.ifnil()
	r.error = nil
	return r
}

/////////////////////////////////////////
// Coding implementation
/////////////////////////////////////////

// Coding returns the coding type to be used by the String function
func (r *Byte) Coding() string {
	r = r.ifnil()
	if r.coding >= len(CodeType) {
		r.coding = 0
		r.SetError("invalid coding type in Byte")
	}
	return CodeType[r.coding]
}

// SetCoding changes the encoding type
func (r *Byte) SetCoding(coding string) (R Buffer) {
	r = r.ifnil()
	found := false
	for i := range CodeType {
		if coding == CodeType[i] {
			found = true
			r.coding = i
		}
	}
	if !found {
		r.SetError("code type not found")
	}
	return r
}

// Codes returns a copy of the array of CodeType
func (r *Byte) Codes() (R []string) {
	r = r.ifnil()
	copy(R, CodeType)
	return
}

/////////////////////////////////////////
// Array implementation
/////////////////////////////////////////

// Elem returns the value of a specified bit as an 8 bit value (0 or 1)..
func (r *Byte) Elem(index int) (R Buffer) {
	r = r.ifnil()
	if index > 7 {
		r.SetError("a byte only has 8 bits")
		return &Byte{}
	}
	temp := r.byte
	for i := 0; i < index; i++ {
		temp >>= 1
	}
	return &Byte{byte: temp & 1}
}

// Len always returns 8
func (r *Byte) Len() int {
	return 8
}

// Cap always returns 8
func (r *Byte) Cap() int {
	return 8
}

// Make allocates a new zero byte
func (r *Byte) Make(int, int) Array {
	r = new(Byte)
	return r
}

// Purge makes the byte zero
func (r *Byte) Purge() Array {
	r = r.ifnil()
	r.byte = 0
	return r
}

// SetElem sets a bit in a byte
func (r *Byte) SetElem(index int, b Buffer) Array {
	r = r.ifnil()
	if index > 7 {
		r.SetError("a byte only has 8 bits")
		return r
	}
	mask := byte(1)
	for i := 0; i < index; i++ {
		mask *= 2
	}
	r.byte &= mask
	return r
}

/////////////////////////////////////////
// Stringer implementation
/////////////////////////////////////////

// String returns a string containing the buffer encoded according to the setting in force
func (r *Byte) String() string {
	r = r.ifnil()
	switch CodeType[r.coding] {
	case "byte":
		return fmt.Sprint(r.Buf())
	case "string":
		return string(r.byte)
	case "decimal":
		return fmt.Sprint(r.byte)
	case "hex":
		return "0x" + hex.EncodeToString(*r.Buf())
	default:
		return r.SetCoding("decimal").String()
	}
}
