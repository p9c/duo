package wallet

import (
	"fmt"
	"github.com/awnumar/memguard"
	"os"
	"testing"
)

var (
	f = "/tmp/wallet"
)

func TestNewDB(t *testing.T) {
	os.RemoveAll(f)
	db, err := NewDB(f)
	if err != nil {
		t.Error("Failed to open")
	}
	for i := range KeyNames {
		db.NewTable(KeyNames[i])
	}
	db.Close()
}

func TestImport(t *testing.T) {
	os.RemoveAll(f)
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
	var es *EncryptedStore
	es, err = Import(pass)
	if err != nil {
		t.Error("failed to import wallet", err)
	}
	fmt.Println(string(ToJSON(es)))
	db.Close()
}

func TestJSON(t *testing.T) {
	es := new(EncryptedStore)
	fmt.Println(string(ToJSON(es)))
}
