package wallet

import (
	"github.com/parallelcointeam/duo/pkg/block"
	"github.com/parallelcointeam/duo/pkg/key"
	"github.com/parallelcointeam/duo/pkg/proto"
	"github.com/parallelcointeam/duo/pkg/tx"
	"github.com/parallelcointeam/duo/pkg/wallet/db"
	"github.com/parallelcointeam/duo/pkg/wallet/db/rec"
)

// AddKeyPair -
func (r *Wallet) AddKeyPair(k *key.Priv) *Wallet {
	return r
}

// AddAccount adds a new account (correspondent address and optionally public key)
func (r *Wallet) AddAccount(accounnt *rec.Account) *Wallet {
	return r
}

// AddScript -
func (r *Wallet) AddScript(script *rec.Script) *Wallet { return w }

// AddTx -
func (r *Wallet) AddTx(tx *tx.Transaction) *Wallet { return w }

// AddToWalletIfInvolvingMe -
func (r *Wallet) AddToWalletIfInvolvingMe(id proto.Address, tx *tx.Transaction, block *block.Block, update bool, findblock bool) *Wallet {
	return r
}

// AvailableCoins -
func (r *Wallet) AvailableCoins([]tx.TxOutput, bool) {}

// ChangeWalletPassphrase -
func (r *Wallet) ChangeWalletPassphrase(string, string) *Wallet { return w }

// CommitTransaction -
func (r *Wallet) CommitTransaction(*tx.Transaction, *ReserveKey) *Wallet { return w }

// CreateTransaction -
func (r *Wallet) CreateTransaction(*rec.Script, int64, *tx.Transaction, *ReserveKey, int64, string) *Wallet {
	return r
}

// CreateTransactions -
func (r *Wallet) CreateTransactions([]map[*rec.Script]int64, *tx.Transaction, *ReserveKey, int64, string) *Wallet {
	return r
}

// DelAddressBookName -
func (r *Wallet) DelAddressBookName(*key.TxDestination) *Wallet { return w }

// EncryptWallet -
func (r *Wallet) EncryptWallet(string) {}

// EraseFromWallet -
func (r *Wallet) EraseFromWallet(proto.Hash) *Wallet { return w }

// GenerateNewKey -
func (r *Wallet) GenerateNewKey() *key.Pub { return nil }

// GetAddressBalances -
func (r *Wallet) GetAddressBalances() map[*key.TxDestination]int64 { return nil }

// GetAddressGroupings -
func (r *Wallet) GetAddressGroupings() []key.TxDestination { return nil }

// GetAllReserveKeys -
func (r *Wallet) GetAllReserveKeys() []proto.Address { return nil }

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
func (r *Wallet) GetKeyBirthTimes(map[*proto.Address]int64) {}

// GetTransaction -
func (r *Wallet) GetTransaction(*proto.Hash, *tx.Transaction) *Wallet { return w }

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
func (r *Wallet) IncOrderPosNext(*walletdb.DB) int64 { return 0 }

// Inventory -
func (r *Wallet) Inventory(*proto.Hash) {}

// IsChange -
func (r *Wallet) IsChange(*tx.Out) *Wallet { return w }

// IsFromMe -
func (r *Wallet) IsFromMe(*tx.Transaction) *Wallet { return w }

// IsLockedCoin -
func (r *Wallet) IsLockedCoin(*proto.Hash, uint) *Wallet { return w }

// IsMyTX -
func (r *Wallet) IsMyTX(*tx.Transaction) *Wallet { return w }

// IsMyTxIn -
func (r *Wallet) IsMyTxIn(*tx.In) *Wallet { return w }

// IsMyTxOut -
func (r *Wallet) IsMyTxOut(*tx.Out) *Wallet { return w }

// KeepKey -
func (r *Wallet) KeepKey(int64) {}

// ListLockedCoins -
func (r *Wallet) ListLockedCoins([]tx.OutPoint) {}

// LoadCryptedKey -
func (r *Wallet) LoadCryptedKey(*key.Pub, []byte) *Wallet { return w }

// LoadKey -
func (r *Wallet) LoadKey(*key.Priv, *key.Pub) *Wallet { return w }

// LoadKeyMetadata -
func (r *Wallet) LoadKeyMetadata(*key.Pub, *KeyMetadata) *Wallet {
	return r
}

// LoadMinVersion -
func (r *Wallet) LoadMinVersion(int) *Wallet { return w }

// LoadScript -
func (r *Wallet) LoadScript(*rec.Script) *Wallet { return w }

// LoadWallet -
func (r *Wallet) LoadWallet(bool) error { return nil }

// LockCoin -
func (r *Wallet) LockCoin(*tx.OutPoint) *Wallet {
	return r
}

// MarkDirty -
func (r *Wallet) MarkDirty() *Wallet {
	return r
}

// NotifyAddressBookChanged -
func (r *Wallet) NotifyAddressBookChanged(*Wallet, *key.TxDestination, string, bool, int) {}

// NotifyTransactionChanged -
func (r *Wallet) NotifyTransactionChanged(*Wallet, *proto.Hash, int) {}

// OrderedTxItems -
func (r *Wallet) OrderedTxItems([]AccountingEntry, string) *TxItems { return nil }

// PrintWallet -
func (r *Wallet) PrintWallet(*block.Block) {}

// ReacceptWalletTransactions -
func (r *Wallet) ReacceptWalletTransactions() {}

// ResendWalletTransactions -
func (r *Wallet) ResendWalletTransactions() {}

// ReturnKey -
func (r *Wallet) ReturnKey(int64) {}

// ScanForWalletTransactions -
func (r *Wallet) ScanForWalletTransactions(*block.Index, bool) int { return 0 }

// SelectCoinsMinConf -
func (r *Wallet) SelectCoinsMinConf(int64, int, int, []tx.TxOutput) (err error) { return }

// SendMoney -
func (r *Wallet) SendMoney(*rec.Script, int64, *tx.Transaction, bool) string { return "" }

// SendMoneyToDestination -
func (r *Wallet) SendMoneyToDestination(*key.TxDestination) string { return "" }

// SetAddressBookName -
func (r *Wallet) SetAddressBookName(*key.TxDestination, string) *Wallet { return w }

// SetBestChain -
func (r *Wallet) SetBestChain(*block.Locator) {}

// SetDefaultKey -
func (r *Wallet) SetDefaultKey(*key.Pub) *Wallet { return w }

// SetMaxVersion -
func (r *Wallet) SetMaxVersion(int) *Wallet { return w }

// SetMinVersion -
func (r *Wallet) SetMinVersion(int, *walletdb.DB, bool) *Wallet { return w }

// Unlock -
func (r *Wallet) Unlock(string) *Wallet { return w }

// UnlockAllCoins -
func (r *Wallet) UnlockAllCoins() {}

// UnlockCoin -
func (r *Wallet) UnlockCoin(*tx.OutPoint) {}

// UpdatedTransaction -
func (r *Wallet) UpdatedTransaction(*proto.Hash) {}

// WalletUpdateSpent -
func (r *Wallet) WalletUpdateSpent(*tx.Transaction) {}