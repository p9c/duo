package walletdb

import (
	"bytes"
	"testing"

	"github.com/parallelcointeam/duo/pkg/blockcrypt"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/proto"
)

func TestMasterKey(t *testing.T) {
	p := []byte("testing password")
	pass := buf.NewSecure().Copy(&p).(*buf.Secure)
	BC := bc.New().Generate(pass).Arm()
	origcrypt := BC.Crypt.Bytes()
	origidx := proto.Hash64(origcrypt)
	origiv := BC.IV.Bytes()
	origiters := BC.Iterations

	wdb := NewWalletDB()
	if wdb.OK() {
		defer wdb.Close()
	}

	defer wdb.deleteAll()
	wdb.WriteMasterKey(BC)
	crypt := BC.Crypt.Bytes()
	idx := proto.Hash64(crypt)
	defer wdb.EraseMasterKey(idx)
	BCs := wdb.ReadMasterKeys()
	for i := range BCs {
		crypt = BCs[i].Crypt.Bytes()
		iv := BCs[i].IV.Bytes()
		iterations := BCs[i].Iterations
		idx = proto.Hash64(crypt)

		if bytes.Compare(*idx, *origidx) != 0 {
			t.Error("idx not properly retrieved")
		}
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
	pass := buf.NewSecure().Copy(&p).(*buf.Secure)
	BC := bc.New().Generate(pass).Arm()
	wdb := NewWalletDB()
	if wdb.OK() {
		defer wdb.Close()
	}
	defer wdb.deleteAll()
	wdb.WriteMasterKey(BC)
	crypt := BC.Crypt.Bytes()
	idx := proto.Hash64(crypt)
	defer wdb.EraseMasterKey(idx)
	BCs := wdb.ReadMasterKeys()
	for i := range BCs {
		crypt = BCs[i].Crypt.Bytes()
		idx = proto.Hash64(crypt)
		BCs[i].Unlock(buf.NewSecure().Copy(&p).(*buf.Secure)).Arm()
	}
}
