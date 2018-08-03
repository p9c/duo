// Package walletdat is a library for reading all of the data out of a standard bitcoin berkeleydb wallet.dat file
package walletdat

import (
	"gitlab.com/parallelcoin/duo/pkg/bdb"
	"errors"
	"time"
	"gitlab.com/parallelcoin/duo/pkg/key"
	// "gitlab.com/parallelcoin/duo/pkg/server/args"
	// "gitlab.com/parallelcoin/duo/pkg/cmd"
	// "gitlab.com/parallelcoin/duo/pkg/bdb"
)

type Name struct {
	Addr string
	Name string
}
type Metadata struct {
	Pub        *key.Pub
	Version    uint32
	CreateTime time.Time
}
type Key struct {
	Pub  *key.Pub
	Priv *key.Priv
}
type WKey struct {
	Pub         *key.Pub
	Priv     *key.Priv
	TimeCreated int64
	TimeExpires int64
	Comment     string
}
type MKey struct {
	MKeyID int64
	EncryptedKey              []byte
	Salt                      []byte
	Method          uint32
	Iterations          uint32
	Other []byte
}
type CKey struct {
	Pub  *key.Pub
	Priv []byte
}
// Imports is a list of structures from each of the types of imported records in a wallet.dat
type Imports struct {
	Names []Name
	Metadata []Metadata
	Keys []Key
	WKeys []WKey
	MKeys []MKey
	CKeys []CKey
}

// Import reads an existing wallet.dat and returns all the keys and address data in it
func Import(filename ...string) (imp *Imports, err error) {
	var db *DB
	if db, err = NewDB(); err != nil {
		return nil, errors.New(err.Error() + " unable to open wallet")
		} else if Db.Filename == "" {
			Db.SetFilename("~/.parallelcoin/wallet.dat")
		} else if err = db.Open(); err != nil {
			return
		} else if cursor, err := db.Cursor(bdb.NoTransaction); err != nil {
			return nil, err 
		} else {
			rec := [2][]byte{}
			err = cursor.First(&rec)
			if err != nil {
				return nil, err
			} else {
				for {
					if res := KVDec(rec); res == nil {
						r := res.([]interface{})
						t := r[0].(string)
						switch t {
						case "name":
							imp.Names = append(imp.Names, Name{})
							l := len(imp.Names)-1
							imp.Names[l].Addr = r[1].(string)
							imp.Names[l].Name = r[2].(string)
						case "keymeta":
							imp.Metadata = append(imp.Metadata, Metadata{})
							l := len(imp.Metadata)-1
							imp.Metadata[l] = Metadata{
								Pub: r[1].(*key.Pub),
								Version: r[2].(uint32),
								CreateTime: time.Unix(r[3].(int64), 0)}
						case "key":
							imp.Keys = append(imp.Keys, Key{})
							l := len(imp.Keys)-1
							imp.Keys[l] = Key{
								Pub:  r[1].(*key.Pub),
								Priv: r[2].(*key.Priv), 
							}
						case "wkey":
							imp.WKeys = append(imp.WKeys, WKey{})
							l := len(imp.WKeys)-1
							imp.WKeys[l] = WKey{
								Pub:  r[1].(*key.Pub),
								Priv: r[2].(*key.Priv),
								TimeCreated: r[3].(int64),
								TimeExpires: r[4].(int64),
								Comment:     r[5].(string),
							}
						case "mkey":
							imp.MKeys = append(imp.MKeys, MKey{})
							l := len(imp.WKeys)-1
							imp.MKeys[l] = MKey{
								MKeyID: r[1].(int64),
								EncryptedKey: r[2].([]byte),
								Salt: r[3].([]byte),
								Method: r[4].(uint32),
								Iterations: r[5].(uint32),
								Other: r[6].([]byte),
							}
						case "ckey":
							imp.CKeys = append(imp.CKeys, CKey{})
							l := len(imp.CKeys)-1
							imp.CKeys[l] = CKey{
								Pub: r[2].(*key.Pub),
								Priv: r[3].([]byte),
							}
						default:
						}
					} else {
						if err = cursor.Next(&rec); err != nil {
							err = cursor.Close()
							break
						}
					}
				}
			}
		}
	return
}
