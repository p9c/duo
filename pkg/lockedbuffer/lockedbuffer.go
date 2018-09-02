// Package lockedbuffer is a wrapper around the memguard LockedBuffer that automatically handles destroying data no longer needed and enables copy, link and move functions on the data contained inside the structure.
package lockedbuffer

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/awnumar/memguard"
	. "gitlab.com/parallelcoin/duo/pkg/byte"
	. "gitlab.com/parallelcoin/duo/pkg/interfaces"
	"math/big"
)

// LockedBuffer is a struct that stores and manages memguard.LockedBuffers, ensuring that buffers are destroyed when no longer in use.
//
// Do not use struct literals and not assign them to a name and null() (deletes and zeroes struct) afterwards, or you could run out of memguard LockedBuffers
//
// All functions except for those exporting buffers will automatically allocate the struct of the receiver if it is nil. This permits the use of struct literals in assignment for one-liners that initialise values and call a function with data. It may introduce side effects in code if you did not intend to create a new variable.
//
// The maximum size of buffer is around 172500 bytes on a linux 4.18, it may be more may be less.
type LockedBuffer struct {
	buf    *memguard.LockedBuffer
	set    bool
	coding int
	err    error
}

func donil(r *Bytes, f1 func(), f2 func()) *Bytes {
	if r == nil {
		f1()
	}
	if f2 != nil {
		f2()
	}
	return r
}

func doif(b bool, fn func()) {
	if b {
		fn()
	}
}

// NewLockedBuffer clears the passed LockedBuffer or creates a new one if null
func NewLockedBuffer() *LockedBuffer {
	return new(LockedBuffer)
}

/////////////////////////////////////////
// Buffer implementations
/////////////////////////////////////////

// Buf returns a pointer to the byte slice in the LockedBuffer.
func (r *LockedBuffer) Buf() *[]byte {
	r = r.ifnil()
	if r.set {
		if r.buf != nil {
			a := r.buf.Buffer()
			return &a
		}
	}
	return nil
}

// Copy duplicates the buffer from another LockedBuffer.
func (r *LockedBuffer) Copy(buf Buffer) Buffer {
	if r == nil {
		r = NewLockedBuffer()
		r.SetError("nil receiver")
	}
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
	r.New(buf.Size())
	r.ForEach(func(i int) {
		r.SetElem(i, buf.Elem(i))
	})
	r.Set()
	return r
}

// ForEach calls a function that is called with an index and allows iteration neatly with a closure
func (r *LockedBuffer) ForEach(f func(int)) Buffer {
	r = ifnil(r)
	if r.buf == nil {
		return r
	}
	for i := range *r.Buf() {
		f(i)
	}
	return r
}

// Free destroys the LockedBuffer and dereferences it
func (r *LockedBuffer) Free() Buffer {
	r = ifnil(r)
	if r.buf != nil {
		r.buf.Destroy()
		r.buf = nil
	}
	r.UnsetError().Unset()
	return r
}

// Link copies the pointer from another LockedBuffer's content, meaning what is written to one will also be visible in the other
func (r *LockedBuffer) Link(buf Buffer) Buffer {
	if r == nil {
		r = NewLockedBuffer()
	}
	r.Null()
	r.Load(buf.Buf())
	return r
}

// Load moves the contents of a byte slice into the LockedBuffer, erasing the original copy.
func (r *LockedBuffer) Load(bytes *[]byte) Buffer {
	r = ifnil(r)
	if bytes == nil {
		r.SetError("nil parameter")
	} else {
		r.Null()
		if r.buf, r.err = memguard.NewMutableFromBytes(*bytes); r.err != nil {
			return r
		}
		r.Set()
	}
	return r
}

// Move copies the pointer to the buffer into the receiver and nulls the passed LockedBuffer
func (r *LockedBuffer) Move(buf Buffer) Buffer {
	if r == nil {
		r = NewLockedBuffer()
	}
	if buf == nil {
		r.err = errors.New("nil parameter")
	} else {
		r.buf = buf.(*LockedBuffer).buf
		r.UnsetError().Set()
		buf.(*LockedBuffer).buf = nil
		buf.UnsetError().Unset()
	}
	return r
}

// New destroys the old Lockedbuffer and assigns a new one with a given length
func (r *LockedBuffer) New(size int) Buffer {
	if r == nil {
		r = NewLockedBuffer()
	}
	r.Null()
	r.buf, r.err = memguard.NewMutable(size)
	r.set = true
	return r
}

// Null zeroes out a LockedBuffer
func (r *LockedBuffer) Null() Buffer {
	r = ifnil(r)
	if r.buf != nil {
		r.buf.Wipe()
	}
	return r
}

