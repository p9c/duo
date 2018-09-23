package key

import (
	"github.com/parallelcointeam/duo/pkg/blockcrypt"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/crypt"
	"github.com/parallelcointeam/duo/pkg/proto"
)

var er = proto.Errors

// Priv is a private key, stored in a Crypt
type Priv struct {
	*crypt.Crypt
	pub   *Pub
	valid bool
}

// Pub is a secp256k1 EC public key which can be represented as a compressed, uncompressed or hybrid for wire and storage
type Pub struct {
	*buf.Byte
}

// Sig is a bitcoin EC signature
type Sig struct {
	buf.Byte
	mh   *buf.Byte
	addr proto.Address
}

// Store is a keychain for public and private keys
type Store struct {
	BC     *blockcrypt.BlockCrypt
	privs  map[proto.Address]*Priv
	pubs   map[proto.Address]*Pub
	Status string
}
