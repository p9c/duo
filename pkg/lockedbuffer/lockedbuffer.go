// Package lb is a wrapper around the native byte slice that automatically handles purging discarded data and enables copy, link and move functions on the data contained inside the structure.
package lb

import (
	"errors"
	"github.com/awnumar/memguard"
)

// LockedBuffer is a struct that stores and manages memguard.LockedBuffers, ensuring that buffers are destroyed when no longer in use.
type LockedBuffer struct {
	val *memguard.LockedBuffer
	set bool
	err error
}

type lockedBuffer interface {
	Len() int
	Null() *LockedBuffer
	Rand(int) *LockedBuffer
	New(int) *LockedBuffer
	Buf() []byte
	Load(*[]byte) *LockedBuffer
	Copy(*LockedBuffer) *LockedBuffer
	Link(*LockedBuffer) *LockedBuffer
	Move(*LockedBuffer) *LockedBuffer
}

// Len returns the length of the LockedBuffer if it has been loaded, or zero
func (r *LockedBuffer) Len() int {
	if r.set {
		return r.val.Size()
	}
	return 0
}

// Null destroys the LockedBuffer if it has been set, and nulls all the variables in the LockedBuffer
func (r *LockedBuffer) Null() *LockedBuffer {
	if r.set {
		r.val.Destroy()
	}
	r.val = nil
	r.set = false
	r.err = nil
	return r
}

// Rand loads the LockedBuffer with cryptographically random bytes to a specified length, destroying existing buffer if it was set
func (r *LockedBuffer) Rand(size int) *LockedBuffer {
	r.Null()
	r.val, r.err = memguard.NewMutableRandom(size)
	if r.err != nil {
		return r
	}
	r.set = true
	return r
}

// New loads a fresh, zero-filled LockedBuffer in and destroys the existing buffer if it was set
func (r *LockedBuffer) New(size int) *LockedBuffer {
	r.Null()
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
	if r.set {
		a := r.val.Buffer()
		return &a
	}
	return nil
}

// Load moves the contents of a byte slice into the LockedBuffer, erasing the original copy.s
func (r *LockedBuffer) Load(bytes *[]byte) *LockedBuffer {
	r.Null()
	r.val, r.err = memguard.NewMutableFromBytes(*bytes)
	if r.err != nil {
		return r
	}
	r.set = true
	return r
}

// Copy duplicates the buffer from another LockedBuffer.
func (r *LockedBuffer) Copy(buf *LockedBuffer) *LockedBuffer {
	if buf == nil {
		r.err = errors.New("nil pointer received")
		return r
	}
	r.Null()
	r.New(buf.Len())
	if r.err != nil {
		return r
	}
	a := r.Buf()
	A := *a
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
	r.Null()
	r.val = buf.val
	r.set = true
	r.err = nil
	return r
}

// Move copies the pointer to the buffer into the receiver and nulls the passed LockedBuffer
func (r *LockedBuffer) Move(buf *LockedBuffer) *LockedBuffer {
	r.Null()
	r.val = buf.val
	r.set = true
	r.err = nil
	buf.val = nil
	buf.set = false
	return r
}
