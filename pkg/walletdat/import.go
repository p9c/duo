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
func Import(filename ...string) (imp *Imports, err error) {
	if db, err := NewDB(); err != nil {
		return nil, errors.New(err, "unable to open wallet")
		} else if Db.Filename == "" {
			Db.SetFilename("~/.parallelcoin/wallet.dat")
		} else if err := db.Open(); err != nil {
			return
		} else if cursor, err := db.Cursor(bdb.NoTransaction); err != nil {
			return "", err 
		} else {
			rec := [2][]byte{}
			err = cursor.First(&rec)
			if err != nil {
				return "", err
			} else {
				dbt, _ := db.Type()
				for {
					if res := KVDec(rec); res == nil {
						r := res.([]interface{})
						t := r[0].(string)
						switch t {
						case "name":
							imp.Name = append(imp.Name, make(imp.Name))
							l := len(imp.Name)-1
							imp.Name[l].Addr = r[1].(string)
							imp.Name[l].Name = r[2].(string)
						case "keymeta":
							imp.Metadata = append(imp.Metadata, make(imp.Metadata))
							l := len(imp.Metadata)-1
							imp.Metadata[l].Pub = r[1].(*key.Pub)
							imp.Metadata[l].Version = r[2].(uint32)
							imp.Metadata[l].CreateTime = r[3].(int64)
						case "key":
							imp.Key = append(imp.Key, make(imp.Key))
							
						case "wkey":
						case "mkey":
						case "ckey":
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
	}
	return
}
