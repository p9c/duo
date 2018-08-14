package wallet

import (
	"fmt"
	"github.com/awnumar/memguard"
	"gitlab.com/parallelcoin/duo/pkg/logger"
	"testing"
)

var (
	f = "/tmp/wallet"
)

func TestNewDB(t *testing.T) {
	db, err := NewDB(f)
	if err != nil {
		t.Error("Failed to open")
	}
	logger.Debug(*db)
	for i := range KeyNames {
		db.NewTable(KeyNames[i])
	}
	db.Close()
}
func TestImport(t *testing.T) {
	db, err := NewDB(f)
	if err != nil {
		t.Error("Failed to open")
	}
	db.Net = "mainnet"
	for i := 0; i < Flast; i++ {
		db.NewTable(KeyNames[i])
	}
	pass, err := memguard.NewImmutableFromBytes([]byte(passwd))
	imp := Import(pass)
	if err != nil {
		t.Error("failed to import wallet", err)
	}
	es := imp.ToEncryptedStore()
	for i := range es.AddressBook {
		fmt.Println(ToJSON(es.AddressBook[i]))
	}
	for i := range es.AddressBook {
		fmt.Println(ToJSON(es.AddressBook[i].Decrypt()))
		es.AddressBook[i].Wipe()
	}
	db.Close()
}

func TestJSON(t *testing.T) {
	es := new(EncryptedStore)
	fmt.Println(string(ToJSON(es)))
}
