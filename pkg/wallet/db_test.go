package wallet
import (
	"bytes"
	"gitlab.com/parallelcoin/duo/pkg/logger"
	"gitlab.com/parallelcoin/duo/pkg/walletdat"
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
	logger.Debug(*db)
	for i := range KeyNames {
		db.NewTable(KeyNames[i])
	}
	imp, err := walletdat.Import()
	if err != nil {
		t.Error("failed to import wallet", err)
	}
	for i := range imp.Names {
		db.WriteName(imp.Names[i].Addr, imp.Names[i].Name)
	}
	md := new(KeyMetadata)
	for i := range imp.Keys {
		for j := range imp.Metadata {
			if bytes.Compare(imp.Keys[i].Pub.Key(), imp.Metadata[j].Pub.Key()) == 0 {
				md.Pub = imp.Keys[i].Pub
				md.Version = imp.Metadata[j].Version
				md.CreateTime = imp.Metadata[j].CreateTime.Unix()
			}
		}
		db.WriteKey(imp.Keys[i].Pub, imp.Keys[i].Priv, md) 
	}
	logger.Debug(imp)
	db.Close()
}
