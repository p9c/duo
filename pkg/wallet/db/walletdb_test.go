package walletdb

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/parallelcointeam/duo/pkg/blockcrypt"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/key"
)

func TestMasterKey(t *testing.T) {
	p := []byte("testing password")
	pass := buf.NewSecure().Copy(&p).(*buf.Secure)
	BC := bc.New().Generate(pass).Arm()
	origcrypt := BC.Crypt.Bytes()
	origiv := BC.IV.Bytes()
	origiters := BC.Iterations

	wdb := NewWalletDB()
	if wdb.OK() {
		defer wdb.Close()
	}
	wdb.WithBC(BC)

	wdb.WriteMasterKey(BC)
	BCs := wdb.ReadMasterKeys()
	for i := range BCs {
		crypt := BCs[i].Crypt.Bytes()
		iv := BCs[i].IV.Bytes()
		iterations := BCs[i].Iterations
		if bytes.Compare(*crypt, *origcrypt) != 0 {
			t.Error("crypt not decrypted properly")
		}
		if bytes.Compare(*iv, *origiv) != 0 {
			t.Error("iv not retrieved properly")
		}
		if iterations != origiters {
			t.Error("iterations not retrieved properly")
		}

		BCs[i].Unlock(buf.NewSecure().Copy(&p).(*buf.Secure)).Arm()

		plaintext := "This is the message in plaintext"
		plainbytes := []byte(plaintext)
		encrypted := BCs[i].Encrypt(&plainbytes)
		decrypted := BCs[i].Decrypt(encrypted)

		if bytes.Compare(plainbytes, *decrypted) != 0 {
			t.Error("encryption/decryption did not work properly")
		}
	}
	wdb.EraseMasterKey(BC.Idx)
}

func TestMultiMasterKey(t *testing.T) {
	p := []byte("testing password")
	p2 := []byte("testing password2")
	pass := buf.NewSecure().Copy(&p).(*buf.Secure)
	pass2 := buf.NewSecure().Copy(&p2).(*buf.Secure)
	wdb := NewWalletDB()
	if wdb.OK() {
		BC := bc.New().Generate(pass).Arm()
		wdb.WithBC(BC)
		BC2 := bc.New().CopyCipher(pass2, BC)
		wdb.WriteMasterKey(BC2)
		BCs := wdb.ReadMasterKeys()
		BCs[0].Unlock(pass).Arm()
		BCs[1].Unlock(pass2).Arm()
		teststring := "the quick brown fox jumped over the lazy dog"
		testbytes := []byte(teststring)
		enc1 := BCs[0].Encrypt(&testbytes)
		enc2 := BCs[1].Encrypt(&testbytes)
		if bytes.Compare(*enc1, *enc2) != 0 {
			t.Error("did not successfully create two identical blockcrypts")
		}
		wdb.EraseMasterKey(BCs[0].Idx)
		wdb.EraseMasterKey(BCs[1].Idx)

		wdb.Close()
	}
}

func TestReadWriteEraseKeyEncryptDecrypt(t *testing.T) {
	p := []byte("testing password")
	pass := buf.NewSecure().Copy(&p).(*buf.Secure)
	BC := bc.New().Generate(pass).Arm()
	wdb := NewWalletDB()
	if wdb.OK() {
		defer wdb.Close()
	}
	wdb.WithBC(BC)
	BCs := wdb.ReadMasterKeys()
	bc := BCs[0]
	bc.Unlock(pass).Arm()
	pk := key.NewPriv()
	pk.WithBC(bc)
	pk.Make()
	wdb.WriteKey(pk)
	addr := pk.GetID()
	address := []byte(addr)
	rpk := wdb.ReadKey(&address)
	if bytes.Compare(*pk.Bytes(), *rpk.Bytes()) != 0 {
		t.Error("failed to write and read back")
	}
	wdb.RemoveBC()
	addr = pk.GetID()
	address = []byte(addr)
	rpk = wdb.ReadKey(&address)
	if bytes.Compare(*pk.Bytes(), *rpk.Bytes()) != 0 {
		t.Error("failed to remove masterkey encryption and read back")
	}
	wdb.WithBC(bc)
	addr = rpk.GetID()
	address = []byte(addr)
	rrpk := wdb.ReadKey(&address)
	if bytes.Compare(*rrpk.Bytes(), *rpk.Bytes()) != 0 {
		t.Error("failed to re-add masterkey encryption and read back")
	}
	wdb.EraseMasterKey(BC.Idx)
	wdb.EraseKey(&address)
	wdb.ReadKey(&address)
	if wdb.Error() != "Key not found" {
		t.Error("failed delete key")
	}
	wdb.deleteAll()
}

func TestEncryptDecrypt(t *testing.T) {
	p := []byte("testing password")
	pass := buf.NewSecure().Copy(&p).(*buf.Secure)
	BC := bc.New().Generate(pass).Arm()
	wdb := NewWalletDB()
	if wdb.OK() {
		defer wdb.Close()
	}
	wdb.dump()
	wdb.WithBC(BC)
	wdb.dump()
	BCs := wdb.ReadMasterKeys()
	wdb.RemoveBC()
	bc := BCs[0]
	bc.Unlock(pass).Arm()
	wdb.WithBC(bc)
	wdb.dump()

	pk := key.NewPriv().WithBC(bc).Make()
	wdb.WriteKey(pk)
	a := key.NewPriv().WithBC(bc).Make()
	aa := []byte(a.GetID())
	label := []byte("random label for random address")
	wdb.WriteName(&aa, &label)
	b := key.NewPriv().WithBC(bc).Make()
	bb := []byte(b.GetID())
	wdb.WriteAccount(&bb, b.PubKey().Bytes())
	wdb.dump()

	wdb.RemoveBC()
	wdb.dump()

	wdb.WithBC(bc)
	wdb.dump()

	wdb.EraseMasterKey(bc.Idx)
	kid := []byte(pk.GetID())
	wdb.EraseKey(&kid)
	wdb.EraseName(&aa)
	wdb.EraseAccount(&bb)

	wdb.deleteAll()
	fmt.Println()
}

func TestJustDump(t *testing.T) {
	wdb := NewWalletDB()
	if wdb.OK() {
		defer wdb.Close()
	}
	wdb.dump()
}

func TestJustDeleteAll(t *testing.T) {
	wdb := NewWalletDB()
	if wdb.OK() {
		wdb.deleteAll()
	}
	wdb.Close()
}
