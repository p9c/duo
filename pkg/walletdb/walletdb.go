package walletdb

import (
	"github.com/dgraph-io/badger"
	"github.com/mitchellh/go-homedir"
)

// NewWalletDB creates a new walletdb.DB
func NewWalletDB() (db *DB) {
	db = new(DB)
	opts := badger.DefaultOptions
	var err error
	db.Path, err = homedir.Dir()
	if err != nil {
		db.Status = err.Error()
		return db
	}
	db.DB, err = badger.Open(opts)
	db.Status = err.Error()
	if err != nil {
		db.Status = err.Error()
		return db
	}
	return db
}

// WriteName is a
func (r *DB) WriteName() {}

func (r *DB) EraseName() {}

// ReadTx is a
func (r *DB) ReadTx() {}

// WriteTx is a
func (r *DB) WriteTx() {}

// EraseTx is a
func (r *DB) EraseTx() {}

// WriteKey is a
func (r *DB) WriteKey() {}

// WriteMasterKey is a
func (r *DB) WriteMasterKey() {}

// WriteScript is a
func (r *DB) WriteScript() {}

// WriteBestBlock is a
func (r *DB) WriteBestBlock() {}

// WriteDefaultKey is a
func (r *DB) WriteDefaultKey() {}

// ReadBestBlock is a
func (r *DB) ReadBestBlock() {}

// ReadDefaultKey is a
func (r *DB) ReadDefaultKey() {}

// ReadPool is a
func (r *DB) ReadPool() {}

// WritePool is a
func (r *DB) WritePool() {}

// ReadSetting is a
func (r *DB) ReadSetting() {}

// WriteSetting is a
func (r *DB) WriteSetting() {}

// EraseSetting is a
func (r *DB) EraseSetting() {}

// ReadMinVersion is a
func (r *DB) ReadMinVersion() {}

// WriteMinVersion is a
func (r *DB) WriteMinVersion() {}

// ReadAccount is a
func (r *DB) ReadAccount() {}

// WriteAccount is a
func (r *DB) WriteAccount() {}

// EraseAccount is a
func (r *DB) EraseAccount() {}

// WriteAccountingEntry is a
func (r *DB) WriteAccountingEntry() {}

// GetAccountCreditDebit is a
func (r *DB) GetAccountCreditDebit() {}

// ListAccCreditDebit is a
func (r *DB) ListAccCreditDebit() {}

// ReorderTransactions is a
func (r *DB) ReorderTransactions() {}

// LoadWallet is a
func (r *DB) LoadWallet() {}

// Recover is a
func (r *DB) Recover() {}
