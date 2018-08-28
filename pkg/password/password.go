// Package password is a LockedBuffer with string conversion functions, which are not recommended to be used. However the distinct type adds other safety benefits for the programmer and can be treated exactly the same as the LockedBuffer otherwise.
package password

import (
	. "gitlab.com/parallelcoin/duo/pkg/lockedbuffer"
)

// Password is just a LockedBuffer but recommended for use in storing password fields in other structures to signify what kind of data it is.
type Password struct {
	*LockedBuffer
}
type passwordI interface {
	ToString() string
	FromString(string) *Password
}

// NewPassword creates a new Password
func NewPassword() *Password {
	return new(Password)
}

// ToString returns the password as a string. Not recommended, as the memory is immutable and may end up being copied several times.
func (r *Password) ToString() *string {
	if r == nil || r.LockedBuffer == nil {
		s := ""
		return &s
	}
	s := r.Buf()
	S := string(*s)
	return &S
}

// FromString loads the Lockedbuffer with the bytes of a string. The string is immutable so it is not removed from memory except automatically.
func (r *Password) FromString(s *string) *Password {
	if r == nil {
		r = new(Password)
		r.LockedBuffer = new(LockedBuffer)
		r.LockedBuffer = r.LockedBuffer.NilGuard(r.LockedBuffer, NullLockedBuffer).(*LockedBuffer)
		r.SetError("receiver was nil")
	}
	if r.LockedBuffer == nil {
		r.LockedBuffer = new(LockedBuffer)
		r.SetError("lockedbuffer was nil")
	}
	rr, S := r.New(len(*s)), []byte(*s)
	R := *rr.Buf()
	for i := range S {
		R[i] = S[i]
	}
	r.LockedBuffer = rr
	return r
}
