// Package lockedbuffer is a wrapper around the memguard LockedBuffer that automatically handles destroying data no longer needed and enables copy, link and move functions on the data contained inside the structure.
package lockedbuffer

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/awnumar/memguard"
)

// LockedBuffer is a struct that stores and manages memguard.LockedBuffers, ensuring that buffers are destroyed when no longer in use.
//
// Do not use struct literals and not assign them to a name and null() (deletes and zeroes struct) afterwards, or you could run out of memguard LockedBuffers
//
// All functions except for those exporting buffers will automatically allocate the struct of the receiver if it is nil. This permits the use of struct literals in assignment for one-liners that initialise values and call a function with data. It may introduce side effects in code if you did not intend to create a new variable.
//
// The maximum size of buffer is around 172500 bytes on a linux 4.18, it may be more may be less.
type LockedBuffer struct {
	val  *memguard.LockedBuffer
	set  bool
	utf8 bool
	err  error
}

// NewLockedBuffer clears the passed LockedBuffer or creates a new one if null
func NewLockedBuffer(r ...*LockedBuffer) *LockedBuffer {
	if len(r) == 0 {
		r = append(r, new(LockedBuffer))
	}
	if r[0] == nil {
		r[0] = new(LockedBuffer)
		r[0].SetError("receiver was nil")
	}
	if r[0].set {
		r[0].val.Destroy()
	}
	r[0].val, r[0].set, r[0].err = nil, false, nil
	return r[0]
}

type lockedBuffer interface {
	Buf() []byte
	Copy(*LockedBuffer) *LockedBuffer
	Delete()
	Error() string
	IsSet() bool
	IsUTF8() bool
	Len() int
	Link(*LockedBuffer) *LockedBuffer
	Load(*[]byte) *LockedBuffer
	MarshalJSON() ([]byte, error)
	Move(*LockedBuffer) *LockedBuffer
	New(int) *LockedBuffer
	Null() *LockedBuffer
	Rand(int) *LockedBuffer
	SetBin() *LockedBuffer
	SetError(string) *LockedBuffer
	SetUTF8() *LockedBuffer
	String() string
}

// Buf returns a pointer to the byte slice in the LockedBuffer.
//
// Note that this buffer cannot be treated as a regular byte slice, or it will likely trample the canaries or leave a dangling pointer if it is.
func (r *LockedBuffer) Buf() *[]byte {
	if r == nil {
		return &[]byte{}
	}
	if r.set {
		a := r.val.Buffer()
		return &a
	}
	return nil
}

// Copy duplicates the buffer from another LockedBuffer.
func (r *LockedBuffer) Copy(buf *LockedBuffer) *LockedBuffer {
	if r == nil {
		r = NewLockedBuffer()
		r.err = errors.New("nil receiver")
	}
	r.err = nil
	if buf == nil {
		r.Null()
		r.err = errors.New("nil parameter")
		return r
	}
	if r == buf {
		r.err = errors.New("parameter is receiver")
		return r
	}
	if buf.Len() == 0 {
		r.Null()
		r.val = buf.val
		r.err = errors.New("empty buffer received")
		return r
	}
	r = r.New(buf.Len())
	a := *r.Buf()
	b := *buf.Buf()
	for i := range a {
		a[i] = b[i]
	}
	r.set = true
	return r
}

// Delete deletes the memguard LockedBuffer and dereferences it
func (r *LockedBuffer) Delete() {
	r.Null()
}

// Error returns the string in the err field
func (r *LockedBuffer) Error() string {
	if r == nil {
		return "nil receiver"
	}
	if r.err != nil {
		return r.err.Error()
	}
	return ""
}

// IsSet returns true if the Lockedbuffer has been loaded with a data
func (r *LockedBuffer) IsSet() bool {
	if r == nil {
		return false
	}
	return r.set
}

// IsUTF8 returns true if the LockedBuffer is set to output UTF8
func (r *LockedBuffer) IsUTF8() bool {
	return r.utf8
}

