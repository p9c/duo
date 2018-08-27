// Package b is a wrapper around the native byte slice that automatically handles purging discarded data and enables copy, link and move functions on the data contained inside the structure.
//
// The purpose of this structure is to enable the chaining of pointer methods and eliminate the need for separate assignments by passing error value within the structure instead of as the last term in the return tuple. This structuring has a similar functionality to default parameters, without the compile-time complexity. The same pattern can be used to extend the type to be incorporated into a new aggregate type that can contain more similar structures or add them in addition to implemented methods.
package bytes

import (
	"crypto/rand"
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

type bytes interface {
	Len() int
	Null() *Bytes
	Rand(int) *Bytes
	New(int) *Bytes
	Buf() []byte
	Load(*[]byte) *Bytes
	Copy(*Bytes) *Bytes
	Link(*Bytes) *Bytes
	Move(*Bytes) *Bytes
}

// NewBytes creates a new empty Bytes
func NewBytes() *Bytes {
	return new(Bytes)
}

// Len returns the length of the *[]byte if it has a value assigned, or zero
func (r *Bytes) Len() int {
	if r == nil {
		return 0
	}
	if r.set {
		return len(*r.val)
	}
	return 0
}

// Null wipes the value stored, and restores the Bytes to the same state as a newly created one (with a nil *[]byte).
func (r *Bytes) Null() *Bytes {
	if r == nil {
		r = new(Bytes)
	}
	if r.set {
		rr := *r.val
		for i := range rr {
			rr[i] = 0
		}
	}
	r.val, r.set, r.err = nil, false, nil
	return r
}

// Rand loads a cryptographically random string of []byte of a specified size.
func (r *Bytes) Rand(size int) *Bytes {
	r = r.New(size)
	// This function simply blocks until it gets all it wants from the random source.
	rand.Read(*r.Buf())
	return r
}

// New nulls the Bytes and assigns an empty *[]byte with a specified size.
func (r *Bytes) New(size int) *Bytes {
	if r == nil {
		r = new(Bytes)
	} else {
		r.Null()
	}
	b := make([]byte, size)
	r.Load(&b)
	return r
}

// Buf returns a variable pointing to the value stored in a Bytes.
func (r *Bytes) Buf() *[]byte {
	if r == nil {
		return &[]byte{}
	}
	return r.val
}

// Load nulls any existing data and sets its pointer to refer to the pointer to byte slice in the parameter.
func (r *Bytes) Load(bytes *[]byte) *Bytes {
	if r == nil {
		r = new(Bytes)
	}
	r.val, r.set, r.err = bytes, true, nil
	return r
}

// Copy duplicates the data from the *[]byte provided and zeroes and replaces its contents, clearing the error value.
func (r *Bytes) Copy(bytes *Bytes) *Bytes {
	if r == nil {
		r = new(Bytes)
	}
	if bytes == nil {
		return r
	}
	b := bytes.Buf()
	if b == nil {
		return r
	}
	temp := make([]byte, bytes.Len())
	for i := range *b {
		temp[i] = (*b)[i]
	}
	return r.Load(&temp)
}

// Link nulls the Bytes and copies the pointer in from another Bytes. Whatever is done to one's []byte will also affect the other, but they keep separate error values
func (r *Bytes) Link(bytes *Bytes) *Bytes {
	if r == nil {
		r = new(Bytes)
	} else {
		r.Null()
	}
	return r.Load(bytes.val)
}

// Move copies the *[]byte pointer into the Bytes structure after removing what was in it, if anything. The value input into this function will be empty afterwards
func (r *Bytes) Move(bytes *Bytes) *Bytes {
	r.Load(bytes.val)
	bytes.Null()
	return r
}
