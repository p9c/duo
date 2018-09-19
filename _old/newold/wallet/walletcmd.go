package wallet

import (
	"github.com/parallelcointeam/duo/pkg/Uint"
	"github.com/parallelcointeam/duo/pkg/block"
	"github.com/parallelcointeam/duo/pkg/key"
	"github.com/parallelcointeam/duo/pkg/tx"
)

// AddCryptedKey -
func (w *Wallet) AddCryptedKey(*key.Pub, *KeyMetadata) (success bool) { return }

// AddKeyPair -
func (w *Wallet) AddKeyPair(*key.Priv, *key.Pub) (success bool) { return }

// AddReserveKey -
func (w *Wallet) AddReserveKey(*KeyPool) int64 { return 0 }

// AddScript -
func (w *Wallet) AddScript(*key.Script) (success bool) { return }

// AddToWallet -
func (w *Wallet) AddToWallet(Tx) (success bool) { return }

// AddToWalletIfInvolvingMe -
func (w *Wallet) AddToWalletIfInvolvingMe(*Uint.U256, *tx.Transaction, *block.Block, bool, bool) (success bool) {
	return
}

// AvailableCoins -
func (w *Wallet) AvailableCoins([]TxOutput, bool) {}

// CanSupportFeature -
func (w *Wallet) CanSupportFeature(int) (success bool) { return }

// ChangeWalletPassphrase -
func (w *Wallet) ChangeWalletPassphrase(string, string) (success bool) { return }

// CommitTransaction -
func (w *Wallet) CommitTransaction(*Tx, *ReserveKey) (success bool) { return }

// CreateTransaction -
func (w *Wallet) CreateTransaction(*key.Script, int64, *Tx, *ReserveKey, int64, string) (success bool) {
	return
}

// CreateTransactions -
func (w *Wallet) CreateTransactions([]map[*key.Script]int64, *Tx, *ReserveKey, int64, string) (success bool) {
	return
}

// DelAddressBookName -
func (w *Wallet) DelAddressBookName(*key.TxDestination) (success bool) { return }

// EncryptWallet -
func (w *Wallet) EncryptWallet(string) {}

// EraseFromWallet -
func (w *Wallet) EraseFromWallet(*Uint.U256) (success bool) { return }

// GenerateNewKey -
func (w *Wallet) GenerateNewKey() *key.Pub { return nil }

// GetAddressBalances -
func (w *Wallet) GetAddressBalances() map[*key.TxDestination]int64 { return nil }

// GetAddressGroupings -
func (w *Wallet) GetAddressGroupings() []key.TxDestination { return nil }

// GetAllReserveKeys -
func (w *Wallet) GetAllReserveKeys() []key.ID { return nil }

// GetBalance -
func (w *Wallet) GetBalance() int64 { return 0 }

// GetChange -
func (w *Wallet) GetChange(*tx.Out) int64 { return 0 }

// GetCredit -
func (w *Wallet) GetCredit(*tx.Out) int64 { return 0 }

// GetDebit -
func (w *Wallet) GetDebit(*tx.In) int64 { return 0 }

// GetImmatureBalance -
func (w *Wallet) GetImmatureBalance() int64 { return 0 }

// GetKeyBirthTimes -
func (w *Wallet) GetKeyBirthTimes(map[*key.ID]int64) {}

// GetKeyFromPool -
func (w *Wallet) GetKeyFromPool(*key.Pub, bool) (success bool) { return }

// GetKeyPoolSize -
func (w *Wallet) GetKeyPoolSize() int { return 0 }

// GetOldestKeyPoolTime -
func (w *Wallet) GetOldestKeyPoolTime() int64 { return 0 }

// GetTransaction -
func (w *Wallet) GetTransaction(*Uint.U256, *Tx) (success bool) { return }

// GetTxChange -
func (w *Wallet) GetTxChange(*tx.Transaction) int64 { return 0 }

// GetTxCredit -
func (w *Wallet) GetTxCredit(*tx.Transaction) int64 { return 0 }

// GetTxDebit -
func (w *Wallet) GetTxDebit(*tx.Transaction) int64 { return 0 }

// GetUnconfirmedBalance -
func (w *Wallet) GetUnconfirmedBalance() int64 { return 0 }

// GetVersion -
func (w *Wallet) GetVersion() int { return 0 }

// IncOrderPosNext -
func (w *Wallet) IncOrderPosNext(*DB) int64 { return 0 }

