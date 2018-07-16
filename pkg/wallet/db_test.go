package wallet

import (
	"errors"
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

func TestPutGetDel(t *testing.T) {
	keyType, acct, label := "name", "aYjNn4hsFZeKoChNifH8m9QLtrxTQU1nN9", "test"
	if db, err := NewDB(f); err != nil {
		t.Error(err)
	} else {
		if _, err := db.Find(keyType, acct); err == nil {
			if err := db.EraseName(acct); err != nil {
				t.Error(err)
			}
		}
		dump, _ := db.Dump()
		logger.Debug(dump)
		if err := db.WriteName(acct, label); err != nil {
			t.Error(err)
		} else {
			dump, _ := db.Dump()
			logger.Debug(dump)
			if _, err := db.Find(keyType, acct); err != nil {
				t.Error(errors.New("Could not find key"))
			} else if err := db.EraseName(acct); err != nil {
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
