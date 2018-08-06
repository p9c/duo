package wallet
import (
	"gitlab.com/parallelcoin/duo/pkg/block"
	"gitlab.com/parallelcoin/duo/pkg/key"
	"gitlab.com/parallelcoin/duo/pkg/tx"
	"gitlab.com/parallelcoin/duo/pkg/Uint"
)
func (w *Wallet) AddCryptedKey(*key.Pub, *KeyMetadata) (success bool) { return }
func (w *Wallet) AddKeyPair(*key.Priv, *key.Pub) (success bool) { return }
func (w *Wallet) AddReserveKey(*KeyPool) int64 { return 0 }
func (w *Wallet) AddScript(*key.Script) (success bool) { return }
func (w *Wallet) AddToWallet(Tx) (success bool) { return }
func (w *Wallet) AddToWalletIfInvolvingMe(*Uint.U256, *tx.Transaction, *block.Block, bool, bool) (success bool) { return }
func (w *Wallet) AvailableCoins([]Output, bool) { }
func (w *Wallet) CanSupportFeature(int) (success bool) { return }
func (w *Wallet) ChangeWalletPassphrase(string, string) (success bool) { return }
func (w *Wallet) CommitTransaction(*Tx, *ReserveKey) (success bool) { return }
func (w *Wallet) CreateTransaction(*key.Script, int64, *Tx, *ReserveKey, int64, string) (success bool) { return }
func (w *Wallet) CreateTransactions([]map[*key.Script]int64, *Tx, *ReserveKey, int64, string) (success bool) { return }
func (w *Wallet) DelAddressBookName(*key.TxDestination) (success bool) { return }
func (w *Wallet) EncryptWallet(string) { }
func (w *Wallet) EraseFromWallet(*Uint.U256) (success bool) { return }
func (w *Wallet) GenerateNewKey() *key.Pub { return nil }
func (w *Wallet) GetAddressBalances() map[*key.TxDestination]int64 { return nil }
func (w *Wallet) GetAddressGroupings() []key.TxDestination { return nil }
func (w *Wallet) GetAllReserveKeys() []key.ID { return nil }
func (w *Wallet) GetBalance() int64 { return 0 }
func (w *Wallet) GetChange(*tx.Out) int64 { return 0 }
func (w *Wallet) GetCredit(*tx.Out) int64 { return 0 }
func (w *Wallet) GetDebit(*tx.In) int64 { return 0 }
func (w *Wallet) GetImmatureBalance() int64 { return 0 }
func (w *Wallet) GetKeyBirthTimes(map[*key.ID]int64) { }
func (w *Wallet) GetKeyFromPool(*key.Pub, bool) (success bool) { return }
func (w *Wallet) GetKeyPoolSize() int { return 0 }
func (w *Wallet) GetOldestKeyPoolTime() int64 { return 0 }
func (w *Wallet) GetTransaction(*Uint.U256, *Tx) (success bool) { return }
func (w *Wallet) GetTxChange(*tx.Transaction) int64 { return 0 }
func (w *Wallet) GetTxCredit(*tx.Transaction) int64 { return 0 }
func (w *Wallet) GetTxDebit(*tx.Transaction) int64 { return 0 }
func (w *Wallet) GetUnconfirmedBalance() int64 { return 0 }
func (w *Wallet) GetVersion() int { return 0 }
func (w *Wallet) IncOrderPosNext(*DB) int64 { return 0 }
func (w *Wallet) Inventory(*Uint.U256) { }
func (w *Wallet) IsChange(*tx.Out) (success bool) { return }
func (w *Wallet) IsFromMe(*tx.Transaction) (success bool) { return }
func (w *Wallet) IsLockedCoin(*Uint.U256, uint) (success bool) { return }
func (w *Wallet) IsMyTX(*tx.Transaction) (success bool) { return }
func (w *Wallet) IsMyTxIn(*tx.In) (success bool) { return }
func (w *Wallet) IsMyTxOut(*tx.Out) (success bool) { return }
func (w *Wallet) KeepKey(int64) { }
func (w *Wallet) ListLockedCoins([]tx.OutPoint) { }
func (w *Wallet) LoadCryptedKey(*key.Pub, []byte) (success bool) { return }
func (w *Wallet) LoadKey(*key.Priv, *key.Pub) (success bool) { return }
func (w *Wallet) LoadKeyMetadata(*key.Pub, *KeyMetadata) (success bool) { return }
func (w *Wallet) LoadMinVersion(int) (success bool) { return }
func (w *Wallet) LoadScript(*key.Script) (success bool) { return }
func (w *Wallet) LoadWallet(bool) error { return nil }
func (w *Wallet) LockCoin(*tx.OutPoint) { }
func (w *Wallet) MarkDirty() { }
func (w *Wallet) NewKeyPool() (success bool) { 	return }
func (w *Wallet) NotifyAddressBookChanged(*Wallet, *key.TxDestination, string, bool, int) { }
func (w *Wallet) NotifyTransactionChanged(*Wallet, *Uint.U256, int) { }
func (w *Wallet) OrderedTxItems([]AccountingEntry, string) *TxItems { return nil }
func (w *Wallet) PrintWallet(*block.Block) { }
func (w *Wallet) ReacceptWalletTransactions() { }
func (w *Wallet) ResendWalletTransactions() { }
func (w *Wallet) ReserveKeyFromKeyPool(int64, *KeyPool) { }
func (w *Wallet) ReturnKey(int64) { }
func (w *Wallet) ScanForWalletTransactions(*block.Index, bool) int { return 0 }
func (w *Wallet) SelectCoinsMinConf(int64, int, int, []Output) (err error) { return }
func (w *Wallet) SendMoney(*key.Script, int64, *Tx, bool) string { return "" }
func (w *Wallet) SendMoneyToDestination(*key.TxDestination) string { return "" }
func (w *Wallet) SetAddressBookName(*key.TxDestination, string) (success bool) { return }
func (w *Wallet) SetBestChain(*block.Locator) { }
func (w *Wallet) SetDefaultKey(*key.Pub) (success bool) { return }
func (w *Wallet) SetMaxVersion(int) (success bool) { return }
func (w *Wallet) SetMinVersion(int, *DB, bool) (success bool) { return }
func (w *Wallet) TopUpKeyPool() (success bool) { return }
func (w *Wallet) Unlock(string) (success bool) { return }
func (w *Wallet) UnlockAllCoins() { }
func (w *Wallet) UnlockCoin(*tx.OutPoint) { }
func (w *Wallet) UpdatedTransaction(*Uint.U256) { }
func (w *Wallet) WalletUpdateSpent(*tx.Transaction) { }
