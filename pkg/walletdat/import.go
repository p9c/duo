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
	}
}

func Import(filename string) (imp Imports) {
	return
}