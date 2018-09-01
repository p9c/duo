package bytes

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	. "gitlab.com/parallelcoin/duo/pkg/interfaces"
	"math/big"
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

// guards against nil pointer receivers
func ifnil(r *Byte) *Byte {
	if r == nil {
		r = new(Byte)
		r.SetError("nil receiver")
	}
	return r
}

/////////////////////////////////////////
// Buffer implementation
/////////////////////////////////////////

// NewByte creates a new byte
func NewByte() Buffer {
	return new(Byte)
}

// Buf returns itself inside a pointer to a slice
func (r *Byte) Buf() *[]byte {
	if r == nil {
		r = NewByte().(*Byte)
		r.Unset()
		r.SetError("nil receiver")
	}
	return &[]byte{r.byte}
}

// Copy copies the (first) byte from another buffer
func (r *Byte) Copy(b Buffer) Buffer {
	if r == nil {
		r = NewByte().(*Byte)
		r.Unset()
		r.SetError("nil receiver")
	}
	r.byte = (*b.Buf())[0]
	return r
}

// ForEach calls a function that is called with an index and allows iteration neatly with a closure
func (r *Byte) ForEach(f func(int)) Buffer {
	for i := range *r.Buf() {
		f(i)
	}
	return r
}

// Free is a no-op for this type, really just an end-of-line for a pipeline statement, only to satisfy the buffer interface
func (r *Byte) Free() Buffer {
	return r
}

// Link points the Byte to the byte of another
func (r *Byte) Link(b Buffer) Buffer {
	if r == nil {
		r = NewByte().(*Byte)
		r.Unset()
		r.SetError("nil receiver")
	}
	r.byte = (*b.Buf())[0]
	return r
}

// Load copies the first byte of a buffer
func (r *Byte) Load(b *[]byte) Buffer {
	if r == nil {
		r = NewByte().(*Byte)
		r.Unset()
		r.SetError("nil receiver")
	}
	r.byte = byte((*b)[0])
	return r
}

// Move copies the first byte of a buffer and zeroes the original
func (r *Byte) Move(b Buffer) Buffer {
	if r == nil {
		r = NewByte().(*Byte)
		r.Unset()
		r.SetError("nil receiver")
	}
	r.byte = (*b.Buf())[0]
	b.Null()
	return r
}

// New creates a new byte (really just zeroes it, size is ignored, 0 input would be usual)
func (r *Byte) New(size int) Buffer {
	if r == nil {
		r = NewByte().(*Byte)
		r.Unset()
		r.SetError("nil receiver")
	}
	r.Null()
	return r
}

// Null makes the Byte zero
func (r *Byte) Null() Buffer {
	r.byte = 0
	r.Unset()
	return r
}

// Rand loads a random byte into the buffer
func (r *Byte) Rand(int) Buffer {
	if r == nil {
		r = NewByte().(*Byte)
		r.Unset()
		r.SetError("nil receiver")
	}
	rand.Read([]byte{r.byte})
	return r
}

// Size returns 1 or if unset, 0
func (r *Byte) Size() int {
	if r.IsSet() {
		return 1
	}
	return 0
}

// Toggle implementation

// IsSet returns the boolean indicating if the variable has been initialised/loaded
func (r *Byte) IsSet() bool {
	return r.isset
}

// Set marks the Byte to set
func (r *Byte) Set() Toggle {
	r.isset = true
	return r
}

// Unset marks the Byte to unset
func (r *Byte) Unset() Toggle {
	r.isset = false
	return r
}

// Status implementation

// SetError sets the status string in the error field
func (r *Byte) SetError(s string) Buffer {
	r.error = errors.New(s)
	return r
}

// UnsetError nils the error in the error field
func (r *Byte) UnsetError() Buffer {
	r.error = nil
	return r
}

/////////////////////////////////////////
// Coding implementation
/////////////////////////////////////////

// Coding returns the coding type to be used by the String function
func (r *Byte) Coding() string {
	if r == nil {
		r = NewByte().(*Byte)
		r.SetError("nil receiver")
	}
	if r.coding >= len(CodeType) {
		r.coding = 0
		r.SetError("invalid coding type in Byte")
	}
	return CodeType[r.coding]
}

// SetCoding changes the encoding type
func (r *Byte) SetCoding(coding string) Buffer {
	if r == nil {
		r = NewByte().(*Byte)
		r.SetError("nil receiver")
	}
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
	copy(R, CodeType)
	return
}

/////////////////////////////////////////
// Array implementation
/////////////////////////////////////////

// Elem returns the value of a specified bit as an 8 bit value (0 or 1)..
func (r *Byte) Elem(index int) Buffer {
	if r == nil {
		r = NewByte().(*Byte)
		r.SetError("nil receiver")
		return &Byte{}
	}
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
	return new(Byte)
}

// Purge makes the byte zero
func (r *Byte) Purge() Array {
	r.byte = 0
	return r
}

// SetElem sets a bit in a byte
func (r *Byte) SetElem(index int, b Buffer) Array {
	if r == nil {
		r = NewByte().(*Byte)
		r.SetError("nil receiver")
		return &Byte{}
	}
	if index > 7 {
		r.SetError("a byte only has 8 bits")
		return &Byte{}
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
	r = ifnil(r)
	switch CodeType[r.coding] {
	case "byte", "string":
		return string(r.byte)
	case "decimal":
		bi := big.NewInt(0)
		bi.SetBytes(*r.Buf())
		return fmt.Sprint(bi)
	case "hex":
		return "0x" + hex.EncodeToString(*r.Buf())
	default:
		return r.SetCoding("decimal").String()
	}
}
