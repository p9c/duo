package wallet
import (
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
	defer db.Close()
	db.Flush()
}
