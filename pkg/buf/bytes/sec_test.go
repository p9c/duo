package buf

import (
	"testing"
)

func TestSecBuf(t *testing.T) {
	n := NewFenced()
	n.Buf()
	r := NewFenced().Copy(NewFenced().Rand(12))
	r.Buf()
	var empty *Fenced
	empty.Buf()
	empty.Rand()
}
