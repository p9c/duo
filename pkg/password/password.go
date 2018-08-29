// Package password is a LockedBuffer with string conversion functions, which are not recommended to be used. However the distinct type adds other safety benefits for the programmer and can be treated exactly the same as the LockedBuffer otherwise.
package password

import (
	"encoding/json"
	. "gitlab.com/parallelcoin/duo/pkg/bytes"
	. "gitlab.com/parallelcoin/duo/pkg/lockedbuffer"
)

// Password is just a LockedBuffer but recommended for use in storing password fields in other structures to signify what kind of data it is.
type Password struct {
	*LockedBuffer
}

// NewPassword creates a new Password
func NewPassword(r ...*Password) *Password {
	if len(r) == 0 {
		r = append(r, new(Password))
		r[0].LockedBuffer = NewLockedBuffer()
	}
	if r[0] == nil {
		r[0] = new(Password)
		r[0].LockedBuffer = NewLockedBuffer(r[0].LockedBuffer)
	}
	r[0].SetUTF8()
	return r[0]
}

type passwordI interface {
	ToString() string
	FromString(string) *Password
	MarshalJSON() ([]byte, error)
	String() string
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
func (r *Password) FromString(s string) *Password {
	if r == nil {
		r = NewPassword(r)
	}
	S := []byte(s)
	r.LockedBuffer = r.New(len(s))
	R := *r.Buf()
	for i := range S {
		R[i] = S[i]
	}
	return r
}

// MarshalJSON marshalls the JSON for a Password
func (r *Password) MarshalJSON() ([]byte, error) {
	if r == nil {
		r = NewPassword()
	}
	if r.LockedBuffer == nil {
		r.LockedBuffer = NewLockedBuffer()
	}
	if r.LockedBuffer.Len() == 0 {
		r.LockedBuffer.Load(NewBytes().Rand(32).Buf())
	}
	return json.Marshal(&struct {
		Value  string
		IsSet  bool
		IsUTF8 bool
		Error  string
	}{
		Value:  string(*r.LockedBuffer.Buf()),
		IsSet:  r.IsSet(),
		IsUTF8: r.IsUTF8(),
		Error:  r.Error(),
	})
}
