package walletdb

import (
	"fmt"
	"github.com/parallelcointeam/duo/pkg/blockcrypt"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/key"
	"github.com/parallelcointeam/duo/pkg/proto"
	"testing"
)

func TestOpenClose(t *testing.T) {
	db := NewWalletDB()
	if db.OK() {
		defer db.Close()
	}
	fmt.Println(db)
}

func TestWriteMasterKey(t *testing.T) {
	p := []byte("testing password")
	pass := buf.NewSecure().Copy(&p).(*buf.Secure)
	BC := bc.New().Generate(pass)
	wdb := NewWalletDB()
	if wdb.OK() {
		defer wdb.Close()
	}
	wdb.WriteMasterKey(BC)
}

func TestReadMasterKey(t *testing.T) {
	wdb := NewWalletDB()
	if wdb.OK() {
		defer wdb.Close()
	}
	BC := wdb.ReadMasterKeys()
	for i := range BC {
		fmt.Println("crypt", BC[i].Crypt.Bytes())
		fmt.Println("iv", BC[i].IV.Bytes())
		fmt.Println("iterations", BC[i].Iterations)
		idx := proto.Hash64(BC[i].Crypt.Bytes())
		wdb.EraseMasterKey(idx)
	}
}

func TestEncryptDecrypt(t *testing.T) {
	p := []byte("testing password")
	pass := buf.NewSecure().Copy(&p).(*buf.Secure)
	BC := bc.New().Generate(pass)
	wdb := NewWalletDB()
	if wdb.OK() {
		defer wdb.Close()
	}
	wdb.WriteMasterKey(BC)
	k := key.NewPriv()
	k.WithBC(BC)
	k.Make()
	fmt.Println("secret", k.Error())
	// wdb.WriteKey(k)
}
