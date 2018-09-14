package proto

import (
	"github.com/awnumar/memguard"
)

// Byte is a simple single byte
type Byte struct {
	Val    *byte
	Status string
	Coding string
}

// Bytes is a simple byte slice
type Bytes struct {
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
