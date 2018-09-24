package buf

import (
	"github.com/awnumar/memguard"
	"github.com/parallelcointeam/duo/pkg/proto"
)

var er = proto.Errors

// Byte is a simple byte slice
type Byte struct {
	Val *[]byte
	*proto.State
	Coding string
}

// Secure is a memguard LockedBuffer
type Secure struct {
	Val    *memguard.LockedBuffer
	Status string
	Coding string
}
