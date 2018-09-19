package buf

import (
	"github.com/awnumar/memguard"
)

// Byte is a simple byte slice
type Byte struct {
	Val    *[]byte
	Status string
	Coding string
}

// Secure is a memguard LockedBuffer
type Secure struct {
	Val    *memguard.LockedBuffer
	Status string
	Coding string
}
