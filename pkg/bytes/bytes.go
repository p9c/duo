// Package bytes is a wrapper around the native byte slice that automatically handles purging discarded data and enables copy, link and move functions on the data contained inside the structure.
//
// The purpose of this structure is to enable the chaining of pointer methods and eliminate the need for separate assignments by passing error value within the structure instead of as the last term in the return tuple. This structuring has a similar functionality to default parameters, without the compile-time complexity. The same pattern can be used to extend the type to be incorporated into a new aggregate type that can contain more similar structures or add them in addition to implemented methods.
package bytes

import (
	"crypto/rand"
	"errors"
)

// Bytes is a struct that stores and manages byte slices for security purposes, automatically wipes old data when new data is loaded.
//
// The structure stores a boolean signifying whether its value is set to point at a valid slice or not, and an error value, which allows one to use the type in assignments without multiple return values, while still allowing one to check the error value of functions performed with it.
//
// To use it, simply new(Bytes) to get pointer to a empty new structure, and then after that you can call the methods of the interface.
type Bytes struct {
	val *[]byte
	set bool
	err error
}

// NewBytes empties an existing bytes or makes a new one
func NewBytes(r ...*Bytes) *Bytes {
	if len(r) == 0 {
		r = append(r, new(Bytes))
	}
	if r[0] == nil {
		r[0] = new(Bytes)
		r[0].SetError("receiver was nil")
	}
	if r[0].val != nil {
		rr := *r[0].val
		if r[0].set {
			for i := range rr {
				rr[i] = 0
			}
		}
	}
	r[0].val, r[0].set, r[0].err = nil, false, nil
	return r[0]
}

type bytes interface {
	Buf() []byte
	Copy(*Bytes) *Bytes
	Delete()
	Error() string
	IsSet() bool
	Len() int
	Link(*Bytes) *Bytes
	Load(*[]byte) *Bytes
	Move(*Bytes) *Bytes
	New(int) *Bytes
	Null() *Bytes
	Rand(int) *Bytes
	SetError(string) *Bytes
}

// Buf returns a variable pointing to the value stored in a Bytes.
func (r *Bytes) Buf() *[]byte {
	if r == nil || r.val == nil {
		return &[]byte{}
	}
	return r.val
}

// Copy duplicates the data from the *[]byte provided and zeroes and replaces its contents, clearing the error value.
func (r *Bytes) Copy(bytes *Bytes) *Bytes {
	if r == nil {
		r = NewBytes()
		r.err = errors.New("nil pointer receiver")
	}
	r.err = nil
	if bytes == nil {
		r.Null()
		r.err = errors.New("nil parameter")
		return r
	}
	if r == bytes {
		r.err = errors.New("parameter is receiver")
		return r
	}
	if len(*bytes.Buf()) == 0 {
		r.Null()
		r.val = bytes.val
		r.err = errors.New("empty buffer received")
		return r
	}
	r = r.New(bytes.Len())
	a := *r.Buf()
	b := *bytes.Buf()
	for i := range b {
		a[i] = b[i]
	}
	r.set = true
	return r
}

// Delete wipes the buffer and dereferences it
func (r *Bytes) Delete() {
	r.Null()
}

// Error gets the error string
func (r *Bytes) Error() string {
	if r == nil {
		return "nil receiver"
	}
	if r.err != nil {
		return r.err.Error()
	}
	return ""
}

// IsSet returns true if the Bytes buffer has been loaded with a slice
func (r *Bytes) IsSet() bool {
	if r == nil {
		return false
	}
	return r.set
}

// Len returns the length of the *[]byte if it has a value assigned, or -1
func (r *Bytes) Len() int {
	if r == nil {
		return -1
	}
	if r.set {
		return len(*r.val)
	}
	return 0
}

// Link nulls the Bytes and copies the pointer in from another Bytes. Whatever is done to one's []byte will also affect the other, but they keep separate error values
func (r *Bytes) Link(bytes *Bytes) *Bytes {
	if r == nil {
		r = NewBytes(nil)
	}
	r.Null()
	r.val, r.set = bytes.val, bytes.set
	return r
}

// Load nulls any existing data and sets its pointer to refer to the pointer to byte slice in the parameter.
func (r *Bytes) Load(bytes *[]byte) *Bytes {
	if r == nil {
		r = NewBytes()
	}
	if bytes == nil {
		r.SetError("nil parameter")
	} else {
		r.Null()
		r.val, r.set, r.err = bytes, true, nil
	}
	return r
}

// Move copies the *[]byte pointer into the Bytes structure after removing what was in it, if anything. The value input into this function will be empty afterwards
func (r *Bytes) Move(bytes *Bytes) *Bytes {
	if r == nil {
		r = NewBytes()
	}
	if bytes == nil {
		r.err = errors.New("nil parameter")
	} else {
		r.val, r.set, r.err = bytes.val, true, nil
		bytes.val, bytes.set, bytes.err = nil, false, nil
	}
	return r
}

// New nulls the Bytes and assigns an empty *[]byte with a specified size.
func (r *Bytes) New(size int) *Bytes {
	if r == nil {
		r = NewBytes()
	}
	r.Null()
	rr := make([]byte, size)
	r.val = &rr
	r.set = true
	return r
}

// Null wipes the value stored, and restores the Bytes to the same state as a newly created one (with a nil *[]byte).
func (r *Bytes) Null() *Bytes {
	return NewBytes(r)
}

// Rand loads a cryptographically random string of []byte of a specified size.
func (r *Bytes) Rand(size int) *Bytes {
	r = r.New(size)
	rand.Read(*r.Buf())
	r.set = true
	return r
}

// SetError sets the error string
func (r *Bytes) SetError(s string) *Bytes {
	if r == nil {
		r = NewBytes(nil)
	}
	r.err = errors.New(s)
	return r
}