// Len returns the length of the LockedBuffer if it has been loaded, or -1 if not
func (r *LockedBuffer) Len() int {
	if r == nil {
		return -1
	}
	if r.set {
		return r.val.Size()
	}
	return 0
}

// Link copies the pointer from another LockedBuffer's content, meaning what is written to one will also be visible in the other
func (r *LockedBuffer) Link(buf *LockedBuffer) *LockedBuffer {
	if r == nil {
		r = NewLockedBuffer()
	}
	r.Null()
	r.val, r.set = buf.val, buf.set
	return r
}

// Load moves the contents of a byte slice into the LockedBuffer, erasing the original copy.
func (r *LockedBuffer) Load(bytes *[]byte) *LockedBuffer {
	if r == nil {
		r = NewLockedBuffer()
	}
	if bytes == nil {
		r.SetError("nil parameter")
	} else {
		r.Null()
		r.val, r.err = memguard.NewMutableFromBytes(*bytes)
		if r.err != nil {
			return r
		}
		r.set = true
	}
	return r
}

// MarshalJSON marshals the data of this object into JSON
func (r *LockedBuffer) MarshalJSON() ([]byte, error) {
	if r == nil {
		r = NewLockedBuffer()
		r.SetError("nil receiver")
	}
	var val string
	if r.IsSet() {
		if r.utf8 {
			val = string(*r.Buf())
		} else {
			val = string(append([]byte("0x"), hex.EncodeToString(*r.Buf())...))
		}
	}
	return json.Marshal(&struct {
		Value string
		IsSet bool
		Error string
	}{
		Value: val,
		IsSet: r.set,
		Error: r.Error(),
	})
}

// Move copies the pointer to the buffer into the receiver and nulls the passed LockedBuffer
func (r *LockedBuffer) Move(buf *LockedBuffer) *LockedBuffer {
	if r == nil {
		r = NewLockedBuffer()
	}
	if buf == nil {
		r.err = errors.New("nil parameter")
	} else {
		r.val, r.set, r.err = buf.val, true, nil
		buf.val, buf.set, buf.err = nil, false, nil
	}
	return r
}

// New destroys the old Lockedbuffer and assigns a new one with a given length
func (r *LockedBuffer) New(size int) *LockedBuffer {
	if r == nil {
		r = NewLockedBuffer()
	}
	r.Null()
	r.val, r.err = memguard.NewMutable(size)
	r.set = true
	return r
}

// Null zeroes out a LockedBuffer
func (r *LockedBuffer) Null() *LockedBuffer {
	return NewLockedBuffer(r)
}

// Rand loads the LockedBuffer with cryptographically random bytes to a specified length, destroying existing buffer if it was set
func (r *LockedBuffer) Rand(size int) *LockedBuffer {
	r = r.New(size)
	r.val, r.err = memguard.NewMutableRandom(size)
	if r.err != nil {
		return r
	}
	r.set = true
	return r
}

// SetBin sets the data type to be binary
func (r *LockedBuffer) SetBin() *LockedBuffer {
	if r == nil {
		r = NewLockedBuffer()
		r.SetError("nil receiver")
	}
	r.utf8 = false
	return r
}

// SetError sets the string of the error in the err field
func (r *LockedBuffer) SetError(s string) *LockedBuffer {
	if r == nil {
		r = new(LockedBuffer)
	}
	r.err = errors.New(s)
	return r
}

// SetUTF8 sets the LockedBuffer to output UTF8 strings
func (r *LockedBuffer) SetUTF8() *LockedBuffer {
	if r == nil {
		r = NewLockedBuffer()
		r.SetError("nil receiver")
	}
	r.utf8 = true
	return r
}

// String returns the JSON representing the data in a LockedBuffer
func (r *LockedBuffer) String() string {
	s, _ := json.MarshalIndent(r, "", "    ")
	return string(s)
}
