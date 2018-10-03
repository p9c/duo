package key

import (
	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/hash160"
)

// NewID creates an ID out of the bytes of an object, used for address (public key ID) and script ID
func NewID(bytes *[]byte) (out core.Address) {
	out = core.Address(*hash160.Sum(bytes))
	return
}
