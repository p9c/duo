package wallet

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/parallelcointeam/duo/pkg/bc"
	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/hash160"

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
	WW.LoadKeyPool()
	WW.DB.Dump()
	priv := WW.GetKeyFromPool(false)
	I := []byte(*priv.PubKey().Bytes())
	addr := hash160.Sum(&I)
	idx := core.Hash64(addr)
	fmt.Println("idx ", hex.EncodeToString(*idx))
	fmt.Println("addr", hex.EncodeToString([]byte(*addr)))
	fmt.Println("priv", len(priv.Hex())/2, priv.Hex())
	fmt.Println("pub ", priv.PubKey().(*buf.Byte).Len(), hex.EncodeToString(*priv.PubKey().Bytes()))
	WW.DB.DeleteAll()
}
