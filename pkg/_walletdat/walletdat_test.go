package walletdat

import (
	"gitlab.com/parallelcoin/duo/pkg/logger"
	"testing"
)

func TestImport(t *testing.T) {
	if imp, err := Import(); err != nil {
		t.Error(err)
	} else {
		logger.Debug(imp)
	}
}