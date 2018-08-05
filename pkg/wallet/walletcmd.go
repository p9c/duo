package wallet

import (
	"gitlab.com/parallelcoin/duo/pkg/Uint"
	"gitlab.com/parallelcoin/duo/pkg/block"
	"gitlab.com/parallelcoin/duo/pkg/key"
	"gitlab.com/parallelcoin/duo/pkg/tx"
)

// CanSupportFeature returns true if a feature is supported
func (w *Wallet) CanSupportFeature(int) (success bool) {
	return
}

// AvailableCoins returns the list of UXTOs in the wallet
func (w *Wallet) AvailableCoins([]Output, bool) {
}

// SelectCoinsMinConf -
func (w *Wallet) SelectCoinsMinConf(int64, int, int, []Output) (err error) {
	return
}

// IsLockedCoin returns true if a coin is locked
func (w *Wallet) IsLockedCoin(*Uint.U256, uint) (success bool) {
	return
}

// LockCoin locks a UXTO
func (w *Wallet) LockCoin(*tx.OutPoint) {
}

// UnlockCoin unlocks a UXTO
func (w *Wallet) UnlockCoin(*tx.OutPoint) {
}

// UnlockAllCoins unlocks all existing locks
func (w *Wallet) UnlockAllCoins() {
}

// ListLockedCoins returns a list of locked coins
func (w *Wallet) ListLockedCoins([]tx.OutPoint) {
}

// GenerateNewKey generates a new key
func (w *Wallet) GenerateNewKey() *key.Pub {
	return nil
}

// AddKeyPair adds a new key pair to the wallet
func (w *Wallet) AddKeyPair(*key.Priv, *key.Pub) (success bool) {
	return
}

// LoadKey loads a key
func (w *Wallet) LoadKey(*key.Priv, *key.Pub) (success bool) {
	return
}

// LoadKeyMetadata loads a key's metadata
func (w *Wallet) LoadKeyMetadata(*key.Pub, *KeyMetadata) (success bool) {
	return
}

// LoadMinVersion sets the min version on the wallet
func (w *Wallet) LoadMinVersion(int) (success bool) {
	return
}

// AddCryptedKey adds an encrypted key to the wallet
func (w *Wallet) AddCryptedKey(*key.Pub, *KeyMetadata) (success bool) {
	return
}

// LoadCryptedKey loads an encrypted key from the wallet
func (w *Wallet) LoadCryptedKey(*key.Pub, []byte) (success bool) {
	return
}

// AddScript adds a script to the wallet
func (w *Wallet) AddScript(*key.Script) (success bool) {
	return
}

// LoadScript loads a script in the wallet
func (w *Wallet) LoadScript(*key.Script) (success bool) {
	return
}

// Unlock unlocks the wallet with a passphrase
func (w *Wallet) Unlock(string) (success bool) {
	return
}

// ChangeWalletPassphrase changes the passphrase given the correct existing passphrase
func (w *Wallet) ChangeWalletPassphrase(string, string) (success bool) {
	return
}

// EncryptWallet encrypts the wallet with a given passphrase
func (w *Wallet) EncryptWallet(string) {
}

// GetKeyBirthTimes gets a list of keys with their birth times
func (w *Wallet) GetKeyBirthTimes(map[*key.ID]int64) {
}

// IncOrderPosNext -
func (w *Wallet) IncOrderPosNext(*DB) int64 {
	return 0
}

// OrderedTxItems -
func (w *Wallet) OrderedTxItems([]AccountingEntry, string) *TxItems {
	return nil
}

// MarkDirty marks the wallet in memory dirty and needing sync
func (w *Wallet) MarkDirty() {
}

// AddToWallet adds a transaction to the wallet
func (w *Wallet) AddToWallet(Tx) (success bool) {
	return
}

// AddToWalletIfInvolvingMe adds a transaction if we have a private key related to the transaction
func (w *Wallet) AddToWalletIfInvolvingMe(*Uint.U256, *tx.Transaction, *block.Block, bool, bool) (success bool) {
	return
}

// EraseFromWallet removes a key from a wallet
func (w *Wallet) EraseFromWallet(*Uint.U256) (success bool) {
	return
}

// WalletUpdateSpent updates a transaction's record
func (w *Wallet) WalletUpdateSpent(*tx.Transaction) {
}

// ScanForWalletTransactions searches for and adds transactions from the chain related to wallet keys that are not yet in the wallet
func (w *Wallet) ScanForWalletTransactions(*block.Index, bool) int {
	return 0
}

// ReacceptWalletTransactions accepts wallet transactions again
func (w *Wallet) ReacceptWalletTransactions() {
}

// ResendWalletTransactions sends wallet transactions again
func (w *Wallet) ResendWalletTransactions() {
}

// GetBalance gets the balance of the wallet
func (w *Wallet) GetBalance() int64 {
	return 0
}

// GetUnconfirmedBalance returns the unconfirmed balance
func (w *Wallet) GetUnconfirmedBalance() int64 {
	return 0
}

// GetImmatureBalance gets the balance including mining payments not yet reached maturity
func (w *Wallet) GetImmatureBalance() int64 {
	return 0
}

// CreateTransactions creates a set of transactions from a set of scripts
func (w *Wallet) CreateTransactions([]map[*key.Script]int64, *Tx, *ReserveKey, int64, string) (success bool) {
	return
}

// CreateTransaction creates a new transaction from a script
func (w *Wallet) CreateTransaction(*key.Script, int64, *Tx, *ReserveKey, int64, string) (success bool) {
	return
}

