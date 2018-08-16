package wallet

import (
	"fmt"
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
	pass, err := NewBufferFromBytes([]byte(passwd))
	es := Import(pass)
	if err != nil {
		t.Error("failed to import wallet", err)
	}
	fmt.Println(es)
}
