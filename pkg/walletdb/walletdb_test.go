package walletdb

import (
	"bytes"
	"testing"

	"github.com/parallelcointeam/duo/pkg/blockcrypt"
	"github.com/parallelcointeam/duo/pkg/buf"
)

// func TestMasterKey(t *testing.T) {
// 	p := []byte("testing password")
// 	pass := buf.NewSecure().Copy(&p).(*buf.Secure)
// 	BC := bc.New().Generate(pass).Arm()
// 	origcrypt := BC.Crypt.Bytes()
// 	origidx := proto.Hash64(origcrypt)
// 	origiv := BC.IV.Bytes()
// 	origiters := BC.Iterations

// 	wdb := NewWalletDB()
// 	if wdb.OK() {
// 		defer wdb.Close()
// 	}

// 	defer wdb.deleteAll()
// 	wdb.WriteMasterKey(BC)
// 	crypt := BC.Crypt.Bytes()
// 	idx := proto.Hash64(crypt)
// 	BCs := wdb.ReadMasterKeys()
// 	for i := range BCs {
// 		crypt = BCs[i].Crypt.Bytes()
// 		iv := BCs[i].IV.Bytes()
// 		iterations := BCs[i].Iterations
// 		idx = proto.Hash64(crypt)

// 		if bytes.Compare(*idx, *origidx) != 0 {
// 			t.Error("idx not properly retrieved")
// 		}
// 		if bytes.Compare(*crypt, *origcrypt) != 0 {
// 			t.Error("crypt not decrypted properly")
// 		}
// 		if bytes.Compare(*iv, *origiv) != 0 {
// 			t.Error("iv not retrieved properly")
// 		}
// 		if iterations != origiters {
// 			t.Error("iterations not retrieved properly")
// 		}

// 		BCs[i].Unlock(buf.NewSecure().Copy(&p).(*buf.Secure)).Arm()

// 		plaintext := "This is the message in plaintext"
// 		plainbytes := []byte(plaintext)
// 		encrypted := BCs[i].Encrypt(&plainbytes)
// 		decrypted := BCs[i].Decrypt(encrypted)

// 		if bytes.Compare(plainbytes, *decrypted) != 0 {
// 			t.Error("encryption/decryption did not work properly")
// 		}
// 	}
// }

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
		// wdb.dump()
		cipher2 := buf.NewSecure().Copy(BC.Ciphertext.Bytes()).(*buf.Secure)
		BC2 := bc.New()
		BC2.IV = BC.IV
		BC2.Iterations = BC.Iterations
		BC2.LoadCiphertext(cipher2, pass2, BC.IV, BC.Iterations)
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
		wdb.deleteAll()
		wdb.Close()
	}
}

// func TestJustDump(t *testing.T) {
// 	wdb := NewWalletDB()
// 	if wdb.OK() {
// 		defer wdb.Close()
// 	}
// 	wdb.dump()
// }

// func TestJustDeleteAll(t *testing.T) {
// 	wdb := NewWalletDB()
// 	if wdb.OK() {
// 		wdb.deleteAll()
// 	}
// 	wdb.Close()
// }
