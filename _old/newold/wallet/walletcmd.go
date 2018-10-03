package wallet

import (
	"github.com/parallelcointeam/duo/pkg/Uint"
	"github.com/parallelcointeam/duo/pkg/block"
	"github.com/parallelcointeam/duo/pkg/key"
	"github.com/parallelcointeam/duo/pkg/tx"
)

// AddCryptedKey -
func (r *Wallet) AddCryptedKey(*key.Pub, *KeyMetadata) (success bool) { return }

// AddKeyPair -
func (r *Wallet) AddKeyPair(*key.Priv, *key.Pub) (success bool) { return }

// AddReserveKey -
func (r *Wallet) AddReserveKey(*KeyPool) int64 { return 0 }

// AddScript -
func (r *Wallet) AddScript(*key.Script) (success bool) { return }

// AddToWallet -
func (r *Wallet) AddToWallet(Tx) (success bool) { return }

// AddToWalletIfInvolvingMe -
func (r *Wallet) AddToWalletIfInvolvingMe(*Uint.U256, *tx.Transaction, *block.Block, bool, bool) (success bool) {
	return
}

// AvailableCoins -
func (r *Wallet) AvailableCoins([]TxOutput, bool) {}

// CanSupportFeature -
func (r *Wallet) CanSupportFeature(int) (success bool) { return }

// ChangeWalletPassphrase -
func (r *Wallet) ChangeWalletPassphrase(string, string) (success bool) { return }

// CommitTransaction -
func (r *Wallet) CommitTransaction(*Tx, *ReserveKey) (success bool) { return }

// CreateTransaction -
func (r *Wallet) CreateTransaction(*key.Script, int64, *Tx, *ReserveKey, int64, string) (success bool) {
	return
}

// CreateTransactions -
func (r *Wallet) CreateTransactions([]map[*key.Script]int64, *Tx, *ReserveKey, int64, string) (success bool) {
	return
}

// DelAddressBookName -
func (r *Wallet) DelAddressBookName(*key.TxDestination) (success bool) { return }

// EncryptWallet -
func (r *Wallet) EncryptWallet(string) {}

// EraseFromWallet -
func (r *Wallet) EraseFromWallet(*Uint.U256) (success bool) { return }

// GenerateNewKey -
func (r *Wallet) GenerateNewKey() *key.Pub { return nil }

// GetAddressBalances -
func (r *Wallet) GetAddressBalances() map[*key.TxDestination]int64 { return nil }

// GetAddressGroupings -
func (r *Wallet) GetAddressGroupings() []key.TxDestination { return nil }

// GetAllReserveKeys -
func (r *Wallet) GetAllReserveKeys() []key.ID { return nil }

// GetBalance -
func (r *Wallet) GetBalance() int64 { return 0 }

// GetChange -
func (r *Wallet) GetChange(*tx.Out) int64 { return 0 }

// GetCredit -
func (r *Wallet) GetCredit(*tx.Out) int64 { return 0 }

// GetDebit -
func (r *Wallet) GetDebit(*tx.In) int64 { return 0 }

// GetImmatureBalance -
func (r *Wallet) GetImmatureBalance() int64 { return 0 }

// GetKeyBirthTimes -
func (r *Wallet) GetKeyBirthTimes(map[*key.ID]int64) {}

// GetKeyFromPool -
func (r *Wallet) GetKeyFromPool(*key.Pub, bool) (success bool) { return }

// GetKeyPoolSize -
func (r *Wallet) GetKeyPoolSize() int { return 0 }

// GetOldestKeyPoolTime -
func (r *Wallet) GetOldestKeyPoolTime() int64 { return 0 }

// GetTransaction -
func (r *Wallet) GetTransaction(*Uint.U256, *Tx) (success bool) { return }

// GetTxChange -
func (r *Wallet) GetTxChange(*tx.Transaction) int64 { return 0 }

// GetTxCredit -
func (r *Wallet) GetTxCredit(*tx.Transaction) int64 { return 0 }

// GetTxDebit -
func (r *Wallet) GetTxDebit(*tx.Transaction) int64 { return 0 }

// GetUnconfirmedBalance -
func (r *Wallet) GetUnconfirmedBalance() int64 { return 0 }

// GetVersion -
func (r *Wallet) GetVersion() int { return 0 }

// IncOrderPosNext -
func (r *Wallet) IncOrderPosNext(*DB) int64 { return 0 }

