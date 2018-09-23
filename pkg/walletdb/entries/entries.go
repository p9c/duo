package rec

import (
	"github.com/parallelcointeam/duo/pkg/proto"
)

// MasterKey is the password-encrypted store for the symmetric ciphertext that seecures a wallet database
type MasterKey struct {
	Crypt      []byte
	IV         [12]byte
	Iterations int64
}

// Name is a human readable string label for an address
type Name struct {
	Address [20]byte
	Label   string
}

// Tx is a transaction connected to an address in the wallet
type Tx struct {
	ID                    [32]byte
	Data                  []byte
	Prev                  proto.MerkleTx
	TimeRecvIsTxTime      int64
	TimeRecv              int64
	FromMe                bool
	FromAccount           [20]byte
	Spent                 []byte
	OrderPos              int64
	DebitCached           bool
	CreditCached          bool
	ImmatureCreditCached  bool
	AvailableCreditCached bool
	ChangeCached          bool
}

// Seed is a 'hierarchic deterministic' wallet seed that can spawn many subkeys
type Seed struct {
	Secret [32]byte
}

// Key is a key pair controlled by the user
type Key struct {
	Address [20]byte
	Priv    [32]byte
	Pub     [33]byte
}

// Script is a script stored in the wallet database
type Script struct {
	ID   [20]byte
	Data []byte
}

// Pool is a wallet key pair that has not yet been put to use
type Pool struct {
	Seq int64
	Key
	Created int64
	Expires int64
}

// Setting is a configuration that can be used by wallet implementations
type Setting struct {
	Name  string
	Value []byte
}

// Account is an address and public key of a counterparty account
type Account struct {
	Address [20]byte
	Pub     [33]byte
}

// Accounting is an entry regarding internal movements of funds
type Accounting struct {
	Sequence     int64
	Account      [20]byte
	CreditDebit  int64
	Timestamp    int64
	OtherAccount [20]byte
	Comment      string
	OrderPos     int64
	EntryNo      int64
	Extra        []byte
}

// CreditDebit is an increase or decrease to the balance in the wallet
type CreditDebit struct {
	Address [20]byte
	Amount  int64
}

// BestBlock is the latest block at the time the wallet was last open
type BestBlock struct {
	ID   [32]byte
	Data []byte
}

// MinVersion stores the minimum version number required to support the wallet
type MinVersion struct {
	Number int64
}

// DefaultKey sets the default key for receiving payments to use in interfaces
type DefaultKey struct {
	Address [20]byte
}
