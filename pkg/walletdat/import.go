// Package walletdat is a library for reading all of the data out of a standard bitcoin berkeleydb wallet.dat file
package walletdat

import (
)

type Imports struct {
	Name []struct {
		Addr []byte
		Name string
	}
	Metadata []struct {
		Pub *key.Pub
		Version uint32
		CreateTime time.Time
	}
	Key []struct {
		Pub *key.Pub
		Priv *key.Priv
	}
	WKey []struct {
		Pub *key.Pub
		PrivKey     *key.Priv
		TimeCreated int64
		TimeExpires: int64
		Comment:     string
	}
	Mkey []struct {
		EncryptedKey []byte
		Salt []byte
		DerivationMethod uint32
		DeriveIterations uint32
		OtherDerivationParameters []byte
	}
	Ckey []struct {
		Pub *key.Pub
		Priv []byte
	}
}

func Import(filename string) (imp Imports) {
	return
}