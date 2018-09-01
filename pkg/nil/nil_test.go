package password

import (
	"testing"
)

type Iface interface {
	Nil() *Typename
}

func Nil(r *Typename) *Typename {
	if r == nil {
		r = new(Typename)
	}
	if r.Member == nil {
		r.Member = new(SubType)
	}
	return r
}

type SubType struct {
	token int
}

type Typename struct {
	Member *SubType
}

func (r *Typename) FuncName() {
	Nil(r).Member.token = 10
}

func TestNil(t *testing.T) {
	var n Typename
	n.FuncName()
}