// CommitTransaction commits a transaction
func (w *Wallet) CommitTransaction(*Tx, *ReserveKey) (success bool) {
	return
}

// SendMoney sends out a transaction
func (w *Wallet) SendMoney(*key.Script, int64, *Tx, bool) string {
	return ""
}

// SendMoneyToDestination sends a transaction to a specific destination
func (w *Wallet) SendMoneyToDestination(*key.TxDestination) string {
	return ""
}

// NewKeyPool creatse a new keypool in the wallet
func (w *Wallet) NewKeyPool() (success bool) {
	return
}

// TopUpKeyPool creates new keys to fill the key pool
func (w *Wallet) TopUpKeyPool() (success bool) {
	return
}

// AddReserveKey adds a new reserve key from a key pool
func (w *Wallet) AddReserveKey(*KeyPool) int64 {
	return 0
}

// ReserveKeyFromKeyPool adds a specific key from a keypool
func (w *Wallet) ReserveKeyFromKeyPool(int64, *KeyPool) {
}

// KeepKey fixates a key currently in a keypool
func (w *Wallet) KeepKey(int64) {
}

// ReturnKey returns a key to the key pool
func (w *Wallet) ReturnKey(int64) {
}

// GetKeyFromPool gets a key from the key pool
func (w *Wallet) GetKeyFromPool(*key.Pub, bool) (success bool) {
	return
}

// GetOldestKeyPoolTime gets the time of the oldest key in the pool
func (w *Wallet) GetOldestKeyPoolTime() int64 {
	return 0
}

// GetAllReserveKeys gets all reserve key ID's
func (w *Wallet) GetAllReserveKeys() []key.ID {
	return nil
}

// GetAddressGroupings returns a collection of destination keys
func (w *Wallet) GetAddressGroupings() []key.TxDestination {
	return nil
}

// GetAddressBalances gets the balances of a set of destination keys
func (w *Wallet) GetAddressBalances() map[*key.TxDestination]int64 {
	return nil
}

// IsMyTxIn returns true if a tx is related to a key in the wallet
func (w *Wallet) IsMyTxIn(*tx.In) (success bool) {
	return
}

// GetDebit returns the value of a transaction
func (w *Wallet) GetDebit(*tx.In) int64 {
	return 0
}

// IsMyTxOut returns true if a payment is related to a key in the wallet
func (w *Wallet) IsMyTxOut(*tx.Out) (success bool) {
	return
}

// GetCredit returns the value of an outbound transaction
func (w *Wallet) GetCredit(*tx.Out) int64 {
	return 0
}

// IsChange returns true if an output is a change output (it claimed a new key from the pool)
func (w *Wallet) IsChange(*tx.Out) (success bool) {
	return
}

// GetChange gets the amount of change in an outbound transaction
func (w *Wallet) GetChange(*tx.Out) int64 {
	return 0
}

// IsMyTX returns true if the wallet contains the private key that created a transaction
func (w *Wallet) IsMyTX(*tx.Transaction) (success bool) {
	return
}

// IsFromMe returns true if a transaction was created using a private key in this wallet
func (w *Wallet) IsFromMe(*tx.Transaction) (success bool) {
	return
}

// GetTxDebit returns the value out of a transaction
func (w *Wallet) GetTxDebit(*tx.Transaction) int64 {
	return 0
}

// GetTxCredit returns the value in of a transaction
func (w *Wallet) GetTxCredit(*tx.Transaction) int64 {
	return 0
}

// GetTxChange gets the change involved in a transaction
func (w *Wallet) GetTxChange(*tx.Transaction) int64 {
	return 0
}

// SetBestChain sets the best chain (head block) record in the wallet
func (w *Wallet) SetBestChain(*block.Locator) {
}

// LoadWallet loads a wallet
func (w *Wallet) LoadWallet(bool) error {
	return nil
}

// SetAddressBookName sets a new name in the address book
func (w *Wallet) SetAddressBookName(*key.TxDestination, string) (success bool) {
	return
}

// DelAddressBookName deletes an address book name
func (w *Wallet) DelAddressBookName(*key.TxDestination) (success bool) {
	return
}

// UpdatedTransaction -
func (w *Wallet) UpdatedTransaction(*Uint.U256) {
}

// PrintWallet -
func (w *Wallet) PrintWallet(*block.Block) {
}

// Inventory -
func (w *Wallet) Inventory(*Uint.U256) {

}

// GetKeyPoolSize returns the size of the key pool
func (w *Wallet) GetKeyPoolSize() int {
	return 0
}

// GetTransaction -
func (w *Wallet) GetTransaction(*Uint.U256, *Tx) (success bool) {
	return
}

// SetDefaultKey sets the default key in the wallet
func (w *Wallet) SetDefaultKey(*key.Pub) (success bool) {
	return
}

// SetMinVersion sets the min version of the wallet
func (w *Wallet) SetMinVersion(int, *DB, bool) (success bool) {
	return
}

// SetMaxVersion sets the max version of the wallet
func (w *Wallet) SetMaxVersion(int) (success bool) {
	return
}

// GetVersion gets the version of the wallet
func (w *Wallet) GetVersion() int {
	return 0
}

// NotifyAddressBookChanged is a function that is called when an address book entry changes
func (w *Wallet) NotifyAddressBookChanged(*Wallet, *key.TxDestination, string, bool, int) {
}

// NotifyTransactionChanged is a function that is called when a transaction changes
func (w *Wallet) NotifyTransactionChanged(*Wallet, *Uint.U256, int) {
}
