// Package lockedbuffer is a wrapper around the memguard LockedBuffer that automatically handles destroying data no longer needed and enables copy, link and move functions on the data contained inside the structure.
package lockedbuffer

import (
	"errors"
	"github.com/awnumar/memguard"
	. "gitlab.com/parallelcoin/duo/pkg/pipe"
)

// LockedBuffer is a struct that stores and manages memguard.LockedBuffers, ensuring that buffers are destroyed when no longer in use.
//
// Do not use struct literals and not assign them to a name and null() (deletes and zeroes struct) afterwards, or you could run out of memguard LockedBuffers
//
// All functions except for those exporting buffers will automatically allocate the struct of the receiver if it is nil. This permits the use of struct literals in assignment for one-liners that initialise values and call a function with data. It may introduce side effects in code if you did not intend to create a new variable.
//
// The maximum size of buffer is around 172500 bytes on a linux 4.18, it may be more may be less.
type LockedBuffer struct {
	*Pipe
	val *memguard.LockedBuffer
	set bool
	err error
}
type lockedBuffer interface {
	Error() error
	SetError(string)
	Null() *LockedBuffer
	Len() int
	Rand(int) *LockedBuffer
	New(int) *LockedBuffer
	Buf() []byte
	Load(*[]byte) *LockedBuffer
	Copy(*LockedBuffer) *LockedBuffer
	Link(*LockedBuffer) *LockedBuffer
	Move(*LockedBuffer) *LockedBuffer
}

// NewLockedBuffer creates a new, empty LockedBuffer
func NewLockedBuffer() *LockedBuffer {
	return new(LockedBuffer)
}

// Null zeroes out a LockedBuffer
func (r *LockedBuffer) Null() *LockedBuffer {
	return NullLockedBuffer(r).(*LockedBuffer)
}

// NullLockedBuffer clears the passed LockedBuffer or creates a new one if null
func NullLockedBuffer(R interface{}) interface{} {
	r := R.(*LockedBuffer)
	if r == nil {
		r = new(LockedBuffer)
		r.SetError("receiver was nil")
	}
	if r.set {
		r.val.Destroy()
	}
	r.val = nil
	r.set = false
	r.err = nil
	return r
}

// Error returns the string in the err field
func (r *LockedBuffer) Error() error {
	if r.err == nil {
		return nil
	}
	return r.err
}

// SetError sets the string of the error in the err field
func (r *LockedBuffer) SetError(s string) *LockedBuffer {

	return r
}

// Len returns the length of the LockedBuffer if it has been loaded, or zero if not
func (r *LockedBuffer) Len() int {
	if r == nil {
		return 0
	}
	if r.set {
		return r.val.Size()
	}
	return 0
}

// Rand loads the LockedBuffer with cryptographically random bytes to a specified length, destroying existing buffer if it was set
func (r *LockedBuffer) Rand(size int) *LockedBuffer {
	if r == nil {
		r = new(LockedBuffer).NilGuard(r, NullLockedBuffer).(*LockedBuffer)
	}
	r.val, r.err = memguard.NewMutableRandom(size)
	if r.err != nil {
		return r
	}
	r.set = true
	return r
}

// New loads a fresh, zero-filled LockedBuffer in and destroys the existing buffer if it was set
func (r *LockedBuffer) New(size int) *LockedBuffer {
	if r == nil {
		r = new(LockedBuffer).NilGuard(r, NullLockedBuffer).(*LockedBuffer)
	}
	r.val, r.err = memguard.NewMutable(size)
	if r.err != nil {
		return r
	}
	r.set = true
	return r
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

// Load moves the contents of a byte slice into the LockedBuffer, erasing the original copy.
func (r *LockedBuffer) Load(bytes *[]byte) *LockedBuffer {
	if r == nil {
		r = new(LockedBuffer).NilGuard(r, NullLockedBuffer).(*LockedBuffer)
	}
	if bytes == nil {
		return NullLockedBuffer(r).(*LockedBuffer)
	}
	r.val, r.err = memguard.NewMutableFromBytes(*bytes)
	if r.err != nil {
		return r
	}
	r.set = true
	return r
}

// Copy duplicates the buffer from another LockedBuffer.
func (r *LockedBuffer) Copy(buf *LockedBuffer) *LockedBuffer {
	if r == nil {
		r = new(LockedBuffer).NilGuard(r, NullLockedBuffer).(*LockedBuffer)
	}
	if r == buf {
		r.err = errors.New("cannot copy, returning same as given")
		return r
	}
	if buf == nil {
		r.err = errors.New("nil pointer received")
		return r
	}
	// This function cannot return an error because len() cannot return negative
	r.New(buf.Len())
	A := *r.Buf()
	b := buf.Buf()
	for i := range A {
		A[i] = (*b)[i]
	}
	r.set = true
	r.err = nil
	return r
}

// Link copies the pointer from another LockedBuffer's content, meaning what is written to one will also be visible in the other
func (r *LockedBuffer) Link(buf *LockedBuffer) *LockedBuffer {
	if r == nil {
		r = new(LockedBuffer).NilGuard(r, NullLockedBuffer).(*LockedBuffer)
	}
	r.val = buf.val
	r.set = true
	r.err = nil
	return r
}

// Move copies the pointer to the buffer into the receiver and nulls the passed LockedBuffer
func (r *LockedBuffer) Move(buf *LockedBuffer) *LockedBuffer {
	if r == nil {
		r = new(LockedBuffer).NilGuard(r, NullLockedBuffer).(*LockedBuffer)
	}
	r.val = buf.val
	r.set = true
	r.err = nil
	buf.val = nil
	buf.set = false
	return r
}
