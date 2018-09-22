package key

import (
	"github.com/parallelcointeam/duo/pkg/hash160"
	"github.com/parallelcointeam/duo/pkg/proto"
)

// NewID creates an ID out of the bytes of an object, used for address (public key ID) and script ID
func NewID(bytes *[]byte) (out proto.ID) {
	out = proto.ID(*hash160.Sum(bytes))
	return
}
