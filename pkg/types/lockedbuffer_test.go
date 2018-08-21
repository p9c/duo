package types

import (
	"fmt"
	"testing"
)

func TestLockedBuffer(t *testing.T) {
	s := "this is a test"
	l := NewLockedBuffer().FromBytes(NewBytes().FromString(&s))
	fmt.Println(l.WithSize(19).ToBytes())
}
