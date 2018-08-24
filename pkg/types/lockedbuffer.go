package types

import (
	"fmt"
	"github.com/awnumar/memguard"
)

// LockedBuffer uses the memguard library to create a memory buffer that protects against reading or writing of the data by other processes running on a machine (except root user's processes of course).
// Because of the need to manually free the buffers since the system doesn't really allow very much memory to have fences around it, there may be cases where you have to manually delete the buffers so they don't run out. If this becomes troublesome we may add tracking via a 'factory' type, to assist debugging out-of-memory errors.
type LockedBuffer struct {
	value *memguard.LockedBuffer
	err   error
}
type lockedBuffer interface {
	Buffer() *[]byte
	Copy() *[]byte
	Len() int
	WithSize(int) *LockedBuffer
	ToLockedBuffer() *LockedBuffer
	FromLockedBuffer(*LockedBuffer) *LockedBuffer
	FromRandomBytes(int) *LockedBuffer
	Delete()
	ToBytes() *Bytes
	FromBytes(*Bytes) *LockedBuffer
}

// -----------------------------------------------------------------------------
// Methods related to the LockedBuffer

// NewLockedBuffer creates an unpopulated LockedBuffer structure
func NewLockedBuffer() *LockedBuffer {
	return new(LockedBuffer)
}

// Buffer returns a pointer to the buffer storing the value inside the Bytes
func (lb *LockedBuffer) Buffer() (B *[]byte) {
	b := lb.value.Buffer()
	return &b
}

// Copy makes a copy and returns it.
// WARNING: This buffer will need to be deleted later, so don't use this method without assigning it to a variable (it should always be the last in a method call chain and part of an assignment or a parameter in another function).
func (lb *LockedBuffer) Copy() (LB *LockedBuffer) {
	LB = NewLockedBuffer().WithSize(lb.Len())
	b := lb.value.Buffer()
	B := LB.value.Buffer()
	for i := range b {
		B[i] = b[i]
	}
	return
}

// Len returns the length of the data stored in a LockedBuffer
func (lb *LockedBuffer) Len() int {
	if lb.value == nil {
		return 0
	}
	return lb.value.Size()
}

// WithSize allocates a memguard LockedBuffer of a given size
func (lb *LockedBuffer) WithSize(size int) *LockedBuffer {
	if lb.value != nil {
		lb.value.Destroy()
	}
	lb.value, lb.err = memguard.NewMutable(size)
	return lb
}

// ToLockedBuffer moves the current buffer into a new one
func (lb *LockedBuffer) ToLockedBuffer() (LB *LockedBuffer) {
	LB = NewLockedBuffer()
	LB.value, LB.err = memguard.NewMutableFromBytes(*lb.Buffer())
	lb.Delete()
	return
}

// FromLockedBuffer copies in the reference to another LockedBuffer
func (lb *LockedBuffer) FromLockedBuffer(in *LockedBuffer) *LockedBuffer {
	if lb.value != nil {
		lb.value.Destroy()
	}
	if in.value != nil {
		lb.value = in.value
		in.value = nil
	}
	return lb
}

// FromRandomBytes creates a new LockedBuffer from a cryptographically secure source
func (lb *LockedBuffer) FromRandomBytes(size int) *LockedBuffer {
	if lb.value != nil {
		lb.value.Destroy()
	}
	lb.value, lb.err = memguard.NewMutableRandom(size)
	return lb
}

// Delete destroys the locked buffer, erasing its contents and deallocating
func (lb *LockedBuffer) Delete() {
	lb.value.Destroy()
}

// -----------------------------------------------------------------------------
// Methods relating to Bytes

// ToBytes moves a LockedBuffer value into a Bytes
func (lb *LockedBuffer) ToBytes() (B *Bytes) {
	if lb.value == nil {
		return nil
	}
	b := lb.value.Buffer()
	B = NewBytes().FromByteSlice(&b)
	return
}

// FromBytes copies the content of a Bytes into a LockedBuffer
func (lb *LockedBuffer) FromBytes(in *Bytes) *LockedBuffer {
	if lb.value != nil {
		lb.Delete()
	}
	b := in.ToByteSlice()
	lb.value, lb.err = memguard.NewMutableFromBytes(*b)
	return lb
}

func init() {
	fmt.Println("pkg/types/lockedbuffer.go initialising")
}

// ----------------------------------------------------------------------------
// Methods relating to builtin byte slices

// ToByteSlice moves the byte slice inside it out
func (lb *LockedBuffer) ToByteSlice() (out *[]byte) {
	L := lb.value.Buffer()
	o := make([]byte, len(L))
	for i := range L {
		o[i] = L[i]
	}
	out = &o
	lb.Delete()
	return
}

// FromByteSlice moves the content of a byte slice into the LockedBuffer
func (lb *LockedBuffer) FromByteSlice(in *[]byte) *LockedBuffer {
	if lb.value != nil {
		lb.Delete()
	}
	return NewLockedBuffer().FromBytes(NewBytes().FromByteSlice(in))
}
