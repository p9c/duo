package walletdb

import (
	"bytes"
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

	defer wdb.deleteAll()
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
}

func TestMultiMasterKey(t *testing.T) {
	p := []byte("testing password")
	p2 := []byte("testing password2")
	pass := buf.NewSecure().Copy(&p).(*buf.Secure)
	pass2 := buf.NewSecure().Copy(&p2).(*buf.Secure)
	BC := bc.New().Generate(pass).Arm()
	wdb := NewWalletDB()
	if wdb.OK() {
		wdb.WithBC(BC)
		wdb.WriteMasterKey(BC)
		BC2 := bc.New()
		BC2.LoadCiphertext(BC.Ciphertext, pass2, BC.IV, BC.Iterations)
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
func TestReadWriteEraseKey(t *testing.T) {
	p := []byte("testing password")
	pass := buf.NewSecure().Copy(&p).(*buf.Secure)
	BC := bc.New().Generate(pass).Arm()
	wdb := NewWalletDB()
	if wdb.OK() {
		defer wdb.Close()
	}
	wdb.WithBC(BC)
	wdb.WriteMasterKey(BC)
	BCs := wdb.ReadMasterKeys()
	bc := BCs[0]
	bc.Unlock(pass).Arm()
	wdb.WithBC(bc)
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
