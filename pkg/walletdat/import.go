// Package walletdat is a library for reading all of the data out of a standard bitcoin berkeleydb wallet.dat file
package walletdat

import (
	"errors"

	"github.com/golang/go/src/pkg/time"
	"gitlab.com/parallelcoin/duo/pkg/key"
	"gitlab.com/parallelcoin/duo/pkg/server/args"
	// "gitlab.com/parallelcoin/duo/pkg/cmd"
	// "gitlab.com/parallelcoin/duo/pkg/bdb"
)

// Imports is a list of structures from each of the types of imported records in a wallet.dat
type Imports struct {
	Name []struct {
		Addr []byte
		Name string
	}
	Metadata []struct {
		Pub        *key.Pub
		Version    uint32
		CreateTime time.Time
	}
	Key []struct {
		Pub  *key.Pub
		Priv *key.Priv
	}
	WKey []struct {
		Pub         *key.Pub
		PrivKey     *key.Priv
		TimeCreated int64
		TimeExpires int64
		Comment     string
	}
	Mkey []struct {
		EncryptedKey              []byte
		Salt                      []byte
		DerivationMethod          uint32
		DeriveIterations          uint32
		OtherDerivationParameters []byte
	}
	Ckey []struct {
		Pub  *key.Pub
		Priv []byte
	}
}

// SetFilename changes the name of the database we want to open
func (db *DB) SetFilename(filename string) {
	db.Filename = filename
}

// Import reads an existing wallet.dat and returns all the keys and address data in it
func Import(filename string) (imp *Imports, err error) {
	if Db.Filename == "" {
		Db.SetFilename(*args.DataDir + "/" + *args.Wallet)
	}
	if db, err := NewDB(); err != nil {
		return nil, errors.New(err, "unable to open wallet")
	} else if err := db.Open(); err != nil {
		return
	} else {
		r += dump + "\n"
	}

	return
}
