package walletdb

// WriteScript writes a script entry to the database
func (r *DB) WriteScript() {}

// EraseScript deletes a script entry from the database
func (r *DB) EraseScript() {}

// WriteDefaultKey updates the default key used by interfaces when receiving payments
func (r *DB) WriteDefaultKey() {

}

// ReadDefaultKey returns the current set default key
func (r *DB) ReadDefaultKey() {

}

// WriteBestBlock gets the current best block entry
func (r *DB) WriteBestBlock() {}

// ReadBestBlock gets the current best block entry
func (r *DB) ReadBestBlock() {}

// ReadPool gets the oldest available pool entry and refreshes the pool after addresses are used
func (r *DB) ReadPool() {

}

// WritePool adds a new pool key to the wallet
func (r *DB) WritePool() {

}

// ErasePool removes a pool key
func (r *DB) ErasePool() {

}

// ReadMinVersion returns the minimum version required to read this database
func (r *DB) ReadMinVersion() {

}

// WriteMinVersion updates the minimum version
func (r *DB) WriteMinVersion() {

}

// ReadAccountingEntry writes an accounting entry based on a transaction
func (r *DB) ReadAccountingEntry() {}

// WriteAccountingEntry writes an accounting entry based on a transaction
func (r *DB) WriteAccountingEntry() {}

// EraseAccountingEntry writes an accounting entry based on a transaction
func (r *DB) EraseAccountingEntry() {}

// GetAccountCreditDebit finds entries in the credit/debit records written related to each input transaction from a list of indexes of accounts of interest
func (r *DB) GetAccountCreditDebit() {}
