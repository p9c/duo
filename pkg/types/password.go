package types

import (
	"fmt"
)

// Password is a LockedBuffer with string conversion functions added
type Password struct {
	*LockedBuffer
}
type password interface {
	Copy() *Password
	Buffer() *Bytes
	ToString() string
	FromString(string) *LockedBuffer
	ToPassword() *Password
	FromPassword(*Password) *Password
}

// NewPassword creates a new empty password object
func NewPassword() (p *Password) {
	p = new(Password)
	p.LockedBuffer = NewLockedBuffer()
	return p
}

// Copy makes a new Password with a copy of the contents of the receiver
func (p *Password) Copy() (P *Password) {
	P = NewPassword()
	P.LockedBuffer = p.LockedBuffer.Copy()
	return
}

// Buffer returns a Bytes containing a link to the buffer inside the Password struct
func (p *Password) Buffer() (B *[]byte) {
	b := p.value.Buffer()
	B = &b
	return
}

// ToString copies the content of the Password buffer into a string.
// WARNING: Go strings are immutable and potentially could be copied many times so be careful!
func (p *Password) ToString() (S *string) {
	S = p.ToBytes().ToString()
	return
}

// FromString copies a string into the LockedBuffer.
// WARNING: Go strings are immutable so be aware that this password will persist and potentially be copied several times before being zeroed again for a new allocation.
func (p *Password) FromString(s *string) *Password {
	b := []byte(*s)
	p.FromBytes(NewBytes().FromByteSlice(&b))
	return p
}

// ToPassword moves the contents of the receiver into a new Password, dereferences and returns the new structure.
// WARNING: This effectively destroys the receiver.
func (p *Password) ToPassword() (P *Password) {
	P = new(Password)
	P.LockedBuffer = p.LockedBuffer
	p.LockedBuffer = nil
	return
}

// FromPassword moves the contents of the parameter into the receiver.
// WARNING: The parameter will become empty as a result!
// This is to avoid the possibility of the variable being modified by two separate goroutines. If you want the original to remain append a .Copy() to the parameter.
func (p *Password) FromPassword(P *Password) *Password {
	p = new(Password)
	p.LockedBuffer = P.LockedBuffer
	P.LockedBuffer = nil
	return p
}

func init() {
	fmt.Println("pkg/types/password.go initialising")
}
