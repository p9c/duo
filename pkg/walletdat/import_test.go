package walletdat

import (
	"gitlab.com/parallelcoin/duo/pkg/logger"
	"testing"
)

func TestImport(t *testing.T) {
	imp, err := Import()
	if err != nil {
		t.Error(err)
	}
	logger.Debug(imp)
}