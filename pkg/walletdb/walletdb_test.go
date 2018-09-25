package walletdb

import (
	"encoding/hex"
	"fmt"
	"github.com/parallelcointeam/duo/pkg/blockcrypt"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/hash160"
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

func TestMasterKey(t *testing.T) {
	p := []byte("testing password")
	pass := buf.NewSecure().Copy(&p).(*buf.Secure)
	BC := bc.New().Generate(pass).Arm()
	teststring := "if you can read this the encryption worked"
	testbytes := []byte(teststring)
	testcipher := BC.Encrypt(&testbytes)
	fmt.Println("string '" + teststring + "'")
	fmt.Println("plain", hex.EncodeToString(testbytes))
	fmt.Println("ciph ", hex.EncodeToString(*testcipher))
	wdb := NewWalletDB()
	if wdb.OK() {
		defer wdb.Close()
	}
	wdb.WriteMasterKey(BC)
	BCs := wdb.ReadMasterKeys()
	for i := range BCs {
		crypt := BCs[i].Crypt.Bytes()
		iv := BCs[i].IV.Bytes()
		iterations := BCs[i].Iterations
		idx := proto.Hash64(crypt)
		fmt.Println("idx  ", len(*idx), hex.EncodeToString(*idx))
		fmt.Println("crypt", len(*crypt), hex.EncodeToString(*crypt))
		fmt.Println("iv   ", len(*iv), hex.EncodeToString(*iv))
		fmt.Println("iters", iterations)
		BCs[i].Unlock(buf.NewSecure().Copy(&p).(*buf.Secure)).Arm()
		fmt.Println("ciph ", len(*testcipher), hex.EncodeToString(*testcipher))
		fmt.Println("decrp", len(*BCs[i].Decrypt(testcipher)), hex.EncodeToString(*BCs[i].Decrypt(testcipher)))
		fmt.Println("decrs '" + string(*BCs[i].Decrypt(testcipher)) + "'")

		k := key.NewPriv()
		k.WithBC(BCs[i])
		k.Make()
		wdb.WithBC(BCs[i])
		wdb.WriteKey(k)
		kh := hash160.Sum(k.PubKey().Bytes())
		pk := wdb.ReadKey(kh)
		fmt.Println("prvkey", *pk.Bytes())
		fmt.Println("pubkey", *pk.PubKey().Bytes())
		wdb.EraseKey(kh)

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
