package wallet

import (
	"github.com/awnumar/memguard"
	"fmt"
	"testing"
)
func TestCrypto(t *testing.T) {
   db, err := NewDB(f)
	if err != nil {
		t.Error("Failed to open")
	}
	db.Net = "mainnet"
	for i := 0; i<Flast; i++ {
		db.NewTable(KeyNames[i])
	}
	imp, err := Import()
	if err != nil {
		t.Error("failed to import wallet", err)
   }
   pass, err := memguard.NewImmutableFromBytes([]byte(passwd))
   b := make([][]byte, len(imp.CKeys))
   for i := range imp.CKeys {
      b[i] = imp.CKeys[i].Priv
   }
   r, _ := imp.MKeys[0].Decrypt(pass, b...)
   for i := range r {
      fmt.Println(imp.CKeys[i].Pub, r[i])
   }
}