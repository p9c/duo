package types

import (
	"fmt"
)

// Password is a LockedBuffer with string conversion functions added
type Password struct {
	*LockedBuffer
}
type password interface {
	ToString() string
	FromString(string) *LockedBuffer
}

// NewPassword creates a new empty password object
func NewPassword() (p *Password) {
	p = new(Password)
	p.LockedBuffer = NewLockedBuffer()
	return p
}

// ToString copies the content of the Password buffer into a string.
// WARNING: Go strings are immutable and potentially could be copied many times so be careful!
func (p *Password) ToString() string {
	return string(p.value.Buffer())
}

// FromString copies a string into the LockedBuffer.
// WARNING: Go strings are immutable so be aware that this password will persist and potentially be copied several times before being zeroed again for a new allocation.
func (p *Password) FromString(s *string) *Password {
	b := []byte(*s)
	p.FromBytes(NewBytes().FromByteSlice(&b))
	return p
}

func init() {
	fmt.Println("pkg/types/password.go initialising")
}
