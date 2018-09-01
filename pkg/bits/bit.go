package bytes

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	. "gitlab.com/parallelcoin/duo/pkg/interfaces"
)

// Bit represents a binary bit, as a byte, either 0 or 1
type Bit struct {
	bit    byte
	err    error
	set    bool
	coding int
}

func NewBit() *Bit {
	return new(Bit)
}

// ifnil guards against nil pointer receivers
func ifnil(r *Bit) *Bit {
	if r == nil {
		r = new(Bit)
		r.SetError("nil receiver")
	}
	return r
}

/////////////////////////////////////////
// Buffer implementations
/////////////////////////////////////////

// Buf returns either a 1 or zero. It implements array but does not do anything
func (r *Bit) Buf() *[]byte {
	if r == nil {
		return &[]byte{0}
	}
	return r.Buf()
}

// Copy duplicates the buffer from another Bit.
func (r *Bit) Copy(buf Buffer) Buffer {
	r = ifnil(r)
	r.UnsetError()
	if buf == nil {
		r.Free()
		r.SetError("nil parameter")
		return r
	}
	if r == buf {
		r.SetError("parameter is receiver")
		return r
	}
	if buf.Size() == 0 {
		r.Null()
		r.Load(buf.Buf())
		r.SetError("empty buffer received")
		return r
	}
	if (*buf.Buf())[0] == 0 {
		r.Load(&[]byte{0})
	}
	r.Load(&[]byte{0})
	return r
}

// ForEach calls a function that is called with an index and allows iteration neatly with a closure
func (r *Bit) ForEach(f func(int)) Buffer {
	f(0)
	return r
}

// Free destroys the Bit and dereferences it
func (r *Bit) Free() {
	r.Null()
	r.UnsetError()
}

// Link copies the pointer from another Bit's content, meaning what is written to one will also be visible in the other
func (r *Bit) Link(buf Buffer) Buffer {
	if r == nil {
		r = NewBit()
	}
	r.SetError("bits cannot be linked")
	return r
}

// Load moves the contents of a byte slice into the Bit, erasing the original copy.
func (r *Bit) Load(bytes *[]byte) Buffer {
	if r == nil {
		r = NewBit()
		r.SetError("nil receiver")
	}
	if bytes == nil {
		r.SetError("nil parameter")
	} else {
		r.bit = (*bytes)[0] & 1
		r.Set()
	}
	return r
}

// Move copies the pointer to the buffer into the receiver and nulls the passed Bit
func (r *Bit) Move(buf Buffer) Buffer {
	if r == nil {
		r = NewBit()
	}
	if buf == nil {
		r.err = errors.New("nil parameter")
	} else {
		r.Load(buf.Buf()).UnsetError()
		buf.Free()
		buf.UnsetError()
	}
	return r
}

// New destroys the old Lockedbuffer and assigns a new one with a given length
func (r *Bit) New(size int) Buffer {
	if r == nil {
		r = NewBit()
	}
	r.Null()
	r.set = true
	return r
}

// Null zeroes out a Bit
func (r *Bit) Null() Buffer {
	r.Load(&[]byte{0})
	r.UnsetError()
	return r
}

// Rand loads the Bit with cryptographically random bytes to a specified length, destroying existing buffer if it was set
func (r *Bit) Rand(size int) Buffer {
	if r == nil {
		r = NewBit()
	}
	b := make([]byte, 1)
	rand.Read(b)
	r.bit = b[0]
	r.set = true
	return r
}

// Size returns the length of the Bit if it has been loaded, or -1 if not
func (r *Bit) Size() int {
	if r == nil {
		return -1
	}
	return 1
}

/////////////////////////////////////////
// Error implementation
/////////////////////////////////////////

// Error returns the string in the err field
func (r *Bit) Error() string {
	if r == nil {
		return "nil receiver"
	}
	if r.err != nil {
		return r.err.Error()
	}
	return ""
}

// SetError sets the string of the error in the err field
func (r *Bit) SetError(s string) Buffer {
	if r == nil {
		r = new(Bit)
	}
	r.err = errors.New(s)
	return r
}

// UnsetError sets the error to nil
func (r *Bit) UnsetError() Buffer {
	if r == nil {
		r = NewBit()
	}
	r.err = nil
	return r
}

/////////////////////////////////////////
// Toggle implementation
/////////////////////////////////////////

// IsSet returns true if the Lockedbuffer has been loaded with data
func (r *Bit) IsSet() bool {
	if r == nil {
		return false
	}
	return r.set
}

// Set signifies that the state of the data is consistent
func (r *Bit) Set() Toggle {
	r.set = true
	return r
}

// Unset changes the set flag in a Bit to false and other functions will treat it as empty
func (r *Bit) Unset() Toggle {
	r = ifnil(r)
	r.set = false
	return r
}

/////////////////////////////////////////
// Array implementation
/////////////////////////////////////////
// These are not implemented for this type as it is atomic

// Elem returns the byte at a given index of the buffer
func (r *Bit) Elem(i int) Buffer {
	r = ifnil(r)
	return NewBit().Load(&[]byte{r.bit})
}

// Len returns the length of the array
func (r *Bit) Len() int {
	return 1
}

// Cap returns the amount of elements allocated (can be larger than the size)
func (r *Bit) Cap() int {
	return 1
}

// Purge zeroes out all of the buffers in the array
func (r *Bit) Purge() Array {
	r = ifnil(r)
	r.Load(&[]byte{0})
	return r
}

// SetElem sets an element in the buffer
func (r *Bit) SetElem(i int, b Buffer) Array {
	r = ifnil(r)
	r.Load(&[]byte{(*b.Buf())[0]})
	return r
}

/////////////////////////////////////////
// Coding implementation
/////////////////////////////////////////

// Coding returns the coding type to be used by the String function
func (r *Bit) Coding() string {
	r = ifnil(r)
	for i := range CodeTypes {
		if r.coding == i {
			return CodeTypes[i]
		}
	}
	r.SetError("").
		SetError("invalid CodeType")
	return ""
}

// SetCoding changes the encoding type
func (r *Bit) SetCoding(coding string) Buffer {
	r = ifnil(r)
	found := false
	for i := range CodeTypes {
		if coding == CodeTypes[i] {
			found = true
			r.coding = i
		}
	}
	if !found {
		r.SetError("code type not found")
	}
	return r
}

// Codes returns a copy of the array of CodeTypes
func (r *Bit) Codes() (R []string) {
	copy(R, CodeTypes)
	return
}

/////////////////////////////////////////
// Stringer implementation
/////////////////////////////////////////

// String returns the Bit in the currently set coding format
func (r *Bit) String() string {
	switch CodeTypes[r.coding] {
	case "utf8":
		return string(*r.Buf())
	case "hex":
		return string(append([]byte("0x"), []byte(hex.EncodeToString(*r.Buf()))...))
	default:
		r.SetError("coding type not implemented")
	}
	return ""
}

/////////////////////////////////////////
// JSON implementation
/////////////////////////////////////////

// MarshalJSON renders the data as JSON
func (r *Bit) MarshalJSON() ([]byte, error) {
	r = ifnil(r)
	if r.coding >= len(CodeTypes) {
		r.SetError("invalid coding type set")
	}
	buf := r.bit & 1
	return json.Marshal(&struct {
		Bit    byte
		IsSet  bool
		Coding string
		Error  string
	}{
		Bit:    buf,
		IsSet:  r.set,
		Coding: CodeTypes[r.coding],
		Error:  r.Error(),
	})
}
