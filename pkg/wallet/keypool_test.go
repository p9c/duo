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
	WW.LoadKeyPool()
	// WW.DB.Dump()
	for i := 0; i < 91; i++ {
		_ = WW.GetKeyFromPool(false)
	}
	// WW.DB.Dump()
	WW.TopUpKeyPool()
	WW.DB.Dump()
	WW.EmptyKeyPool()
	WW.DB.Dump()
	WW.DB.DeleteAll()
}

func TestJustErasePool(t *testing.T) {
	wdb := db.NewWalletDB()
	W := New(wdb)
	defer wdb.Close()
	W.EmptyKeyPool()
}
