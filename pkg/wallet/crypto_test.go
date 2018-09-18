package wallet

import (
	"fmt"
	"github.com/awnumar/memguard"
	"testing"
)

func TestCrypto(t *testing.T) {
	db, err := NewDB(f)
	if err != nil {
		t.Error("Failed to open")
	}
	db.Net = "mainnet"
	for i := 0; i < Flast; i++ {
		db.NewTable(KeyNames[i])
	}
	var pass *memguard.LockedBuffer
	pass, err = NewBufferFromBytes([]byte(passwd))
	if err != nil {
		t.Error("failed allocate secure buffer for password", err)
	}
	var es *EncryptedStore
	es, err = Import(pass)
	if err != nil {
		t.Error("failed to import wallet", err)
	}
	fmt.Println(es)
}
