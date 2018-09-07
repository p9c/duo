package buf

import (
	"testing"
)

func TestUnsafe(t *testing.T) {
	al := &Unsafe{}
	// an := NewUnsafe()
	// var ay *Unsafe
	// ab := make([]byte, 0)
	// ap := &ab
	// az := NewUnsafe().Load(ap).(*Unsafe)
	al.Status().Set("test")
}
