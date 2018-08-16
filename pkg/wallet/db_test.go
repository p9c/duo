package wallet

import (
	"fmt"
	"os"
	"reflect"
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
	pass, err := NewBufferFromBytes([]byte(passwd))
	es := Import(pass)
	if err != nil {
		t.Error("failed to import wallet", err)
	}
	v := reflect.ValueOf(es)
	for i := 0; i < v.NumField(); i++ {
		fmt.Println(v.Field(i).Interface())
	}
	db.Close()
}

func TestJSON(t *testing.T) {
	es := new(EncryptedStore)
	fmt.Println(string(ToJSON(es)))
}