// Inventory -
func (r *Wallet) Inventory(*Uint.U256) {}

// IsChange -
func (r *Wallet) IsChange(*tx.Out) (success bool) { return }

// IsFromMe -
func (r *Wallet) IsFromMe(*tx.Transaction) (success bool) { return }

// IsLockedCoin -
func (r *Wallet) IsLockedCoin(*Uint.U256, uint) (success bool) { return }

// IsMyTX -
func (r *Wallet) IsMyTX(*tx.Transaction) (success bool) { return }

// IsMyTxIn -
func (r *Wallet) IsMyTxIn(*tx.In) (success bool) { return }

// IsMyTxOut -
func (r *Wallet) IsMyTxOut(*tx.Out) (success bool) { return }

// KeepKey -
func (r *Wallet) KeepKey(int64) {}

// ListLockedCoins -
func (r *Wallet) ListLockedCoins([]tx.OutPoint) {}

// LoadCryptedKey -
func (r *Wallet) LoadCryptedKey(*key.Pub, []byte) (success bool) { return }

// LoadKey -
func (r *Wallet) LoadKey(*key.Priv, *key.Pub) (success bool) { return }

// LoadKeyMetadata -
func (r *Wallet) LoadKeyMetadata(*key.Pub, *KeyMetadata) (success bool) { return }

// LoadMinVersion -
func (r *Wallet) LoadMinVersion(int) (success bool) { return }

// LoadScript -
func (r *Wallet) LoadScript(*key.Script) (success bool) { return }

// LoadWallet -
func (r *Wallet) LoadWallet(bool) error { return nil }

// LockCoin -
func (r *Wallet) LockCoin(*tx.OutPoint) {}

// MarkDirty -
func (r *Wallet) MarkDirty() {}

// NewKeyPool -
func (r *Wallet) NewKeyPool() (success bool) { return }

// NotifyAddressBookChanged -
func (r *Wallet) NotifyAddressBookChanged(*Wallet, *key.TxDestination, string, bool, int) {}

// NotifyTransactionChanged -
func (r *Wallet) NotifyTransactionChanged(*Wallet, *Uint.U256, int) {}

// OrderedTxItems -
func (r *Wallet) OrderedTxItems([]AccountingEntry, string) *TxItems { return nil }

// PrintWallet -
func (r *Wallet) PrintWallet(*block.Block) {}

// ReacceptWalletTransactions -
func (r *Wallet) ReacceptWalletTransactions() {}

// ResendWalletTransactions -
func (r *Wallet) ResendWalletTransactions() {}

// ReserveKeyFromKeyPool -
func (r *Wallet) ReserveKeyFromKeyPool(int64, *KeyPool) {}

// ReturnKey -
func (r *Wallet) ReturnKey(int64) {}

// ScanForWalletTransactions -
func (r *Wallet) ScanForWalletTransactions(*block.Index, bool) int { return 0 }

// SelectCoinsMinConf -
func (r *Wallet) SelectCoinsMinConf(int64, int, int, []TxOutput) (err error) { return }

// SendMoney -
func (r *Wallet) SendMoney(*key.Script, int64, *Tx, bool) string { return "" }

// SendMoneyToDestination -
func (r *Wallet) SendMoneyToDestination(*key.TxDestination) string { return "" }

// SetAddressBookName -
func (r *Wallet) SetAddressBookName(*key.TxDestination, string) (success bool) { return }

// SetBestChain -
func (r *Wallet) SetBestChain(*block.Locator) {}

// SetDefaultKey -
func (r *Wallet) SetDefaultKey(*key.Pub) (success bool) { return }

// SetMaxVersion -
func (r *Wallet) SetMaxVersion(int) (success bool) { return }

// SetMinVersion -
func (r *Wallet) SetMinVersion(int, *DB, bool) (success bool) { return }

// TopUpKeyPool -
func (r *Wallet) TopUpKeyPool() (success bool) { return }

// Unlock -
func (r *Wallet) Unlock(string) (success bool) { return }

// UnlockAllCoins -
func (r *Wallet) UnlockAllCoins() {}

// UnlockCoin -
func (r *Wallet) UnlockCoin(*tx.OutPoint) {}

// UpdatedTransaction -
func (r *Wallet) UpdatedTransaction(*Uint.U256) {}

// WalletUpdateSpent -
func (r *Wallet) WalletUpdateSpent(*tx.Transaction) {}
