package wallet

import (
	"encoding/hex"
	"errors"
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"gitlab.com/parallelcoin/duo/pkg/logger"
	"os"
	"testing"
)

var (
	f = "/tmp/wallet.dat"
)

func TestNewDB(t *testing.T) {
	if db, err := NewDB(f); err != nil {
		t.Error(err)
	} else if err := db.Close(); err != nil {
		t.Error(err)
	} else if err = os.Remove(f); err != nil {
		t.Error(err)
	}
}

func TestPutGetDelName(t *testing.T) {
	keyType, acctS, label := "name", "aYjNn4hsFZeKoChNifH8m9QLtrxTQU1nN9", "test"
	acctB := []byte(acctS)
	if db, err := NewDB(f); err != nil {
		t.Error(err)
	} else {
		if _, err := db.Find(keyType, acctB); err == nil {
			if err := db.EraseName(acctS); err != nil {
				t.Error(err)
			}
		}
		dump, _ := db.Dump()
		logger.Debug(dump)
		if err := db.WriteName(acctS, label); err != nil {
			t.Error(err)
		} else {
			dump, _ := db.Dump()
			logger.Debug(dump)
			if _, err := db.Find(keyType, acctB); err != nil {
				t.Error(errors.New("Could not find key"))
			} else if err := db.EraseName(acctS); err != nil {
				t.Error(err)
			} else {
				dump, _ := db.Dump()
				logger.Debug(dump)
				if err := db.Close(); err != nil {
					t.Error(err)
				} else if err = os.Remove(f); err != nil {
					t.Error(err)
				}
			}
		}
	}
}

func TestPutGetDelTx(t *testing.T) {
	keyType, thX, txX := "tx", "a4a6ffa9acc9fcf383f81215784b6b728099924c700e771f58347d11d7b24b1b", "60900f001d072693cf71df52074c7af8089fe4c6d9edeb8433a8e9b1c03000000000000000f13a4846cab468de79ea99ea2724e5d5fd639d53d858b4bf5000000000000000080621ef4aec21b770e9fcc145600f5507cb6e71b432fd87250000000000000097ff5e3be987fa9a86fb53591b38dfc6c6f48acc4d26ed1a0f0000000000000027bd7a34005d973a00cd862882384422669c3eb97aa2636630000000000000000"
	thB, _ := hex.DecodeString(thX)
	txHash := *Uint.Zero256()
	txHash.SetBytes(thB)
	tx, _ := hex.DecodeString(txX)
	if db, err := NewDB(f); err != nil {
		t.Error(err)
	} else {
		if _, err := db.Find(keyType, thB); err == nil {
			if err := db.EraseTx(txHash); err != nil {
				t.Error(err)
			}
		}
		dump, _ := db.Dump()
		logger.Debug(dump)
		if err := db.WriteTx(txHash, tx); err != nil {
			t.Error(err)
		} else {
			dump, _ := db.Dump()
			logger.Debug(dump)
			if _, err := db.Find(keyType, thB); err != nil {
				t.Error(errors.New("Could not find key"))
			} else if err := db.EraseTx(txHash); err != nil {
				t.Error(err)
			} else {
				dump, _ := db.Dump()
				logger.Debug(dump)
				if err := db.Close(); err != nil {
					t.Error(err)
				} else if err = os.Remove(f); err != nil {
					t.Error(err)
				}
			}
		}

	}
}
