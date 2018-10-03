package buf

import (
	"github.com/awnumar/memguard"
	"github.com/parallelcointeam/duo/pkg/core"
)

var er = core.Errors

// Byte is a simple byte slice
type Byte struct {
	Val *[]byte
	core.State
	Coding string
}

// Secure is a memguard LockedBuffer
type Secure struct {
	Val *memguard.LockedBuffer
	core.State
	Coding string
}
