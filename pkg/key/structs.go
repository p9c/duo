package key

import (
	"github.com/parallelcointeam/duo/pkg/blockcrypt"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/crypt"
)

var er = core.Errors

// Priv is a private key, stored in a Crypt
type Priv struct {
	crypt.Crypt
	pub   *Pub
	valid bool
	core.State
}

// Pub is a secp256k1 EC public key which can be represented as a compressed, uncompressed or hybrid for wire and storage
type Pub struct {
	buf.Byte
}

// Sig is a bitcoin EC signature
type Sig struct {
	buf.Byte
	mh   *buf.Byte
	addr core.Address
}

// Store is a keychain for public and private keys
type Store struct {
	BC    *bc.BlockCrypt
	privs map[core.Address]*Priv
	pubs  map[core.Address]*Pub
	core.State
}