// Rand loads the LockedBuffer with cryptographically random bytes to a specified length, destroying existing buffer if it was set
func (r *LockedBuffer) Rand(size int) Buffer {
	r = ifnil(r)
	r.New(size)
	r.buf, r.err = memguard.NewMutableRandom(size)
	if r.err != nil {
		return r
	}
	r.set = true
	return r
}

// Size returns the length of the LockedBuffer if it has been loaded, or -1 if not
func (r *LockedBuffer) Size() int {
	if r == nil {
		return -1
	}
	if r.buf == nil {
		return 0
	}
	return r.buf.Size()
}

/////////////////////////////////////////
// Coding implementation
/////////////////////////////////////////

// Coding returns the coding type to be used by the String function
func (r *LockedBuffer) Coding() string {
	r = ifnil(r)
	if r.coding >= len(CodeType) {
		r.coding = 0
		r.SetError("invalid coding type in LockedBuffer")
	}
	return CodeType[r.coding]
}

// SetCoding changes the encoding type
func (r *LockedBuffer) SetCoding(coding string) Buffer {
	r = ifnil(r)
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
func (r *LockedBuffer) Codes() (R []string) {
	copy(R, CodeType)
	return
}

/////////////////////////////////////////
// Status implementation
/////////////////////////////////////////

// Error returns the string in the err field
func (r *LockedBuffer) Error() string {
	r = ifnil(r)
	if r.err != nil {
		return r.err.Error()
	}
	return ""
}

// SetError sets the string of the error in the err field
func (r *LockedBuffer) SetError(s string) Buffer {
	r = ifnil(r)
	r.err = errors.New(s)
	return r
}

// UnsetError sets the error to nil
func (r *LockedBuffer) UnsetError() Buffer {
	r = ifnil(r)
	r.err = nil
	return r
}

/////////////////////////////////////////
// Array implementation
/////////////////////////////////////////

// Elem returns the byte at a given index of the buffer
func (r *LockedBuffer) Elem(i int) Buffer {
	r = ifnil(r)
	if r.buf == nil {
		r.SetError("nil buffer")
		return &LockedBuffer{}
	}
	return NewByte().Load(&[]byte{r.buf.Buffer()[i]})
}

// Len returns the length of the array
func (r *LockedBuffer) Len() int {
	return r.Size()
}

// Cap returns the amount of elements allocated (can be larger than the size)
func (r *LockedBuffer) Cap() int {
	if r == nil || r.buf == nil {
		return 0
	}
	return cap(*r.Buf())
}

// SetElem sets an element in the buffer
func (r *LockedBuffer) SetElem(i int, b Buffer) Array {
	r = ifnil(r)
	if r.buf == nil {
		r.SetError("nil value")
		return r
	}
	R := r.buf.Buffer()
	R[i] = (*b.Buf())[0]
	return r
}

// Purge zeroes out all of the buffers in the array
func (r *LockedBuffer) Purge() Array {
	r = ifnil(r)
	if r.buf == nil {
		r.SetError("nil buffer")
		return r
	}
	r.buf.Wipe()
	return r
}

/////////////////////////////////////////
// Toggle implementation
/////////////////////////////////////////

// IsSet returns true if the Lockedbuffer has been loaded with data
func (r *LockedBuffer) IsSet() bool {
	if r == nil {
		return false
	}
	return r.set
}

// Set signifies that the state of the data is consistent
func (r *LockedBuffer) Set() Toggle {
	r = ifnil(r)
	r.set = true
	return r
}

// Unset changes the set flag in a LockedBuffer to false and other functions will treat it as empty
func (r *LockedBuffer) Unset() Toggle {
	r = ifnil(r)
	r.set = false
	return r
}

/////////////////////////////////////////
// JSON
/////////////////////////////////////////

// MarshalJSON marshals the data of this object into JSON
func (r *LockedBuffer) MarshalJSON() ([]byte, error) {
	r = ifnil(r)
	return json.Marshal(&struct {
		Value  string
		IsSet  bool
		Coding string
		Error  string
	}{
		Value:  r.String(),
		IsSet:  r.set,
		Coding: CodeType[r.coding],
		Error:  r.Error(),
	})
}

/////////////////////////////////////////
// Stringer implementation
/////////////////////////////////////////

// String returns the JSON representing the data in a LockedBuffer
func (r *LockedBuffer) String() string {
	r = ifnil(r)
	if r.buf == nil {
		return "<nil>"
	}
	if r.coding > len(CodeType) {
		r.SetError("invalid coding")
		r.SetCoding("decimal")
	}
	switch CodeType[r.coding] {
	case "byte":
		return fmt.Sprint(*r.Buf())
	case "string":
		return string(*r.Buf())
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
