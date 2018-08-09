package wallet
import (
	"bytes"
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
	for i := 0; i<Flast; i++ {
		db.NewTable(KeyNames[i])
	}
	imp, err := Import()
	if err != nil {
		t.Error("failed to import wallet", err)
	}
	for i := range imp.Names {
		db.WriteName(&imp.Names[i])
	}
	md := new(KeyMetadata)
	for i := range imp.CKeys {
		for j := range imp.Metadata {
			if bytes.Compare(imp.CKeys[i].Pub.Key(), imp.Metadata[j].Pub.Key()) == 0 {
				md.Pub = imp.CKeys[i].Pub
				md.Version = imp.Metadata[j].Version
				md.CreateTime = imp.Metadata[j].CreateTime.Unix()
				break
			}
		}
		db.WriteCryptedKey(imp.CKeys[i].Pub, imp.CKeys[i].Priv, md) 
	}
	for i := range imp.Keys {
		for j := range imp.Metadata {
			if bytes.Compare(imp.Keys[i].Pub.Key(), imp.Metadata[j].Pub.Key()) == 0 {
				md.Pub = imp.Keys[i].Pub
				md.Version = imp.Metadata[j].Version
				md.CreateTime = imp.Metadata[j].CreateTime.Unix()
				break
			}
		}
		db.WriteKey(&imp.Keys[i], md) 
	}
	for i := range imp.WKeys {
		db.WriteWalletKey(&imp.WKeys[i]) 
	}
	for i := range imp.MKeys {
		db.WriteMasterKey(&imp.MKeys[i])
	}
	logger.Debug("Dump:\n", db.Dump())
	db.Close()
}
