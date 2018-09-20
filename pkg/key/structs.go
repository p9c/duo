package key

import (
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/crypt"
	"github.com/parallelcointeam/duo/pkg/proto"
)

var er = proto.Errors

// Priv is a private key, stored in a Crypt
type Priv struct {
	crypt.Crypt
	pub   *Pub
	valid bool
}

// Pub is a secp256k1 EC public key which can be represented as a compressed, uncompressed or hybrid for wire and storage
type Pub struct {
	buf.Byte
}
