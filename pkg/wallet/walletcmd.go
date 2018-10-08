package wallet

import (
	"bytes"

	"github.com/parallelcointeam/duo/pkg/bc"
	"github.com/parallelcointeam/duo/pkg/block"
	"github.com/parallelcointeam/duo/pkg/buf"
	"github.com/parallelcointeam/duo/pkg/core"
	"github.com/parallelcointeam/duo/pkg/key"
	"github.com/parallelcointeam/duo/pkg/tx"
	"github.com/parallelcointeam/duo/pkg/wallet/db"
	"github.com/parallelcointeam/duo/pkg/wallet/db/rec"
)

// AddKeyPair -
func (r *Wallet) AddKeyPair(k *key.Priv) *Wallet {
	return r
}

// AddAccount adds a new account (correspondent address and optionally public key)
func (r *Wallet) AddAccount(account *rec.Account) *Wallet {
	allIDs := r.DB.GetAllAccountIDs()
	for i := range allIDs {
		if bytes.Compare(allIDs[i], account.Idx) == 0 {
			r.SetStatus("account already in wallet")
			return r
		}
	}
	r.DB.WriteAccount(account.Address, account.Pub)
	return r
}

// AddScript -
func (r *Wallet) AddScript(script *rec.Script) *Wallet { return r }

// AddTx -
func (r *Wallet) AddTx(tx *tx.Transaction) *Wallet { return r }

// AddToWalletIfInvolvingMe -
func (r *Wallet) AddToWalletIfInvolvingMe(id core.Address, tx *tx.Transaction, block *block.Block, update bool, findblock bool) *Wallet {
	return r
}

// AvailableCoins -
func (r *Wallet) AvailableCoins([]tx.Output, bool) {}

// ChangeWalletPassphrase removes any old master keys and creates a new one based on a given password. If the crypt is not locked the old password is required to change it, and if it's not encrypted we just return an error
func (r *Wallet) ChangeWalletPassphrase(oldp, newp *buf.Secure) *Wallet {
	var BC *bc.BlockCrypt
	if r.DB.BC == nil {
		mks := r.DB.ReadMasterKeys()
		if len(mks) < 1 {
			r.SetStatus("wallet is not encrypted")
			return r
		}
		for i := range mks {
			mks[i].Unlock(oldp).Arm()
			if !mks[i].OK() {
				r.SetStatus("password was incorrect for locked wallet")
			} else {
				BC = mks[i]
				r.UnsetStatus()
				break
			}
		}
		if BC == nil {
			r.SetStatus("did not find key unlocked by given password")
			return r
		}
		for i := range mks {
			r.DB.EraseMasterKey(mks[i].Idx)
		}
	}
	BC = bc.New().CopyCipher(newp, r.DB.BC)
	r.DB.WriteMasterKey(BC)
	return r
}

// CommitTransaction -
func (r *Wallet) CommitTransaction(*tx.Transaction, *ReserveKey) *Wallet { return r }

// CreateTransaction -
func (r *Wallet) CreateTransaction(*rec.Script, int64, *tx.Transaction, *ReserveKey, int64, string) *Wallet {
	return r
}

// CreateTransactions -
func (r *Wallet) CreateTransactions([]map[*rec.Script]int64, *tx.Transaction, *ReserveKey, int64, string) *Wallet {
	return r
}

// DelAddressBookName -
func (r *Wallet) DelAddressBookName(*tx.Destination) *Wallet { return r }

// EncryptWallet -
func (r *Wallet) EncryptWallet(string) {}

// EraseFromWallet -
func (r *Wallet) EraseFromWallet(core.Hash) *Wallet { return r }

// GenerateNewKey -
func (r *Wallet) GenerateNewKey() *key.Pub { return nil }

// GetAddressBalances -
func (r *Wallet) GetAddressBalances() map[*tx.Destination]int64 { return nil }

// GetAddressGroupings -
func (r *Wallet) GetAddressGroupings() []tx.Destination { return nil }

// GetAllReserveKeys -
func (r *Wallet) GetAllReserveKeys() []core.Address { return nil }

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
func (r *Wallet) GetKeyBirthTimes(map[*core.Address]int64) {}

// GetTransaction -
func (r *Wallet) GetTransaction(*core.Hash, *tx.Transaction) *Wallet { return r }

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
func (r *Wallet) IncOrderPosNext(*db.DB) int64 { return 0 }

// Inventory -
func (r *Wallet) Inventory(*core.Hash) {}

// IsChange -
func (r *Wallet) IsChange(*tx.Out) *Wallet { return r }

// IsFromMe -
func (r *Wallet) IsFromMe(*tx.Transaction) *Wallet { return r }

// IsLockedCoin -
func (r *Wallet) IsLockedCoin(*core.Hash, uint) *Wallet { return r }

// IsMyTX -
func (r *Wallet) IsMyTX(*tx.Transaction) *Wallet { return r }

// IsMyTxIn -
func (r *Wallet) IsMyTxIn(*tx.In) *Wallet { return r }

// IsMyTxOut -
func (r *Wallet) IsMyTxOut(*tx.Out) *Wallet { return r }

// KeepKey -
func (r *Wallet) KeepKey(int64) {}

// ListLockedCoins -
func (r *Wallet) ListLockedCoins([]tx.OutPoint) {}

// LoadCryptedKey -
func (r *Wallet) LoadCryptedKey(*key.Pub, []byte) *Wallet { return r }

// LoadKey -
func (r *Wallet) LoadKey(*key.Priv, *key.Pub) *Wallet { return r }

// LoadKeyMetadata -
func (r *Wallet) LoadKeyMetadata(*key.Pub, *KeyMetadata) *Wallet {
	return r
}

// LoadMinVersion -
func (r *Wallet) LoadMinVersion(int) *Wallet { return r }

// LoadScript -
func (r *Wallet) LoadScript(*rec.Script) *Wallet { return r }

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
func (r *Wallet) NotifyAddressBookChanged(*Wallet, *tx.Destination, string, bool, int) {}

// NotifyTransactionChanged -
func (r *Wallet) NotifyTransactionChanged(*Wallet, *core.Hash, int) {}

// OrderedTxItems -
func (r *Wallet) OrderedTxItems([]rec.Accounting, string) *tx.Items { return nil }

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
func (r *Wallet) SelectCoinsMinConf(int64, int, int, []tx.Output) (err error) { return }

// SendMoney -
func (r *Wallet) SendMoney(*rec.Script, int64, *tx.Transaction, bool) string { return "" }

// SendMoneyToDestination -
func (r *Wallet) SendMoneyToDestination(*tx.Destination) string { return "" }

// SetAddressBookName -
func (r *Wallet) SetAddressBookName(*tx.Destination, string) *Wallet { return r }

// SetBestChain -
func (r *Wallet) SetBestChain(*block.Locator) {}

// SetDefaultKey -
func (r *Wallet) SetDefaultKey(*key.Pub) *Wallet { return r }

// SetMaxVersion -
func (r *Wallet) SetMaxVersion(int) *Wallet { return r }

// SetMinVersion -
func (r *Wallet) SetMinVersion(int, *db.DB, bool) *Wallet { return r }

// Unlock -
func (r *Wallet) Unlock(string) *Wallet { return r }

// UnlockAllCoins -
func (r *Wallet) UnlockAllCoins() {}

// UnlockCoin -
func (r *Wallet) UnlockCoin(*tx.OutPoint) {}

// UpdatedTransaction -
func (r *Wallet) UpdatedTransaction(*core.Hash) {}

// WalletUpdateSpent -
func (r *Wallet) WalletUpdateSpent(*tx.Transaction) {}
