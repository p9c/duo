// Package password is a LockedBuffer with string conversion functions, which are not recommended to be used. However the distinct type adds other safety benefits for the programmer and can be treated exactly the same as the LockedBuffer otherwise.
package password

import (
	"encoding/json"
	. "gitlab.com/parallelcoin/duo/pkg/byte"
	. "gitlab.com/parallelcoin/duo/pkg/bytes"
	. "gitlab.com/parallelcoin/duo/pkg/interfaces"
	. "gitlab.com/parallelcoin/duo/pkg/lockedbuffer"
)

// Password is just a LockedBuffer but recommended for use in storing password fields in other structures to signify what kind of data it is.
type Password struct {
	*LockedBuffer
}

// guards against nil pointer receivers
func ifnil(r *Password) *Password {
	if r == nil {
		r = new(Password)
		r.LockedBuffer = new(LockedBuffer)
		r.SetError("nil receiver")
	}
	return r
}

// NewPassword creates a new Password
func NewPassword(r ...*Password) *Password {
	if len(r) == 0 {
		r = append(r, new(Password))
		r[0].LockedBuffer = NewLockedBuffer()
	}
	if r[0] == nil {
		r[0] = new(Password)
		r[0].LockedBuffer = NewLockedBuffer()
	}
	r[0].SetCoding("string")
	return r[0]
}

type passwordI interface {
	FromString(string) *Password
	MarshalJSON() ([]byte, error)
	String() string
}

// FromString loads the Lockedbuffer with the bytes of a string. The string is immutable so it is not removed from memory except automatically.
func (r *Password) FromString(s string) *Password {
	if r == nil {
		r = NewPassword(r)
	}
	r.LockedBuffer = r.New(len(s)).(*LockedBuffer)
	r.ForEach(func(i int) {
		r.SetElem(i, NewByte().Load(&[]byte{s[i]}))
	})
	// for i := range S {
	// 	R[i] = S[i]
	// }
	return r
}

// MarshalJSON marshalls the JSON for a Password
func (r *Password) MarshalJSON() ([]byte, error) {
	r = ifnil(r)
	if r.LockedBuffer == nil {
		r.LockedBuffer = NewLockedBuffer()
	}
	if r.LockedBuffer.Len() == 0 {
		r.LockedBuffer.Load(NewBytes().Rand(32).Buf())
	}
	return json.Marshal(&struct {
		Value  string
		IsSet  bool
		Coding string
		Error  string
	}{
		Value:  string(*r.LockedBuffer.Buf()),
		IsSet:  r.IsSet(),
		Coding: r.Coding(),
		Error:  r.Error(),
	})
}

/////////////////////////////////////////
// Stringer implementation
/////////////////////////////////////////

// String returns value encoded according to the current coding mode
func (r *Password) String() string {
	r = ifnil(r)
	return r.LockedBuffer.String()
}
