package key

import (
	"github.com/parallelcointeam/duo/pkg/buf"
)

// ID is the RipeMD160 hash of a public key
type ID struct {
	buf.Byte
}
