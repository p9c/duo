package wallet
import (
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