// Inventory -
func (w *Wallet) Inventory(*Uint.U256) {}

// IsChange -
func (w *Wallet) IsChange(*tx.Out) (success bool) { return }

// IsFromMe -
func (w *Wallet) IsFromMe(*tx.Transaction) (success bool) { return }

// IsLockedCoin -
func (w *Wallet) IsLockedCoin(*Uint.U256, uint) (success bool) { return }

// IsMyTX -
func (w *Wallet) IsMyTX(*tx.Transaction) (success bool) { return }

// IsMyTxIn -
func (w *Wallet) IsMyTxIn(*tx.In) (success bool) { return }

// IsMyTxOut -
func (w *Wallet) IsMyTxOut(*tx.Out) (success bool) { return }

// KeepKey -
func (w *Wallet) KeepKey(int64) {}

// ListLockedCoins -
func (w *Wallet) ListLockedCoins([]tx.OutPoint) {}

// LoadCryptedKey -
func (w *Wallet) LoadCryptedKey(*key.Pub, []byte) (success bool) { return }

// LoadKey -
func (w *Wallet) LoadKey(*key.Priv, *key.Pub) (success bool) { return }

// LoadKeyMetadata -
func (w *Wallet) LoadKeyMetadata(*key.Pub, *KeyMetadata) (success bool) { return }

// LoadMinVersion -
func (w *Wallet) LoadMinVersion(int) (success bool) { return }

// LoadScript -
func (w *Wallet) LoadScript(*key.Script) (success bool) { return }

// LoadWallet -
func (w *Wallet) LoadWallet(bool) error { return nil }

// LockCoin -
func (w *Wallet) LockCoin(*tx.OutPoint) {}

// MarkDirty -
func (w *Wallet) MarkDirty() {}

// NewKeyPool -
func (w *Wallet) NewKeyPool() (success bool) { return }

// NotifyAddressBookChanged -
func (w *Wallet) NotifyAddressBookChanged(*Wallet, *key.TxDestination, string, bool, int) {}

// NotifyTransactionChanged -
func (w *Wallet) NotifyTransactionChanged(*Wallet, *Uint.U256, int) {}

// OrderedTxItems -
func (w *Wallet) OrderedTxItems([]AccountingEntry, string) *TxItems { return nil }

// PrintWallet -
func (w *Wallet) PrintWallet(*block.Block) {}

// ReacceptWalletTransactions -
func (w *Wallet) ReacceptWalletTransactions() {}

// ResendWalletTransactions -
func (w *Wallet) ResendWalletTransactions() {}

// ReserveKeyFromKeyPool -
func (w *Wallet) ReserveKeyFromKeyPool(int64, *KeyPool) {}

// ReturnKey -
func (w *Wallet) ReturnKey(int64) {}

// ScanForWalletTransactions -
func (w *Wallet) ScanForWalletTransactions(*block.Index, bool) int { return 0 }

// SelectCoinsMinConf -
func (w *Wallet) SelectCoinsMinConf(int64, int, int, []TxOutput) (err error) { return }

// SendMoney -
func (w *Wallet) SendMoney(*key.Script, int64, *Tx, bool) string { return "" }

// SendMoneyToDestination -
func (w *Wallet) SendMoneyToDestination(*key.TxDestination) string { return "" }

// SetAddressBookName -
func (w *Wallet) SetAddressBookName(*key.TxDestination, string) (success bool) { return }

// SetBestChain -
func (w *Wallet) SetBestChain(*block.Locator) {}

// SetDefaultKey -
func (w *Wallet) SetDefaultKey(*key.Pub) (success bool) { return }

// SetMaxVersion -
func (w *Wallet) SetMaxVersion(int) (success bool) { return }

// SetMinVersion -
func (w *Wallet) SetMinVersion(int, *DB, bool) (success bool) { return }

// TopUpKeyPool -
func (w *Wallet) TopUpKeyPool() (success bool) { return }

// Unlock -
func (w *Wallet) Unlock(string) (success bool) { return }

// UnlockAllCoins -
func (w *Wallet) UnlockAllCoins() {}

// UnlockCoin -
func (w *Wallet) UnlockCoin(*tx.OutPoint) {}

// UpdatedTransaction -
func (w *Wallet) UpdatedTransaction(*Uint.U256) {}

// WalletUpdateSpent -
func (w *Wallet) WalletUpdateSpent(*tx.Transaction) {}
