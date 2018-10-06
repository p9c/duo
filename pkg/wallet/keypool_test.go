package wallet

import (
	"testing"

	"github.com/parallelcointeam/duo/pkg/bc"

	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/wallet/db"
)

func TestNewKeyPool(t *testing.T) {
	p := []byte("testing password")
	pass := buf.NewSecure().Copy(&p).(*buf.Secure)
	BC := bc.New().Generate(pass).Arm()
	wdb := db.NewWalletDB()
	wdb.WithBC(BC)
	W := New(wdb)
	W.NewKeyPool()
	W.DB.Close()
	wwdb := db.NewWalletDB()
	wwdb.WithBC(BC)
	WW := New(wwdb)
	defer WW.DB.Close()
	defer WW.DB.DeleteAll()
	// W.DB.WithBC(BC)
	WW.LoadKeyPool()
	// fmt.Println(WW.KeyPool.Size)
	// WW.DB.RemoveBC()
	// WW.DB.Dump()
	WW.EmptyKeyPool()
	WW.DB.Dump()
}
