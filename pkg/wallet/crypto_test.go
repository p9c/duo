package wallet

import (
	"fmt"
	"github.com/awnumar/memguard"
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
	pass, err := memguard.NewImmutableFromBytes([]byte(passwd))
	imp := Import(pass)
	if err != nil {
		t.Error("failed to import wallet", err)
	}
	b := make([][]byte, len(imp.CKeys))
	for i := range imp.CKeys {
		b[i] = imp.CKeys[i].Priv
	}
	r, _ := imp.masterKey[0].Decrypt(pass, b...)
	fmt.Println("Decrypted")
	for i := range r {
		fmt.Println(imp.CKeys[i].Pub, r[i])
	}
	pub := make([][]byte, len(imp.CKeys))
	for i := range imp.CKeys {
		pub[i] = append(imp.CKeys[i].Pub, make([]byte, 31)...)
		for j := 33; j < 64; j++ {
			pub[i][j] = 31
		}
	}
	fmt.Println("\nEncrypted:")
	s, _ := imp.masterKey[0].Encrypt(pass, pub...)
	for i := range s {
		fmt.Println(s[i], r[i])
	}
}
