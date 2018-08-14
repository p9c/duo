package wallet

import (
	"fmt"
	"github.com/awnumar/memguard"
	"gitlab.com/parallelcoin/duo/pkg/logger"
	"testing"
)

var (
	f = "/tmp/wallet"
)

func TestNewDB(t *testing.T) {
	db, err := NewDB(f)
	if err != nil {
		t.Error("Failed to open")
	}
	logger.Debug(*db)
	for i := range KeyNames {
		db.NewTable(KeyNames[i])
	}
	db.Close()
}
func TestImport(t *testing.T) {
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
	es := imp.ToEncryptedStore()
	for i := range es.AddressBook {
		fmt.Println("original   ", string(imp.Names[i].Addr))
		fmt.Println("encrypted  ", (es.AddressBook[i].Pub))
		a := es.AddressBook[i].Decrypt()
		fmt.Println("decrypted  ", string(a.Pub))
		// a = a.Encrypt()
		b := make([]byte, 14)
		for j := range b {
			b[j] = 14
		}
		pub := make([]byte, 48)
		es.EncryptData(pub, append(a.Pub, b...))
		fmt.Println("reencrypted", pub)
		a.Destroy()
	}
	test := []byte("this is a test! ")
	fmt.Println(len(test), test, string(test))
	testenc := make([]byte, 16)
	testdec := make([]byte, 32)
	es.EncryptData(testenc, test)
	fmt.Println(len(testenc), testenc, string(testenc))
	es.DecryptData(testdec, testenc)
	fmt.Println(len(testdec), testdec, string(testdec))
	db.Close()
}

func TestJSON(t *testing.T) {
	es := new(EncryptedStore)
	fmt.Println(string(ToJSON(es)))
}
