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

// Len returns the length of the data stored in a LockedBuffer
func (lb *LockedBuffer) Len() int {
	return lb.value.Size()
}

// WithSize allocates a memguard LockedBuffer of a given size
func (lb *LockedBuffer) WithSize(size int) *LockedBuffer {
	if lb.value != nil {
		lb.value.Destroy()
	}
	lb.value, lb.err = memguard.NewImmutable(size)
	return lb
}

// ToLockedBuffer copies the current buffer into a new one
func (lb *LockedBuffer) ToLockedBuffer() (LB *LockedBuffer) {
	LB = NewLockedBuffer()
	LB.value, LB.err = memguard.NewMutableFromBytes(lb.value.Buffer())
	return
}

// FromLockedBuffer copies in the reference to another LockedBuffer
func (lb *LockedBuffer) FromLockedBuffer(in *LockedBuffer) *LockedBuffer {
	lb.value.Copy(in.value.Buffer())
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

// ToBytes copies a LockedBuffer value into a Bytes
func (lb *LockedBuffer) ToBytes() (B *Bytes) {
	b := lb.value.Buffer()
	fmt.Println(b)
	B = NewBytes().FromByteSlice(&b)
	return
}

// FromBytes copies the content of a Bytes into a LockedBuffer
func (lb *LockedBuffer) FromBytes(in *Bytes) *LockedBuffer {
	if lb.value != nil {
		lb.Delete()
	}
	b := in.ToByteSlice()
	fmt.Println(*b)
	lb.value, lb.err = memguard.NewMutableFromBytes(*b)
	return lb
}

func init() {
	fmt.Println("pkg/types/lockedbuffer.go initialising")
}
